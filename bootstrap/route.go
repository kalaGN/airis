// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
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

}
