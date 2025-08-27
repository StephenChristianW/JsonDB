package UtilsFile

import (
	"encoding/json"
	"errors"
	"os"
)

func ReadJsonFile(filePath string, v interface{}) error {
	if !IsPathExist(filePath) {
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
