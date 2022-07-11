package main

import (
	"github.com/kalaGN/airis/bootstrap"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.Default()
	bootstrap.SetupRoute(app)
	app.Listen(":8080")
}
