//登陆时进行身份验证并生成令牌
//解决api's可以任意访问的问题

package api

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	//"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")


	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	// data := make(map[string]interface{})
	// code := e.INVALID_PARAMS

	// if ok {
	// 	//表单验证通过后，验证用户名和密码是否有效
	// 	isExist := models.CheckAuth(username, password)
	// 	if isExist {
	// 		//用户名和密码有效，生成token
	// 		token, err := util.GenerateToken(username, password)
	// 		if err != nil {
	// 			code = e.ERROR_AUTH_TOKEN
	// 		} else {
	// 			data["token"] = token
	// 			code = e.SUCCESS
	// 		}
	// 	} else {
	// 		code = e.ERROR_AUTH
	// 	}
	// } else {
	// 	for _, err := range valid.Errors {
	// 		logging.Info(err.Key, err.Message)
	// 	}
	//}

	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": data,
	// })

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	isExist, err := models.CheckAuth(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}

	//客户端在后续的请求中将该Token包含在请求头或请求参数中，以证明自己的身份
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}
