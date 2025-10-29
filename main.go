package main

import (
    "context"
    "github.com/kalaGN/airis/bootstrap"
    "github.com/kalaGN/airis/pkg/config"
    "github.com/kataras/iris/v12"
)

func main() {
	app := iris.Default()
	bootstrap.SetupRoute(app)
	port, _ := config.LoadPort()
    iris.RegisterOnInterrupt(func() {
        _ = app.Shutdown(context.Background())
    })
    app.Listen(":"+port, iris.WithoutServerError(iris.ErrServerClosed))
}
