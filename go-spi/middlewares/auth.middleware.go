package middlewares

import (
	"fmt"
	"go-spi/db"
	"go-spi/utils"

	"github.com/garyburd/redigo/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
)

// VerifyToken
// 验证token
func VerifyToken(ctx iris.Context) {
	jwtToken := ctx.Values().Get("jwt").(*jwt.Token)
	tokenString, _ := utils.DecodeToken(jwtToken) // 拿到签名
	account := ctx.URLParamDefault("account", "")
	fmt.Println("account: ", account)
	rdPtr := db.GetRedisDbInstance()
	// pika里拿不到该账号的token
	token, err := redis.String(rdPtr.Do("GET", account))
	if err != nil || tokenString != token {
		fmt.Println("err: ", token)
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}
	// 验证通过
	ctx.Next()
}
