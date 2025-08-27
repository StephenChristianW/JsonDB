package configFileIO

import (
	UtilsTime "JsonDB/utils/time"
	"errors"
	"fmt"
)

// ==================== 数据库操作 ====================

// DBCreateConfig 创建数据库配置
func DBCreateConfig(dbName string) error {
	conf := getConfig()
	if _, ok := conf.Databases[dbName]; ok {
		return errors.New("数据库: " + dbName + " 已存在")
	}

	conf.Databases[dbName] = dbConfig{
		DBName:      dbName,
		CreateAt:    UtilsTime.TimeNow(),
		UpdateAt:    UtilsTime.TimeNow(),
		Collections: make(map[string]collectionConfig),
	}

	return saveConfig(*conf)
}

// DBUpdateConfig 更新数据库更新时间
func DBUpdateConfig(dbName string) error {

	conf := getConfig()
	db, err := getDB(conf, dbName)
	if err != nil {
		return err
	}

	db.UpdateAt = UtilsTime.TimeNow()
	conf.Databases[dbName] = *db
	return saveConfig(*conf)
}

// DBDeleteConfig 删除数据库配置
func DBDeleteConfig(dbName string) error {
	conf := getConfig()
	if _, err := getDB(conf, dbName); err != nil {
		return err
	}

	delete(conf.Databases, dbName)
	fmt.Println("数据库: " + dbName + " 配置数据已删除")
	return saveConfig(*conf)
}

// ReNameDBConfig 重命名数据库
func ReNameDBConfig(dbName, newDBName string) error {
	conf := getConfig()
	db, err := getDB(conf, dbName)
	if err != nil {
		return err
	}

	db.DBName = newDBName
	db.UpdateAt = UtilsTime.TimeNow()
	delete(conf.Databases, dbName)
	conf.Databases[newDBName] = *db

	return saveConfig(*conf)
}
