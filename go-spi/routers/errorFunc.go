package routers

import "github.com/kataras/iris/v12"

func notFound(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusNotFound, iris.Map{"success": false, "message": "无效的请求"})
}

func internalServerError(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{"success": false, "message": "请求异常，请稍后重试"})
	// ctx.WriteString("Oups something went wrong, try again")
}

func badRequest(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"success": false, "message": "请求参数有误"})
}

func unAuthorized(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{"success": false, "message": "token无效或已过期"})
}
