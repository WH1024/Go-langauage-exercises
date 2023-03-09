package routers

import (
	"go-spi/controller"
	"go-spi/db"
	"go-spi/middlewares"
	"go-spi/utils"

	"github.com/kataras/iris/v12"
)

func RegisterV1(app *iris.Application) {
	v1 := app.Party("/v1")
	// 注册jwt验证
	// v1.Use(utils.GetJwt().Serve)
	{
		v1.Get("/history", controller.MeshQuery)
		v1.Get("/getToken", func(ctx iris.Context) {
			token := utils.GenerateToken("19kefei", "leasekefei")
			rdPtr := db.GetRedisDbInstance()
			rdPtr.Send("SET", "19kefei", token)
			rdPtr.Send("EXPIRE", "19kefei", 10) // 设置存活时间100秒
			rdPtr.Flush()
			ctx.WriteString(token)
		})
		v1.Get("/testToken", utils.GetJwt().Serve, middlewares.VerifyToken, func(ctx iris.Context) {
			ctx.Writef("哈哈哈")
		})
	}
}

// {z}/{x}/{y}
