package services

import (
	"errors"
	"github.com/StephenChristianW/JsonDB/config"
	"github.com/StephenChristianW/JsonDB/fileIO"
	ConfigFile "github.com/StephenChristianW/JsonDB/fileIO/configFileIO"
	UtilsFile "github.com/StephenChristianW/JsonDB/utils/file"
	"os"
	"regexp"
	"strings"
)

const (
	dbServicePath = "JsonDB/services/dbServices.go"
)

// sanitizeName 验证并清理数据库或集合名
func sanitizeName(name string) (string, error) {
	name = strings.TrimSpace(name) // 去掉前后空格
	if name == "" {
		return "", errors.New("名称不能为空")
	}

	// 只保留字母、数字和下划线
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	name = re.ReplaceAllString(name, "")

	if name == "" {
		return "", errors.New("名称中没有合法字符")
	}

	// 不能以数字开头
	if name[0] >= '0' && name[0] <= '9' {
		return "", errors.New("名称不能以数字开头")
	}

	return name, nil
}

// DBService 数据库操作接口
type DBService interface {
	DBCreate(dbName string) error                      // 创建一个新的数据库
	DBSwitch(dbName string) error                      // 切换当前操作的数据库
	DBRename(oldDBName string, newDBName string) error // 重命名数据库
	DBDelete(dbName string) error                      // 删除数据库
	DBList() ([]string, error)                         // 展示数据库列表
}

// writeDBError 统一记录数据库操作错误
// - err: 错误对象
// - funcName: 出错的函数名
// - msg: 额外错误信息
func writeDBError(err error, funcName string, msg string) error {
	var errInfo string
	if msg != "" {
		errInfo = err.Error() + " | msg: " + msg
	} else {
		errInfo = err.Error()
	}
	fileIO.WriteErrorInfo(errInfo, funcName, dbServicePath)
	return err
}

// DBCreate 创建一个新的数据库
// - dbName: 数据库名
func (db *DBContext) DBCreate(dbName string) error {
	name, err := sanitizeName(dbName)
	if err != nil {
		return err
	}
	dbName = name

	funcName := "CreateDB"
	if dbName == "index" {
		return errors.New("数据库名不能为: index")
	}

	// 检查数据库名是否有效
	if err := validateName(dbName); err != nil {
		return writeDBError(err, funcName, dbName)
	}

	// 获取数据库路径
	dbPath, err := db.getDBFilePath(dbName)
	if err != nil {
		return writeDBError(err, funcName, "")
	}

	// 创建数据库目录
	if err = fileIO.CreateDirectory(dbPath); err != nil {
		return writeDBError(err, funcName, "")
	}

	// 在配置文件中创建数据库记录
	if err = ConfigFile.DBCreateConfig(dbName); err != nil {
		return writeDBError(err, funcName, "")
	}

	return nil
}

// DBSwitch 切换当前操作的数据库
// - dbName: 要切换的数据库名
func (db *DBContext) DBSwitch(dbName string) error {
	funcName := "UseDB"

	// 获取数据库路径
	getDbPath, err := db.getDBFilePath(dbName)
	if err != nil {
		return writeDBError(err, funcName, "")
	}

	// 检查数据库目录是否存在
	if !UtilsFile.IsPathExist(getDbPath) {
		return errors.New("数据库: " + dbName + " 不存在")
	}

	// 切换当前数据库上下文
	if err = db.switchDB(dbName); err != nil {
		return writeDBError(err, funcName, "")
	}

	return nil
}

// DBRename 重命名数据库
// - oldDBName: 原数据库名
// - newDBName: 新数据库名
func (db *DBContext) DBRename(oldDBName, newDBName string) error {
	funcName := "RenameDB"

	// 检查旧数据库名
	if oldDBName == "" {
		return writeDBError(errors.New("请输入正确的原数据库名"), funcName, "")
	}

	// 检查新数据库名
	if err := validateName(newDBName); err != nil {
		return writeDBError(errors.New("请输入正确的新数据库名"), funcName, "")
	}

	// 获取旧数据库路径
	oldDbPath, err := db.getDBFilePath(oldDBName)
	if err != nil {
		return writeDBError(err, funcName, "")
	}

	// 获取新数据库路径
	newDbPath, err := db.getDBFilePath(newDBName)
	if err != nil {
		return writeDBError(err, funcName, "")
	}

	// 检查旧数据库是否存在
	if !UtilsFile.IsPathExist(oldDbPath) {
		return writeDBError(errors.New("源数据库不存在: "+oldDBName), funcName, "")
	}

	// 检查新数据库是否已存在
	if UtilsFile.IsPathExist(newDbPath) {
		return writeDBError(errors.New("同名数据库已存在: "+newDBName), funcName, "")
	}

	// 执行目录重命名
	if err = os.Rename(oldDbPath, newDbPath); err != nil {
		return writeDBError(err, funcName, "")
	}

	// 更新配置文件
	if err = ConfigFile.ReNameDBConfig(oldDBName, newDBName); err != nil {
		return writeDBError(err, funcName, "")
	}

	return nil
}

// DBDelete 删除指定数据库及其所有集合，同时更新配置文件
// - dbName: 要删除的数据库名
func (db *DBContext) DBDelete(dbName string) error {
	funcName := "DeleteDB"

	// 检查数据库名是否为空
	if dbName == "" {
		return writeDBError(errors.New("数据库名为空"), funcName, "")
	}

	// 获取数据库目录路径
	filePath, err := db.getDBFilePath(dbName)
	if err != nil {
		return writeDBError(err, funcName, "")
	}

	// 删除数据库目录及其内容
	if err = os.RemoveAll(filePath); err != nil {
		return writeDBError(err, funcName, "")
	}

	// 更新配置文件（删除记录）
	if err = ConfigFile.DBDeleteConfig(dbName); err != nil {
		return writeDBError(err, funcName, "")
	}

	return nil
}

// DBList 展示数据库列表
func (db *DBContext) DBList() ([]string, error) {
	funcName := "DBList"

	// 读取数据库根目录
	dirs, err := os.ReadDir(config.GetRootDir())
	if err != nil {
		return nil, writeDBError(err, funcName, "")
	}

	var dbNames []string
	for _, dir := range dirs {
		if dir.IsDir() {
			if dir.Name() == "index" {
				continue
			}
			dbNames = append(dbNames, dir.Name())
		}
	}

	if len(dbNames) == 0 {
		return nil, writeDBError(errors.New("当前没有数据库"), funcName, "")
	}

	return dbNames, nil
}
