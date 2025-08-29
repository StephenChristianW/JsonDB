package services

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	ConfigFile "github.com/StephenChristianW/JsonDB/fileIO/configFileIO"
)

var JsonMu sync.RWMutex

type Document map[string]interface{}
type DocumentList []Document

type FindOptions struct {
	Sort   map[string]int // 1升序，-1降序
	Skip   int
	Limit  int
	Fields []string // 投影，可选
}

type DocServices interface {
	Find(filter map[string]interface{}, opts *FindOptions) (DocumentList, error)
	FindOne(filter map[string]interface{}) (Document, error)
	InsertOne(doc Document) (Document, error)
	InsertMany(docs []Document) ([]Document, error)
	UpdateOne(filter map[string]interface{}, update Document) (Document, error)
	UpdateMany(filter map[string]interface{}, update Document) ([]Document, error)
	Delete(filter map[string]interface{}) (int, error)
}

// ---------------- DBContext 文档操作 ----------------

func (db *DBContext) Find(filter map[string]interface{}, opts *FindOptions) (DocumentList, error) {
	JsonMu.RLock()
	defer JsonMu.RUnlock()

	data, err := loadCollection(db)
	if err != nil {
		return nil, err
	}

	// 检查索引字段
	indexFields := ConfigFile.GetIndexFields(db.CurrentDB, db.CurrentCollection)
	candidateIDs := make(map[string]struct{})
	for _, field := range indexFields {
		if val, ok := filter[field]; ok {
			indexMap, _ := loadIndex(db, field)
			if ids, exists := indexMap[val]; exists {
				for _, id := range ids {
					candidateIDs[id] = struct{}{}
				}
			}
		}
	}

	var result DocumentList
	for id, doc := range data {
		if len(candidateIDs) > 0 {
			if _, ok := candidateIDs[id]; !ok {
				continue
			}
		}
		if matchDoc(doc, filter) {
			result = append(result, doc)
		}
	}

	// 排序
	if opts != nil && len(opts.Sort) > 0 {
		sort.Slice(result, func(i, j int) bool {
			for k, order := range opts.Sort {
				vi, _ := getNestedValue(result[i], k)
				vj, _ := getNestedValue(result[j], k)
				if compareNumber(vi, vj) == 0 {
					continue
				}
				if order >= 0 {
					return compareNumber(vi, vj) < 0
				} else {
					return compareNumber(vi, vj) > 0
				}
			}
			return true
		})
	}

	// 分页
	if opts != nil && (opts.Skip > 0 || opts.Limit > 0) {
		start := opts.Skip
		end := len(result)
		if opts.Limit > 0 && start+opts.Limit < end {
			end = start + opts.Limit
		}
		result = result[start:end]
	}

	return result, nil
}

func (db *DBContext) FindOne(filter map[string]interface{}) (Document, error) {
	res, err := db.Find(filter, &FindOptions{Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, errors.New("not found")
}

func (db *DBContext) InsertOne(doc Document) (Document, error) {
	JsonMu.Lock()
	defer JsonMu.Unlock()

	data, err := loadCollection(db)
	if err != nil {
		return nil, err
	}

	// 唯一字段校验
	uniqueFields, err := ConfigFile.GetUniqueFields(db.CurrentDB, db.CurrentCollection)
	for _, field := range uniqueFields {
		val, _ := getNestedValue(doc, field)
		for _, d := range data {
			v, _ := getNestedValue(d, field)
			if fmt.Sprintf("%v", val) == fmt.Sprintf("%v", v) {
				return nil, fmt.Errorf("唯一字段冲突: %s", field)
			}
		}
	}

	id := generateObjectID()
	doc["_id"] = id
	data[id] = doc

	if err := saveCollection(db, data); err != nil {
		return nil, err
	}

	// 更新索引
	indexFields := ConfigFile.GetIndexFields(db.CurrentDB, db.CurrentCollection)
	updateIndex(db, id, doc, indexFields, false)

	// 更新文档数量
	_ = db.flashDocCount(db.CurrentCollection)

	return doc, nil
}

func (db *DBContext) InsertMany(docs []Document) ([]Document, error) {
	var result []Document
	for _, doc := range docs {
		newDoc, err := db.InsertOne(doc)
		if err != nil {
			return nil, err
		}
		result = append(result, newDoc)
	}
	return result, nil
}

func (db *DBContext) UpdateOne(filter map[string]interface{}, update Document) (Document, error) {
	updatedDocs, err := db.UpdateMany(filter, update)
	if err != nil {
		return nil, err
	}
	if len(updatedDocs) > 0 {
		return updatedDocs[0], nil
	}
	return nil, errors.New("not found")
}

func (db *DBContext) UpdateMany(filter map[string]interface{}, update Document) ([]Document, error) {
	JsonMu.Lock()
	defer JsonMu.Unlock()

	data, err := loadCollection(db)
	if err != nil {
		return nil, err
	}

	uniqueFields, err := ConfigFile.GetUniqueFields(db.CurrentDB, db.CurrentCollection)
	if err != nil {
		return nil, err
	}
	indexFields := ConfigFile.GetIndexFields(db.CurrentDB, db.CurrentCollection)
	var updated []Document
	for id, doc := range data {
		if matchDoc(doc, filter) {
			// 删除旧索引
			updateIndex(db, id, doc, indexFields, true)

			// 部分更新
			for k, v := range update {
				if k != "_id" {
					doc[k] = v
				}
			}

			// 唯一字段检查
			for _, field := range uniqueFields {
				val, _ := getNestedValue(doc, field)
				for otherID, otherDoc := range data {
					if otherID == id {
						continue
					}
					v, _ := getNestedValue(otherDoc, field)
					if val == v {
						return nil, fmt.Errorf("唯一字段冲突: %s", field)
					}
				}
			}

			data[id] = doc

			// 添加新索引
			updateIndex(db, id, doc, indexFields, false)

			updated = append(updated, doc)
		}
	}

	if err := saveCollection(db, data); err != nil {
		return nil, err
	}

	_ = db.flashDocCount(db.CurrentCollection)

	return updated, nil
}

func (db *DBContext) Delete(filter map[string]interface{}) (int, error) {
	JsonMu.Lock()
	defer JsonMu.Unlock()

	data, err := loadCollection(db)
	if err != nil {
		return 0, err
	}

	indexFields := ConfigFile.GetIndexFields(db.CurrentDB, db.CurrentCollection)

	deleted := 0
	for id, doc := range data {
		if matchDoc(doc, filter) {
			// 删除索引
			updateIndex(db, id, doc, indexFields, true)

			delete(data, id)
			deleted++
		}
	}

	if err := saveCollection(db, data); err != nil {
		return 0, err
	}

	_ = db.flashDocCount(db.CurrentCollection)

	return deleted, nil
}
