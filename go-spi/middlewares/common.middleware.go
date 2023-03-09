package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/kataras/iris/v12"
)

// UseLog
// 写入日志文件
func UseLog(app *iris.Application) {
	// app.Logger().SetLevel("debug")
	// app.Logger().Debugf(`Log level set to "debug"`)
	dir, _ := os.Getwd()
	t := time.Now()
	filename := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	logFileName := dir + "/log/" + filename + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if err == nil {
		logger := app.Logger()
		app.ConfigureHost(func(su *iris.Supervisor) {
			su.RegisterOnShutdown(func() {
				// 项目关闭时，关闭文件
				logFile.Close()
			})
		})
		logger.SetLevel("debug")
		logger.SetOutput(logFile)
	}

	// app.Logger().SetOutput(logFile)
	// app.Logger().SetOutput(io.MultiWriter(logFile, os.Stdout))

}
