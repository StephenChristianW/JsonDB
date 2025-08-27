package services

import (
	UtilsFile "JsonDB/utils/file"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"regexp"
	"strings"
)

// ---------------- utils ----------------

func generateObjectID() string {
	return primitive.NewObjectID().Hex()
}

func getCollectionFilePath(db *DBContext) (string, error) {
	if db.CurrentDB == "" || db.CurrentCollection == "" {
		return "", errors.New("数据库或集合未选择")
	}
	return db.getCollectionFilePath(db.CurrentCollection)
}

func getIndexFilePath(db *DBContext, field string) (string, error) {
	if db.CurrentDB == "" || db.CurrentCollection == "" {
		return "", errors.New("数据库或集合未选择")
	}
	return fmt.Sprintf("MyDB/%s/%s.%s.index", db.CurrentDB, db.CurrentCollection, field), nil
}

// ---------------- load/save ----------------

func loadCollection(db *DBContext) (map[string]Document, error) {
	colPath, err := getCollectionFilePath(db)
	if err != nil {
		return nil, err
	}

	data := make(map[string]Document)
	if UtilsFile.IsPathExist(colPath) {
		bytes, err := os.ReadFile(colPath)
		if err != nil {
			return nil, err
		}
		if len(bytes) > 0 {
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
		}
	}
	return data, nil
}

func saveCollection(db *DBContext, data map[string]Document) error {
	colPath, err := getCollectionFilePath(db)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(colPath, bytes, 0666)
}

// ---------------- index utils ----------------

func loadIndex(db *DBContext, field string) (map[interface{}][]string, error) {
	path, err := getIndexFilePath(db, field)
	if err != nil {
		return nil, err
	}
	index := make(map[interface{}][]string)
	if UtilsFile.IsPathExist(path) {
		bytes, _ := os.ReadFile(path)
		if len(bytes) > 0 {
			_ = json.Unmarshal(bytes, &index)
		}
	}
	return index, nil
}

func saveIndex(db *DBContext, field string, index map[interface{}][]string) error {
	path, err := getIndexFilePath(db, field)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0666)
}

func updateIndex(db *DBContext, docID string, doc Document, fields []string, remove bool) {
	for _, field := range fields {
		val, _ := getNestedValue(doc, field)
		index, _ := loadIndex(db, field)
		if remove {
			if arr, ok := index[val]; ok {
				newArr := make([]string, 0, len(arr))
				for _, id := range arr {
					if id != docID {
						newArr = append(newArr, id)
					}
				}
				index[val] = newArr
			}
		} else {
			index[val] = append(index[val], docID)
		}
		_ = saveIndex(db, field, index)
	}
}

// ---------------- match ----------------

func matchDoc(doc Document, filter map[string]interface{}) bool {
	for k, v := range filter {
		switch k {
		case "$and":
			arr, ok := v.([]interface{})
			if !ok {
				return false
			}
			for _, cond := range arr {
				if condMap, ok := cond.(map[string]interface{}); ok {
					if !matchDoc(doc, condMap) {
						return false
					}
				}
			}
		case "$or":
			arr, ok := v.([]interface{})
			if !ok {
				return false
			}
			matched := false
			for _, cond := range arr {
				if condMap, ok := cond.(map[string]interface{}); ok {
					if matchDoc(doc, condMap) {
						matched = true
						break
					}
				}
			}
			if !matched {
				return false
			}
		case "$not":
			if condMap, ok := v.(map[string]interface{}); ok {
				if matchDoc(doc, condMap) {
					return false
				}
			}
		default:
			val, ok := getNestedValue(doc, k)
			if !ok {
				// 模糊匹配字段名
				for field := range doc {
					if strings.Contains(field, k) {
						val, ok = doc[field]
						break
					}
				}
				if !ok {
					return false
				}
			}

			if condMap, ok := v.(map[string]interface{}); ok {
				for op, cond := range condMap {
					if !matchOperator(val, op, cond) {
						return false
					}
				}
			} else {
				if !matchOperator(val, "$eq", v) {
					return false
				}
			}
		}
	}
	return true
}

func getNestedValue(doc Document, field string) (interface{}, bool) {
	parts := strings.Split(field, ".")
	var val interface{} = doc
	for _, p := range parts {
		if m, ok := val.(map[string]interface{}); ok {
			val, ok = m[p]
			if !ok {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return val, true
}

func matchOperator(value interface{}, op string, cond interface{}) bool {
	switch op {
	case "$eq":
		return value == cond
	case "$ne":
		return value != cond
	case "$gt":
		return compareNumber(value, cond) > 0
	case "$gte":
		return compareNumber(value, cond) >= 0
	case "$lt":
		return compareNumber(value, cond) < 0
	case "$lte":
		return compareNumber(value, cond) <= 0
	case "$in":
		arr, ok := cond.([]interface{})
		if !ok {
			return false
		}
		for _, v := range arr {
			if value == v {
				return true
			}
		}
		return false
	case "$nin":
		arr, ok := cond.([]interface{})
		if !ok {
			return false
		}
		for _, v := range arr {
			if value == v {
				return false
			}
		}
		return true
	case "$regex":
		s, ok := value.(string)
		if !ok {
			return false
		}
		pattern, ok := cond.(string)
		if !ok {
			return false
		}
		matched, _ := regexp.MatchString(pattern, s)
		return matched
	default:
		return false
	}
}

func compareNumber(a, b interface{}) int {
	var fa, fb float64
	switch v := a.(type) {
	case int:
		fa = float64(v)
	case float64:
		fa = v
	default:
		return 0
	}
	switch v := b.(type) {
	case int:
		fb = float64(v)
	case float64:
		fb = v
	default:
		return 0
	}
	if fa > fb {
		return 1
	} else if fa < fb {
		return -1
	}
	return 0
}
