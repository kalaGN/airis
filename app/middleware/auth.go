package middleware

import (
	"github.com/kalaGN/airis/pkg/rescode"
	"github.com/kataras/iris/v12"
)

// AuthRequired 认证中间件
func AuthRequired() iris.Handler {
	return func(ctx iris.Context) {
		// 从请求头获取 token
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{
				"status": 401,
				"msg":    "Authorization token required",
			})
			return
		}

		// TODO: 验证 token 的有效性
		// 这里可以添加 JWT 验证、Redis 验证等逻辑

		// 验证通过，继续处理
		ctx.Next()
	}
}

// APIKeyAuth API Key 认证中间件
func APIKeyAuth() iris.Handler {
	return func(ctx iris.Context) {
		apiKey := ctx.URLParam("apikey")
		if apiKey == "" {
			apiKey = ctx.GetHeader("X-API-Key")
		}

		if apiKey == "" {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{
				"status": rescode.Err1101,
				"msg":    rescode.GetCodeMsg(rescode.Err1101),
			})
			return
		}

		// TODO: 验证 API Key 的有效性
		// 可以从数据库或配置中验证

		ctx.Next()
	}
}
