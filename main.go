package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kalaGN/airis/bootstrap"
	"github.com/kalaGN/airis/pkg/config"
	"github.com/kalaGN/airis/pkg/logger"
	"github.com/kalaGN/airis/pkg/mongo"
	"github.com/kataras/iris/v12"
)

func main() {
	// 加载环境变量
	_ = godotenv.Load()

	// 初始化日志系统
	logger.Init()
	logger.Log.Info("Starting application...")

	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Log.Info("Configuration loaded")

	// 初始化 Iris 应用
	app := iris.Default()
	bootstrap.SetupRoute(app)

	// 获取服务端口
	port, _ := config.LoadPort()

	// 注册优雅关闭
	iris.RegisterOnInterrupt(func() {
		logger.Log.Info("Shutting down gracefully...")
		
		// 关闭 MongoDB 连接
		if err := mongo.Close(context.Background()); err != nil {
			logger.Log.Errorf("Error closing MongoDB: %v", err)
		}
		
		_ = app.Shutdown(context.Background())
		logger.Log.Info("Application stopped")
	})

	logger.Log.Infof("Server starting on port %s", port)
	app.Listen(":"+port, iris.WithoutServerError(iris.ErrServerClosed))
}
