package app

import (
	"github.com/astaxie/beego/validation"

	"github.com/EDDYCJY/go-gin-example/pkg/logging"
)

//用于处理验证错误，并将错误信息记录到日志中
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		//对每个验证错误进行日志记录
		logging.Info(err.Key, err.Message)
	}

	//return
}