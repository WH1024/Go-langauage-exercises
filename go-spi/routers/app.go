package routers

import (
	// . "go-spi/middlewares"

	"go-spi/controller"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

// NewIris
// 创建一个路由实例
func NewIris() *iris.Application {
	app := iris.New()
	HandlerError(app)
	return app
}

// RegisterRouter
// 注册路由
func RegisterRouterapp(app *iris.Application) {
	RegisterV1(app)
	app.Get("/tiles/{z}/{x}/{y}/{account}/{key}", controller.Tiles)
}

// func RegisterRouter(app *iris.Application) {
// 	// 查询+发布表单
// 	app.Post("/post", func(ctx iris.Context) {
// 		id, err := ctx.URLParamInt("id")
// 		if err != nil {
// 			ctx.StopWithStatus(iris.StatusBadRequest)
// 			return
// 		}
// 		page := ctx.URLParamInt32Default("page", 0)
// 		name := ctx.PostValue("name")
// 		message := ctx.PostValue("message")
// 		ctx.Writef("id:%d; page:%d; name:%s; message:%s;", id, page, name, message)
// 	})

// 	app.Post("/queryForm", func(ctx iris.Context) {
// 		ids := ctx.URLParamSlice("id")
// 		names, err := ctx.PostValues("names")
// 		if err != nil {
// 			ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "参数有误"})
// 			return
// 		}

// 		data := struct {
// 			Ids   []string `json:"ids"`
// 			Names []string `json:"names"`
// 		}{
// 			ids,
// 			names,
// 		}
// 		ctx.StopWithJSON(iris.StatusOK, iris.Map{"status": 200, "data": data})
// 	})

// 	app.Get("/hi", func(ctx iris.Context) {
// 		ctx.Writef("哈哈哈")
// 	})

// 	app.Get("/ping", func(ctx iris.Context) {
// 		ctx.WriteString("hahapong")
// 		ctx.Host()
// 		ctx.Application().Logger().Infof("Request path: %s, host: %s", ctx.Path(), ctx.Host())
// 	})

// 	// app.Done(func(ctx iris.Context) {
// 	// 	app.Logger().Infof("sent: %s", string(ctx.Recorder().Body()))
// 	// })
// }

// RegisterGlobalMiddleware
// 注册全局中间件
func RegisterGlobalMiddleware(app *iris.Application) {
	// 捕获panic并返回500
	app.Use(recover.New())
	// 注册日志 需做日志分割
	// UseLog(app)
	// 缓存
	// app.Use(iris.Cache304(time.Second * 3))
}

// HandlerError
// iris实例对一些通用的状态码的全局处理
func HandlerError(app *iris.Application) {
	app.OnErrorCode(iris.StatusNotFound, notFound)
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)
	app.OnErrorCode(iris.StatusBadRequest, badRequest)
	app.OnErrorCode(iris.StatusUnauthorized, unAuthorized)
}
