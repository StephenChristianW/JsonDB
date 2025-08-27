package UtilsFile

import "os"

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// 文件或目录存在
		return true
	}
	if os.IsNotExist(err) {
		// 文件或目录不存在
		return false
	}
	// 其他错误（比如权限问题）也返回 false 或根据需求处理
	return false
}
