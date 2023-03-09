package main

import (
	"context"
	"time"

	. "go-spi/routers"

	"github.com/kataras/iris/v12"
)

// _ "github.com/ClickHouse/clickhouse-go"

// func main() {
// 	dbPtr := db.GetDB()
// 	fmt.Printf("Err: %v", dbPtr.Ping())
// }

func main() {
	app := NewIris()
	// 正常shutdown
	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts.
		app.Shutdown(ctx)
	})

	RegisterGlobalMiddleware(app)
	RegisterRouterapp(app)
	app.Listen(":54087", iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed))
}

// func myMiddleware(ctx iris.Context) {
// 	ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
// 	ctx.Next()
// }
