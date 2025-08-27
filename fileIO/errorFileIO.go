package fileIO

import (
	"github.com/StephenChristianW/JsonDB/config"
	UtilsFile "github.com/StephenChristianW/JsonDB/utils/file"
	UtilsTime "github.com/StephenChristianW/JsonDB/utils/time"
	"sync"
)

var errorMu sync.RWMutex

// errorInfo 错误信息结构体
type errorInfo struct {
	ErrorMsg  string `json:"error_msg"`
	ErrorLoc  string `json:"error_loc"`
	ErrorFunc string `json:"error_func"`
	RunTIme   string `json:"run_time"`
}

func initErrorInfo(errMsg string, funcName string, errorLoc string) errorInfo {
	errInfo := errorInfo{
		ErrorMsg:  errMsg,
		ErrorLoc:  errorLoc,
		ErrorFunc: funcName,
		RunTIme:   UtilsTime.TimeNow(),
	}
	return errInfo
}

func WriteErrorInfo(errMsg string, funcName string, errorLoc string) {
	errorMu.Lock()
	defer errorMu.Unlock()
	errInfo := initErrorInfo(errMsg, funcName, errorLoc)
	errorsPath := config.GetErrorFilePath()
	var jsonObj []errorInfo
	if UtilsFile.IsPathExist(errorsPath) {
		err := ReadJsonFile(errorsPath, &jsonObj)
		if err != nil {
			panic(err)
		}
	}
	jsonObj = append(jsonObj, errInfo)
	err := WriteJsonFile(errorsPath, jsonObj)
	if err != nil {
		panic(err)
	}
}
