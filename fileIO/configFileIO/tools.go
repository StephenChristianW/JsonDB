package configFileIO

import (
	"encoding/json"
	"github.com/StephenChristianW/JsonDB/config"
	"github.com/StephenChristianW/JsonDB/fileIO"
	UtilsFile "github.com/StephenChristianW/JsonDB/utils/file"
	"os"
	"path/filepath"
)

// ==================== 工具函数 ====================

// FileError 记录错误信息到指定路径

func FileErrorDBIO(msg string, funcName string) {
	fileIO.WriteErrorInfo(msg, dbIOError, funcName)
}
func FileErrorCollectionIO(msg string, funcName string) {
	fileIO.WriteErrorInfo(msg, collectionIOError, funcName)
}
func FileErrorCollectionSettingsIO(msg string, funcName string) {
	fileIO.WriteErrorInfo(msg, collectionSettingsError, funcName)
}

// docCount 返回指定数据库和集合的文档数量
func docCount(dbName, collectionName string) (int, error) {
	colPath := filepath.Join(config.GetRootDir(), dbName, collectionName+".json")
	var documents map[string]interface{}
	if err := fileIO.ReadJsonFile(colPath, &documents); err != nil {
		return 0, err
	}
	if documents == nil {
		return 0, nil
	}
	return len(documents), nil
} // ==================== 配置文件操作 ====================

// initConfig 初始化配置文件路径和默认内容，返回配置文件路径
func initConfig() (string, error) {
	rootDir := config.GetRootDir()

	// 确保根目录存在
	if !UtilsFile.IsPathExist(rootDir) {
		if err := os.MkdirAll(rootDir, 0755); err != nil {
			return "", err
		}
	}

	// 确保配置文件存在，若不存在则创建默认空配置
	configPath := config.GetConfigFilePath()
	if !UtilsFile.IsPathExist(configPath) {
		err := os.WriteFile(configPath, []byte(`{"databases":{}}`), 0644)
		if err != nil {
			return "", err
		}
	}
	return configPath, nil
}

// getConfig 读取并解析配置文件，如果没有配置文件则初始化
func getConfig() *configuration {
	configPath, err := initConfig()
	configMu.RLock()
	defer configMu.RUnlock()

	if err != nil {
		panic(err)
	}

	fileObj, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var conf configuration
	if err := json.Unmarshal(fileObj, &conf); err != nil {
		panic(err)
	}

	if conf.Databases == nil {
		conf.Databases = make(map[string]dbConfig)
	}

	return &conf
}

// saveConfig 保存配置到文件
func saveConfig(conf configuration) error {
	configMu.RLock()
	defer configMu.RUnlock()
	configPath, err := initConfig()
	if err != nil {
		return err
	}

	jsonObj, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, jsonObj, 0644)
}

func GetUniqueFields(dbName string, collectionName string) ([]string, error) {
	// 读取当前配置
	conf := getConfig()

	// 获取指定数据库对象
	db, err := getDB(conf, dbName)
	if err != nil {
		return []string{}, err // 数据库不存在时返回错误
	}
	col, err := getCollection(db, collectionName)

	if err != nil {
		return []string{}, err
	}
	var fields []string
	for field := range col.Settings.UniqueField {
		fields = append(fields, field)
	}
	return fields, nil
}
