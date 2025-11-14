package middleware

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

// Recovery 异常恢复中间件
func Recovery() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"path":  ctx.Path(),
				}).Error("Panic recovered")

				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"status": 500,
					"msg":    fmt.Sprintf("Internal server error: %v", err),
				})
			}
		}()
		ctx.Next()
	}
}
