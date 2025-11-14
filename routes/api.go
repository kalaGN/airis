/*
 * @Author: afei
 * @Date: 2022-07-14 09:22:45
 * @LastEditors: afei
 * @LastEditTime: 2022-07-18 10:55:38
 * @Description:
 *
 * Copyright (c) 2022 , All Rights Reserved.
 */
package routes

import (
	loanc "github.com/kalaGN/airis/app/Http/controllers/loan"
	"github.com/kataras/iris/v12"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(app *iris.Application) {
	app.Get("/health", func(ctx iris.Context) { ctx.WriteString("ok") })

	loan := app.Party("/loan")
	{
		loan.Post("/", loanc.Create)
	}
}
