/*
 * @Author: afei
 * @Date: 2022-07-14 09:22:45
 * @LastEditors: afei
 * @LastEditTime: 2022-07-14 14:23:48
 * @Description:
 *
 * Copyright (c) 2022 , All Rights Reserved.
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

	// Simple group: ivr
	ivr := app.Party("/ivr")
	{
		ivr.Get("/inter/create", demo)
	}
}

type CommonRes struct {
	Code string `json:"code"`
	Id   string `json:"id"`
}

func demo(ctx iris.Context) {
	d1 := CommonRes{"0", "wf000001"}
	ctx.JSON(d1)
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
