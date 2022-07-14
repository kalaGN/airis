package main

import (
	"github.com/kalaGN/airis/bootstrap"
	"github.com/kalaGN/airis/pkg/config"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.Default()
	bootstrap.SetupRoute(app)
	port, _ := config.LoadPort()
	app.Listen(":" + port)
}
