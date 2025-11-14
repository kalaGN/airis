package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataDocument MongoDB 文档结构
type DataDocument struct {
	T string `bson:"t"` // 标识符
	V []byte `bson:"v"` // gzip 压缩后的数据
}

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 获取 MongoDB 配置
	dsn := os.Getenv("MONGODB_DSN")
	database := os.Getenv("MONGODB_DATABASE")

	if dsn == "" {
		log.Fatal("MONGODB_DSN environment variable is required")
	}
	if database == "" {
		log.Fatal("MONGODB_DATABASE environment variable is required")
	}

	// 连接 MongoDB
	ctx := context.Background()
	client, err := connectMongoDB(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(database)

	// 初始化集合 data_20251101_0 到 data_20251101_f (0-15，十六进制)
	collections := []string{
		"data_20251101_0", "data_20251101_1", "data_20251101_2", "data_20251101_3",
		"data_20251101_4", "data_20251101_5", "data_20251101_6", "data_20251101_7",
		"data_20251101_8", "data_20251101_9", "data_20251101_a", "data_20251101_b",
		"data_20251101_c", "data_20251101_d", "data_20251101_e", "data_20251101_f",
	}

	rand.Seed(time.Now().UnixNano())

	for _, collName := range collections {
		log.Printf("Initializing collection: %s", collName)

		collection := db.Collection(collName)

		// 删除已存在的集合（可选）
		collection.Drop(ctx)

		// 生成并插入 1000 条数据
		documents := make([]interface{}, 1000)
		for i := 0; i < 1000; i++ {
			doc, err := generateDocument(i)
			if err != nil {
				log.Printf("Failed to generate document %d: %v", i, err)
				continue
			}
			documents[i] = doc
		}

		// 批量插入
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			log.Printf("Failed to insert documents into %s: %v", collName, err)
			continue
		}

		log.Printf("Successfully inserted 1000 documents into %s", collName)
	}

	log.Println("MongoDB initialization completed!")
}

// connectMongoDB 连接到 MongoDB
func connectMongoDB(ctx context.Context, dsn string) (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(dsn).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	return client, nil
}

// generateDocument 生成单条文档数据
func generateDocument(index int) (DataDocument, error) {
	// 生成随机标识符 (20位字母数字)
	t := generateRandomID(20)

	// 生成随机数据值 (6个字段)
	dataValues := make([]string, 6)
	for i := 0; i < 6; i++ {
		// 生成 0-1000 之间的随机整数
		dataValues[i] = strconv.Itoa(rand.Intn(1001))
	}

	// 将数据拼接成逗号分隔的字符串
	dataString := strings.Join(dataValues, ",")

	// 压缩数据
	compressedData, err := gzipCompress([]byte(dataString))
	if err != nil {
		return DataDocument{}, fmt.Errorf("failed to compress data: %v", err)
	}

	return DataDocument{
		T: t,
		V: compressedData,
	}, nil
}

// generateRandomID 生成随机 ID
func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// gzipCompress 压缩数据
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
