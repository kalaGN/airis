# MongoDB 初始化工具

## 功能说明

自动创建 MongoDB 测试数据：
- 创建 16 个集合：`data_20251101_0` 到 `data_20251101_f`
- 每个集合插入 1000 条模拟数据
- 数据结构与项目中 `mongo.go` 保持一致

## 使用方法

### 1. 配置环境变量

确保 `.env` 文件中包含以下配置：

```env
MONGODB_DSN=mongodb://localhost:27017
MONGODB_DATABASE=test_db
```

### 2. 执行初始化

**Linux/macOS:**
```bash
./cmd/init_mongo.sh
```

**Windows:**
```cmd
cmd\init_mongo.bat
```

**直接运行 Go 程序:**
```bash
go run cmd/init_mongo.go
```

## 数据结构

每条文档包含以下字段：

```json
{
  "t": "随机20位字母数字ID",
  "v": "gzip压缩后的数据（6个逗号分隔的随机整数0-1000）"
}
```

解压后的 `v` 字段数据格式：
```
"123,456,789,234,567,890"
```

对应的数据映射（与 `mongo.go` 中的 `varList` 对应）：
- var100001: 第1个值 (索引0)
- var100002: 第2个值 (索引1)
- var100003: 第3个值 (索引2)
- var100004: 第4个值 (索引3)
- var100005: 第5个值 (索引4)
- var100006: 第6个值 (索引5)

## 注意事项

⚠️ **警告**: 执行脚本会先删除已存在的同名集合！

- 确保 MongoDB 服务已启动
- 确保有足够的数据库权限
- 建议在测试环境执行
- 16个集合 × 1000条数据 = 共16,000条文档

## 示例输出

```
Starting MongoDB initialization...
==================================
Connected to MongoDB successfully!
Initializing collection: data_20251101_0
Successfully inserted 1000 documents into data_20251101_0
Initializing collection: data_20251101_1
Successfully inserted 1000 documents into data_20251101_1
...
MongoDB initialization completed!
==================================
MongoDB initialization completed successfully!
Created 16 collections (data_20251101_0 to data_20251101_f)
Each collection contains 1000 documents
```
