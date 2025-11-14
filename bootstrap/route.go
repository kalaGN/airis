// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"time"

	"github.com/kalaGN/airis/app/middleware"
	"github.com/kalaGN/airis/routes"
	"github.com/kataras/iris/v12"
)

// SetupRoute 路由初始化
func SetupRoute(router *iris.Application) {

	// 注册全局中间件
	registerGlobalMiddleWare(router)

	//  注册 API 路由
	routes.RegisterAPIRoutes(router)
}

func registerGlobalMiddleWare(router *iris.Application) {
	// 异常恢复中间件（最先执行）
	router.Use(middleware.Recovery())

	// CORS 跨域中间件
	router.Use(middleware.CORS())

	// 请求日志中间件
	router.Use(middleware.Logger())

	// 全局限流：每分钟 100 个请求
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(rateLimiter.RateLimit())
}
