//编写JWT工具包，用于生成和解析JWT令牌
//实现身份验证和授权功能

package util

import (
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/dgrijalva/jwt-go"
)

// 用于签名和验证JWT的密钥
var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	jwt.StandardClaims        //标准声明部分
}

// 生成JWT令牌
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog", //发行人
		},
	}

	//jwt.NewWithClaims创建了一个包含自定义声明claims的Token对象
	//jwt.SigningMethodHS256表示使用HS256算法对令牌进行签名
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//将令牌对象使用预先设置的jwtSecret密钥进行签名并生成字符串形式的JWT令牌
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// 解析Token令牌
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	//检查tokenClaims是否成果被解析
	if tokenClaims != nil {
		//检查tokenClaims是否包含有效的声明，并将声明强制转换为*Claims类型的结构体指针claims
		 if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
            return claims, nil
        }
	}

	return nil, err
}
