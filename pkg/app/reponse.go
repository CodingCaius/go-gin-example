package app

import (
	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

//用于统一处理JSON格式响应
type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})

	//return
}