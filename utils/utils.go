package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ParseJSON 解析 JSON 文件
func ParseJSON(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MustParseJSON 解析 JSON，失败时 panic
func MustParseJSON(path string) map[string]interface{} {
	data, err := ParseJSON(path)
	if err != nil {
		panic(err)
	}
	return data
}

// SaveJSON 保存数据到 JSON 文件
func SaveJSON(path string, data map[string]interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// PrettyPrintJSON 美观打印 JSON
func PrettyPrintJSON(data map[string]interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("打印 JSON 失败:", err)
		return
	}
	fmt.Println(string(bytes))
}
