package services

import (
	"JsonDB/config"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ==================== 数据库上下文操作 ====================

// DBContextInterface 定义数据库上下文操作接口
type DBContextInterface interface {
	GetDBFilePath(dbName string) (string, error)  // 获取数据库文件路径
	SwitchDB(dbName string) error                 // 切换当前数据库
	SwitchCollection(collectionName string) error // 切换当前集合
	GetCollectionFilePath() (string, error)       // 获取当前集合文件路径
}

// DBContext 数据库上下文结构体，记录当前数据库和集合
type DBContext struct {
	CurrentDB         string // 当前选中的数据库
	CurrentCollection string // 当前选中的集合
}

// ==================== 数据库路径相关 ====================

// getDBFilePath 获取指定数据库的路径
//
// 参数：
//
//	dbName - 数据库名称
//
// 返回值：
//
//	string - 数据库路径
//	error - 数据库名为空时返回错误
func (db *DBContext) getDBFilePath(dbName string) (string, error) {
	if dbName == "" {
		return "", errors.New("请输入正确数据库名")
	}
	// 如果当前上下文未指定数据库，则默认使用传入的数据库
	if db.CurrentDB == "" {
		db.CurrentDB = dbName
	}
	return filepath.Join(config.GetRootDir(), dbName), nil
}

// ==================== 数据库切换 ====================

// switchDB 切换当前上下文的数据库
//
// 参数：
//
//	dbName - 数据库名称
//
// 返回值：
//
//	error - 数据库名为空时返回错误
func (db *DBContext) switchDB(dbName string) error {
	if dbName == "" {
		return errors.New("请输入正确数据库名")
	}
	db.CurrentDB = dbName
	return nil
}

// ==================== 集合切换 ====================

// switchCollection 切换当前上下文的集合
//
// 参数：
//
//	collectionName - 集合名称
//
// 返回值：
//
//	error - 当前未选择数据库或集合名为空时返回错误
func (db *DBContext) switchCollection(collectionName string) error {
	if db.CurrentDB == "" {
		return errors.New("请输入正确数据库名")
	}
	if collectionName == "" {
		return errors.New("请输入正确集合名")
	}
	db.CurrentCollection = collectionName
	return nil
}

// getCollectionFilePath 获取当前集合的文件路径
//
// 参数：
//
//	collectionName - 集合名称
//
// 返回值：
//
//	string - 集合文件路径
//	error - 数据库未选择或集合名为空时返回错误
func (db *DBContext) getCollectionFilePath(collectionName string) (string, error) {
	if db.CurrentDB == "" {
		return "", errors.New("未选择数据库")
	}
	if collectionName == "" {
		return "", errors.New("未选择集合")
	}
	path := filepath.Join(config.GetRootDir(), db.CurrentDB, collectionName+".json")
	return path, nil
}

// ==================== 数据库列表 ====================

// getDBs 获取所有数据库名称
//
// 返回值：
//
//	map[string]struct{} - 数据库名称集合
//	error - 根目录读取失败或无数据库时返回错误
func getDBs() (map[string]struct{}, error) {
	dbNames, err := os.ReadDir(config.GetRootDir())
	if err != nil {
		return nil, err
	}
	if len(dbNames) == 0 {
		return nil, errors.New("无数据库")
	}

	var dbs = make(map[string]struct{})
	for _, v := range dbNames {
		dbs[v.Name()] = struct{}{}
	}
	return dbs, nil
}

// ==================== 集合列表 ====================

// getCollectionNames 获取指定数据库中的集合名称列表
//
// 参数：
//
//	dbName - 数据库名称
//
// 返回值：
//
//	[]string - 集合名称切片
//	error - 数据库不存在或无集合时返回错误
func (db *DBContext) getCollectionNames(dbName string) ([]string, error) {
	dbPath, err := db.getDBFilePath(dbName)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	nameSlice := make([]string, 0)
	for _, file := range files {
		colName := strings.TrimSuffix(file.Name(), ".json") // 去掉 .json 后缀
		nameSlice = append(nameSlice, colName)
	}

	if len(nameSlice) == 0 {
		return nil, errors.New("数据库: " + dbName + "中无集合")
	}
	return nameSlice, nil
}

// ==================== 命名校验 ====================

// validateName 校验数据库或集合名称是否合法
//
// 参数：
//
//	name - 待校验名称
//
// 返回值：
//
//	error - 名称为空或包含非法字符时返回错误，合法返回 nil
func validateName(name string) error {
	if name == "" {
		return errors.New("命名不能为空")
	}

	// 只允许字母、数字和下划线
	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validName.MatchString(name) {
		return errors.New("命名只能包含字母、数字和下划线")
	}
	return nil
}
