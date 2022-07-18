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
	"math/rand"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(app *iris.Application) {

	// Simple group: v1
	v1 := app.Party("/v1")
	{
		v1.Get("/login", demo)
	}

	// Simple group: ivr
	ivr := app.Party("/ivr")
	{
		ivr.Get("/inter/create", demo)
		ivr.Get("/inter/getflowname", getflowname)

		ivr.Post("/inter/copy", copy)
		ivr.Post("/inter/deleteprocess", del)

	}
}

type CommonRes struct {
	Code string `json:"code"`
	Id   string `json:"id"`
}

type CommonResg struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func demo(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())

	d1 := CommonRes{"0", "cret0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}

func getflowname(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())
	d1 := CommonResg{"0", "name0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}
func copy(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())
	d1 := CommonRes{"0", "copy0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}
func del(ctx iris.Context) {
	d1 := CommonRes{"200", "ok"}
	ctx.JSON(d1)
}
