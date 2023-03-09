package controller

import "github.com/kataras/iris/v12"

func Test(ctx iris.Context) {
	ctx.GetReferrer()
	ctx.Request().Referer()
}
