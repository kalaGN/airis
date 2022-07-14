/*
 * @Author: afei
 * @Date: 2022-07-14 09:22:45
 * @LastEditors: afei
 * @LastEditTime: 2022-07-14 10:30:22
 * @Description:
 *
 * Copyright (c) 2022 by Infobird, All Rights Reserved.
 */
package routes

import (
	"github.com/kataras/iris/v12"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(app *iris.Application) {

	// Simple group: v1
	v1 := app.Party("/v1")
	{
		v1.Get("/login", list)
	}
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func list(ctx iris.Context) {
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}
