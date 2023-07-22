package logging

import (
	"fmt"
	//"log"
	"os"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"      //日志文件拓展名
	TimeFormat  = "20060102" //时间格式
)

func getLogFilePath() string {
	// return fmt.Sprintf("%s", LogSavePath)

	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// // 构建包含时间戳的完整日志文件路径
// // runtime/logs/log20230715.log
// func getLogFileFullPath() string {
// 	prefixPath := getLogFilePath()
// 	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

// 	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
// }

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

// // 打开日志文件并返回文件句柄
// func openLogFile(filePath string) *os.File {
// 	_, err := os.Stat(filePath)
// 	switch {
// 		case os.IsNotExist(err):
// 			mkDir()
// 		case os.IsPermission(err):
// 			log.Fatalf("Permission :%v", err)
// 	}

// 	handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatalf("Fail to OpenFile :%v", err)
// 	}

// 	return handle
// }

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
