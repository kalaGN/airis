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
	ivrc "github.com/kalaGN/airis/app/Http/controllers/ivr"
	"github.com/kataras/iris/v12"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(app *iris.Application) {

	// Simple group: v1
	v1 := app.Party("/v1")
	{
		v1.Get("/login", ivrc.Create)
	}

	// Simple group: ivr
	ivr := app.Party("/ivr")
	{
		ivr.Get("/inter/create", ivrc.Create)
		ivr.Get("/inter/getflowname", ivrc.Getflowname)

		ivr.Post("/inter/copy", ivrc.Copy)
		ivr.Post("/inter/deleteprocess", ivrc.Del)

	}
}
