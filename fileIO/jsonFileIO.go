package fileIO

import (
	UtilsFile "JsonDB/utils/file"
	"encoding/json"
	"errors"
	"os"
)

func ReadJsonFile(filePath string, v interface{}) error {
	if !UtilsFile.IsPathExist(filePath) {
		return errors.New("文件不存在")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func WriteJsonFile(filePath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ") // 格式化输出，更易读
	if err != nil {
		return err
	}
	// 不管文件是否存在，都写入
	err = os.WriteFile(filePath, data, 0666)
	return err
}

func CreateDirectory(directoryName string) error {
	if directoryName == "" {
		return errors.New("目录错误")
	}

	// 如果目录已存在，直接返回 nil
	if UtilsFile.IsPathExist(directoryName) {
		return nil
	}

	// 创建目录
	if err := os.MkdirAll(directoryName, 0755); err != nil {
		return err
	}
	return nil
}
