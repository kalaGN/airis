#!/bin/bash

# MongoDB 初始化脚本
# 功能：创建 data_20251101_0 到 data_20251101_f 共 16 个集合，每个集合插入 1000 条模拟数据

echo "Starting MongoDB initialization..."
echo "=================================="

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "Warning: .env file not found!"
    echo "Please create .env file with MONGODB_DSN and MONGODB_DATABASE"
    exit 1
fi

# 执行初始化
cd "$(dirname "$0")/.."
go run cmd/init_mongo.go

if [ $? -eq 0 ]; then
    echo "=================================="
    echo "MongoDB initialization completed successfully!"
    echo "Created 16 collections (data_20251101_0 to data_20251101_f)"
    echo "Each collection contains 1000 documents"
else
    echo "=================================="
    echo "MongoDB initialization failed!"
    exit 1
fi
