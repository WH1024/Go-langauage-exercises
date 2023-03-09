package utils

import (
	"fmt"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var (
	signer = []byte("soa&r*ab(ili)ty@#") // 密钥
)

func GetJwt() *jwt.Middleware {
	return jwt.New(jwt.Config{
		ErrorHandler: func(ctx iris.Context, err error) {
			if err == nil {
				return
			}
			fmt.Println("zheli 出错了", err.Error())
			ctx.StopWithStatus(iris.StatusUnauthorized)
		},
		// 从请求参数token中提取
		Extractor: jwt.FromParameter("token"),
		// 这里从请求头的“Authorization”的Bearer {token}拿
		// Extractor: jwt.FromAuthHeader,
		// 设置一个函数返回秘钥，关键在于return []byte("这里设置秘钥")
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signer, nil
		},
		// 签名方式
		SigningMethod: jwt.SigningMethodHS256,
	})
}

// GenerateToken
// 生成token
func GenerateToken(account, password string) string {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 签发时间
		"iat": time.Now().Unix(),
		// 有效期
		"exp": time.Now().Add(time.Second * time.Duration(10)).Unix(),
		// 自定义字段
		"account":  account,
		"password": password,
	})
	tokenStr, _ := token.SignedString(signer)
	return tokenStr
}

func DecodeToken(token *jwt.Token) (interface{}, error) {
	return token.SignedString(signer)
}
