package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"      //日志文件拓展名
	TimeFormat  = "20060102" //时间格式
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

// 构建包含时间戳的完整日志文件路径
// runtime/logs/log20230715.log
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// 打开日志文件并返回文件句柄
func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
		case os.IsNotExist(err):
			mkDir()
		case os.IsPermission(err):
			log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir + "/" + getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
