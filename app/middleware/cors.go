package middleware

import (
	"github.com/kataras/iris/v12"
)

// CORS 跨域中间件
func CORS() iris.Handler {
	return func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Header("Access-Control-Max-Age", "3600")

		if ctx.Method() == iris.MethodOptions {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
