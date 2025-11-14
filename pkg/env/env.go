/*
 * Copyright (c) 2022 , All Rights Reserved.
 */
package env

import (
    "github.com/joho/godotenv"
    "os"
)

// ConfigFunc 动态加载配置信息
type ConfigFunc func() map[string]interface{}

func GetQa() (string, string, string, error) {
	// 加载 .env 文件
	_ = godotenv.Load()
	dsn := os.Getenv("MONGODB_DSN")
	dataBase := os.Getenv("MONGODB_DATABASE")
	timeOut := os.Getenv("MONGODB_TIMEOUT")
	return dsn, dataBase, timeOut, nil
}
