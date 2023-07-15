package jwt

import (
	"net/http"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
)

// 返回一个处理请求的中间件处理器函数。在这个函数中可以对请求进行令牌验证和处理
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			//解析token
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt { //token过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{ //401,未授权
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort() //用于终止中间件链的执行，阻止后续的处理器函数继续执行
			//验证JWT失败或者过期时，会立即返回响应，并终止请求链的处理过程
			return
		}

		c.Next()
	}
}
