package configFileIO

import (
	"errors"
	"fmt"
	UtilsTime "github.com/StephenChristianW/JsonDB/utils/time"
)

// ==================== 集合操作 ====================

// CollectionCreateConfig 创建集合配置
//
// 参数：
//
//	dbName - 数据库名称
//	collectionName - 集合名称
//
// 返回值：
//
//	error - 如果集合已存在或创建失败，会返回对应错误；成功返回 nil
func CollectionCreateConfig(dbName, collectionName string) error {

	// 读取当前配置
	conf := getConfig()

	// 获取指定数据库对象
	db, err := getDB(conf, dbName)
	if err != nil {
		return err // 数据库不存在时返回错误
	}

	// 检查集合是否已存在
	if _, ok := db.Collections[collectionName]; ok {
		return errors.New("集合: " + collectionName + " 已存在")
	}

	// 初始化集合配置对象
	col := collectionConfig{
		CreateAt: UtilsTime.TimeNow(),
		UpdateAt: UtilsTime.TimeNow(),
		Settings: collectionSettings{
			UniqueField: nil,
			Index:       nil,
		},
		DocsCount: 0,
	}

	// 将集合对象写入数据库对象
	db.Collections[collectionName] = col
	// 将更新后的数据库对象写回配置对象
	conf.Databases[dbName] = *db

	// 保存配置到文件
	return saveConfig(*conf)
}

// CollectionDeleteConfig 删除集合配置
//
// 参数：
//
//	dbName - 数据库名称
//	collectionName - 集合名称
//
// 返回值：
//
//	error - 如果数据库或集合不存在，返回错误；成功返回 nil
func CollectionDeleteConfig(dbName, collectionName string) error {

	// 读取当前配置
	conf := getConfig()

	// 获取指定数据库对象
	db, err := getDB(conf, dbName)
	if err != nil {
		return err // 数据库不存在时返回错误
	}

	// 检查集合是否存在
	if _, err := getCollection(db, collectionName); err != nil {
		return err
	}

	// 删除集合
	delete(db.Collections, collectionName)

	// 将更新后的数据库对象写回配置对象
	conf.Databases[dbName] = *db

	// 提示输出
	fmt.Println("集合: " + collectionName + " 配置数据已删除")

	// 保存配置到文件
	return saveConfig(*conf)
}

// UpdateCollectionStats 更新集合的统计信息，包括文档数量和更新时间
//
// 参数：
//
//	dbName - 数据库名称
//	collectionName - 集合名称
//
// 返回值：
//
//	error - 如果数据库、集合不存在，或读取文档数量失败，返回对应错误；成功返回 nil
func UpdateCollectionStats(dbName, collectionName string) error {
	// 读取当前配置文件
	conf := getConfig()

	// 获取指定数据库对象
	db, err := getDB(conf, dbName)
	if err != nil {
		return err // 数据库不存在时返回错误
	}

	// 获取指定集合对象
	col, err := getCollection(db, collectionName)
	if err != nil {
		return err // 集合不存在时返回错误
	}

	// 统计集合中文档数量
	count, err := docCount(dbName, collectionName)
	if err != nil {
		return err // 读取文档失败时返回错误
	}

	// 更新集合对象的文档数量和更新时间
	col.DocsCount = count
	col.UpdateAt = UtilsTime.TimeNow()

	// 将更新后的集合对象写回数据库对象
	db.Collections[collectionName] = *col
	// 更新数据库更新时间
	db.UpdateAt = UtilsTime.TimeNow()

	// 将更新后的数据库对象写回配置对象
	conf.Databases[dbName] = *db

	// 保存配置到文件
	return saveConfig(*conf)
}

// CollectionRenameConfig 重命名集合
//
// 参数：
//
//	dbName - 数据库名称
//	collectionName - 旧集合名称
//	newCollectionName - 新集合名称
//
// 返回值：
//
//	error - 如果数据库或旧集合不存在，返回错误；成功返回 nil
func CollectionRenameConfig(dbName, collectionName, newCollectionName string) error {

	// 读取当前配置
	conf := getConfig()

	// 获取指定数据库对象
	db, err := getDB(conf, dbName)
	if err != nil {
		return err // 数据库不存在时返回错误
	}

	// 获取旧集合对象
	col, err := getCollection(db, collectionName)
	if err != nil {
		return err // 集合不存在时返回错误
	}

	// 更新集合更新时间
	col.UpdateAt = UtilsTime.TimeNow()

	// 将集合对象写入新名称
	db.Collections[newCollectionName] = *col
	// 删除旧集合名称
	delete(db.Collections, collectionName)

	// 将更新后的数据库对象写回配置对象
	conf.Databases[dbName] = *db

	// 保存配置到文件
	return saveConfig(*conf)
}
