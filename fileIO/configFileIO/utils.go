package configFileIO

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StephenChristianW/JsonDB/config"
	UtilsFile "github.com/StephenChristianW/JsonDB/utils/file"
	UtilsTime "github.com/StephenChristianW/JsonDB/utils/time"
	"os"
	"path/filepath"
)

// ==================== 内部工具函数 ====================

// getDB 获取数据库配置，如果不存在返回错误
func getDB(conf *configuration, dbName string) (*dbConfig, error) {
	db, ok := conf.Databases[dbName]
	if !ok {
		return nil, errors.New("数据库: " + dbName + " 不存在")
	}
	return &db, nil
}

// getCollection 获取集合配置，如果不存在返回错误
func getCollection(db *dbConfig, collectionName string) (*collectionConfig, error) {
	col, ok := db.Collections[collectionName]
	if !ok {
		return nil, errors.New("集合: " + collectionName + " 不存在")
	}
	return &col, nil
}

// updateFieldMap 批量设置或取消某个字段 map（UniqueField / Index）
// set=true表示设置字段，false表示取消字段
func updateFieldMap(dbName, collectionName string, fieldNames []string, fieldMapType string, set bool) error {
	conf := getConfig()

	db, err := getDB(conf, dbName)
	if err != nil {
		return err
	}

	col, err := getCollection(db, collectionName)
	if err != nil {
		return err
	}

	var targetMap map[string]struct{}
	switch fieldMapType {
	case "UniqueField":
		if col.Settings.UniqueField == nil {
			col.Settings.UniqueField = make(map[string]struct{})
		}
		targetMap = col.Settings.UniqueField
	case "Index":
		if col.Settings.Index == nil {
			col.Settings.Index = make(map[string]struct{})
		}
		targetMap = col.Settings.Index
	default:
		return errors.New("未知字段类型: " + fieldMapType)
	}

	for _, f := range fieldNames {
		if set {
			if _, exists := targetMap[f]; exists {
				fmt.Println(fieldMapType+" 已存在:", f)
				continue
			}
			targetMap[f] = struct{}{}
		} else {
			if _, exists := targetMap[f]; !exists {
				fmt.Println(fieldMapType+" 不存在:", f)
				continue
			}
			delete(targetMap, f)
		}
	}

	col.UpdateAt = UtilsTime.TimeNow()
	db.Collections[collectionName] = *col
	conf.Databases[dbName] = *db
	return saveConfig(*conf)
}

// GetIndexFields 获取指定集合的索引字段
func GetIndexFields(dbName, collectionName string) []string {
	path := getIndexMetaFilePath(dbName, collectionName)
	if !UtilsFile.IsPathExist(path) {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil || len(bytes) == 0 {
		return nil
	}

	var fields []string
	_ = json.Unmarshal(bytes, &fields)
	return fields
}

// CreateIndex 添加索引字段
func CreateIndex(dbName, collectionName, field string) error {
	fields := GetIndexFields(dbName, collectionName)
	if !contains(fields, field) {
		fields = append(fields, field)
	}
	return saveIndexMeta(dbName, collectionName, fields)
}

// DropIndex 删除索引字段
func DropIndex(dbName, collectionName, field string) error {
	fields := GetIndexFields(dbName, collectionName)
	newFields := make([]string, 0, len(fields))
	for _, f := range fields {
		if f != field {
			newFields = append(newFields, f)
		}
	}
	return saveIndexMeta(dbName, collectionName, newFields)
}

// ---------------- utils ----------------

func getIndexMetaFilePath(dbName, collectionName string) string {
	dir := filepath.Join(config.GetRootDir(), "index")
	_ = os.MkdirAll(dir, 0755)
	fileName := dbName + "_" + collectionName + ".index"
	return filepath.Join(dir, fileName)
}

func saveIndexMeta(dbName, collectionName string, fields []string) error {
	path := getIndexMetaFilePath(dbName, collectionName)
	bytes, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0666)
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
