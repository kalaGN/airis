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
