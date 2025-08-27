package config

import (
	"os"
	"path/filepath"
	"strings"
)

func GetErrorFilePath() string {
	return filepath.Join(GetRootDir(), ".errors")
}
func GetConfigFilePath() string {
	return filepath.Join(GetRootDir(), ".config")
}

func GetRootDir() string {
	// 获取可执行文件绝对路径
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	//fmt.Println("程序目录:", exeDir)

	var baseDir string
	lowerDir := strings.ToLower(exeDir)

	if strings.Contains(lowerDir, "tmp") || strings.Contains(lowerDir, "goland") {
		// 在 Goland 临时目录 → 使用当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		baseDir = wd
		//fmt.Println("程序在临时目录运行，使用工作目录:", baseDir)
	} else {
		baseDir = exeDir
	}

	// 拼接 JsonDataBase
	dir := filepath.Join(baseDir, "JsonDataBase")

	// 确保目录存在
	if err = os.MkdirAll(dir, 0777); err != nil {
		panic(err)
	}
	return dir
}
