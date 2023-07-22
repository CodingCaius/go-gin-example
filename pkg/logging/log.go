//日志库实现
//定义不同的日志级别，并提供相应的日志输入函数

package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// 日志级别,包括 DEBUG、INFO、WARNING、ERROR、FATAL
type Level int

var (
	F *os.File //日志文件

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger //用于实际的日志输出
	logPrefix  = ""          //当前日志前缀
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	// filePath := getLogFileFullPath()
	// F = openLogFile(filePath)
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}

	//创建新的日志记录器 logger，并将其绑定到日志文件 F 上
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
