@echo off
REM MongoDB 初始化脚本（Windows）
REM 功能：创建 data_20251101_0 到 data_20251101_f 共 16 个集合，每个集合插入 1000 条模拟数据

echo Starting MongoDB initialization...
echo ==================================

REM 检查 .env 文件
if not exist .env (
    echo Warning: .env file not found!
    echo Please create .env file with MONGODB_DSN and MONGODB_DATABASE
    exit /b 1
)

REM 执行初始化
cd /d "%~dp0\.."
go run cmd\init_mongo.go

if %errorlevel% equ 0 (
    echo ==================================
    echo MongoDB initialization completed successfully!
    echo Created 16 collections ^(data_20251101_0 to data_20251101_f^)
    echo Each collection contains 1000 documents
) else (
    echo ==================================
    echo MongoDB initialization failed!
    exit /b 1
)
