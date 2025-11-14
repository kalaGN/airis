package middleware

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

// Logger 请求日志中间件
func Logger() iris.Handler {
	return func(ctx iris.Context) {
		start := time.Now()

		// 记录请求信息
		logrus.WithFields(logrus.Fields{
			"method": ctx.Method(),
			"path":   ctx.Path(),
			"ip":     ctx.RemoteAddr(),
		}).Info("Request started")

		ctx.Next()

		// 记录响应信息
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"method":   ctx.Method(),
			"path":     ctx.Path(),
			"status":   ctx.GetStatusCode(),
			"duration": duration.Milliseconds(),
		}).Info("Request completed")
	}
}
