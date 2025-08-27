package services

import (
	"errors"
	"fmt"
	"github.com/StephenChristianW/JsonDB/fileIO"
	ConfigFile "github.com/StephenChristianW/JsonDB/fileIO/configFileIO"
	UtilsFile "github.com/StephenChristianW/JsonDB/utils/file"
	"os"
)

// CollectionService 集合操作接口
type CollectionService interface {
	CollectionSwitch(collectionName string) error                       // 切换当前集合
	CollectionList(dbName string) ([]string, error)                     // 列出指定数据库下的所有集合
	CollectionCreate(collectionName string) error                       // 在当前数据库中创建集合
	CollectionDelete(collectionName string) error                       // 删除指定集合
	CollectionRename(oldCollectionName, newCollectionName string) error // 重命名集合
}

// 常量定义
const (
	collectionServicePath = "JsonDB/services/collectionServices.go"
)

// ========== 公共方法 ==========

// flashDocCount 更新指定集合的文档数量
// - collectionName: 集合名
func (db *DBContext) flashDocCount(collectionName string) error {
	err := ConfigFile.UpdateCollectionStats(db.CurrentDB, collectionName)
	if err != nil {
		return err
	}
	return nil
}

// writeCollectionError 统一记录集合操作错误
// - err: 错误对象
// - funcName: 出错的函数名
// - msg: 额外错误信息
func writeCollectionError(err error, funcName string, msg string) error {
	var errInfo string
	if msg != "" {
		errInfo = err.Error() + " | msg: " + msg
	} else {
		errInfo = err.Error()
	}
	fileIO.WriteErrorInfo(errInfo, funcName, collectionServicePath)
	return err
}

// ========== 接口实现 ==========

// CollectionSwitch 切换当前集合
// - collectionName: 要切换的集合名
func (db *DBContext) CollectionSwitch(collectionName string) error {
	// 切换集合
	err := db.switchCollection(collectionName)
	if err != nil {
		return writeCollectionError(err, "CollectionSwitch", collectionName)
	}

	// 更新文档数量
	if err := db.flashDocCount(collectionName); err != nil {
		return writeCollectionError(err, "CollectionSwitch", collectionName)
	}
	return nil
}

// CollectionList 显示指定数据库的所有集合
// - dbName: 数据库名
func (db *DBContext) CollectionList(dbName string) ([]string, error) {
	funcName := "CollectionList"

	if dbName == "" {
		return nil, writeCollectionError(errors.New("请输入正确的数据库名称"), funcName, dbName)
	}

	dbs, _ := getDBs()
	if _, ok := dbs[dbName]; !ok {
		return nil, writeCollectionError(errors.New("数据库不存在"), funcName, dbName)
	}

	collections, err := db.getCollectionNames(dbName)
	if collections == nil || len(collections) == 0 {
		return nil, writeCollectionError(errors.New("数据库: "+dbName+" 内无集合"), funcName, dbName)
	}

	if err != nil {
		return nil, writeCollectionError(err, funcName, dbName)
	}

	return collections, nil
}

// CollectionCreate 在当前数据库中创建新的集合
// - collectionName: 新集合名
func (db *DBContext) CollectionCreate(collectionName string) error {
	funcName := "CollectionCreate"

	// 检查是否选择数据库
	if db.CurrentDB == "" {
		return errors.New("请选择数据库")
	}

	// 验证集合名称
	if err := validateName(collectionName); err != nil {
		return errors.New("请输入正确集合名称")
	}

	// 获取集合文件路径
	colPath, err := db.getCollectionFilePath(collectionName)
	if err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}

	// 检查是否已存在
	if UtilsFile.IsPathExist(colPath) {
		return writeCollectionError(errors.New("集合: "+collectionName+" 已存在于: "+db.CurrentDB+" 中"), funcName, collectionName)
	}

	// 创建空 JSON 文件作为集合
	if err = os.WriteFile(colPath, []byte("{}"), 0666); err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}
	fmt.Printf("集合: %s.%s 已创建 \n", db.CurrentDB, collectionName)

	// 更新配置文件
	if err = ConfigFile.CollectionCreateConfig(db.CurrentDB, collectionName); err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}
	if err = ConfigFile.DBUpdateConfig(db.CurrentDB); err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}

	// 更新文档数量
	if err := db.flashDocCount(collectionName); err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}

	return nil
}

// CollectionDelete 删除指定集合
// - collectionName: 要删除的集合名
func (db *DBContext) CollectionDelete(collectionName string) error {
	funcName := "CollectionDelete"

	// 检查数据库是否选择
	if db.CurrentDB == "" {
		return errors.New("未选择数据库")
	}

	// 获取集合路径
	colPath, err := db.getCollectionFilePath(collectionName)
	if err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}

	// 如果集合存在则删除
	if UtilsFile.IsPathExist(colPath) {
		// 删除集合文件
		if err := os.Remove(colPath); err != nil {
			return writeCollectionError(err, funcName, collectionName)
		}
		// 更新配置文件
		if err = ConfigFile.CollectionDeleteConfig(collectionName, colPath); err != nil {
			return writeCollectionError(err, funcName, collectionName)
		}
		if err = ConfigFile.DBUpdateConfig(db.CurrentDB); err != nil {
			return writeCollectionError(err, funcName, collectionName)
		}
		fmt.Printf("集合: %s.%s 已删除 \n", db.CurrentDB, collectionName)
	} else {
		fmt.Printf("未找到: %s.%s 集合 \n", db.CurrentDB, collectionName)
	}

	// 更新文档数量
	if err = db.flashDocCount(collectionName); err != nil {
		return writeCollectionError(err, funcName, collectionName)
	}
	return nil
}

// CollectionRename 重命名集合
// - oldCollectionName: 原集合名
// - newCollectionName: 新集合名
func (db *DBContext) CollectionRename(oldCollectionName, newCollectionName string) error {
	funcName := "CollectionRename"

	// 检查数据库是否选择
	if db.CurrentDB == "" {
		return errors.New("未选择数据库")
	}

	// 获取集合路径
	oldColPath, err := db.getCollectionFilePath(oldCollectionName)
	if err != nil {
		return writeCollectionError(err, funcName, oldCollectionName)
	}
	newColPath, err := db.getCollectionFilePath(newCollectionName)
	if err != nil {
		return writeCollectionError(err, funcName, newCollectionName)
	}

	// 执行重命名
	if err = os.Rename(oldColPath, newColPath); err != nil {
		return writeCollectionError(err, funcName, oldCollectionName+"->"+newCollectionName)
	}

	// 更新配置文件
	if err = ConfigFile.CollectionRenameConfig(db.CurrentDB, oldCollectionName, newCollectionName); err != nil {
		return writeCollectionError(err, funcName, oldCollectionName+"->"+newCollectionName)
	}
	if err = ConfigFile.DBUpdateConfig(db.CurrentDB); err != nil {
		return writeCollectionError(err, funcName, oldCollectionName+"->"+newCollectionName)
	}

	fmt.Printf("集合: %s.%s 已改名为: %s.%s \n", db.CurrentDB, oldCollectionName, db.CurrentDB, newCollectionName)

	// 更新文档数量
	if err := db.flashDocCount(newCollectionName); err != nil {
		return writeCollectionError(err, funcName, newCollectionName)
	}

	return nil
}
