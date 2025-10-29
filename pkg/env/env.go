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

func GetQa() (string, string, string, error) { // 加载 .env 文件
    _ = godotenv.Load()
    qaDsn := os.Getenv("MONGODB_QA1_DSN")
    dataBase := os.Getenv("MONGODB_QA1_DATABASE")
    timeOut := os.Getenv("MONGODB_QA1_TIMEOUT")
    return qaDsn, dataBase, timeOut, nil
}
