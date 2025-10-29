package mongo

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/kalaGN/airis/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"strings"
	"time"
)

// Config 包含 MongoDB 连接信息和查询条件
type Config struct {
	DSN   string
	DB    string
	Query string
}

func GetMongo(ctx context.Context, config Config) (map[string]int, error) {
    dsn, db, _, _ := env.GetQa()
    if dsn == "" || db == "" {
        return nil, fmt.Errorf("invalid DSN or DB configuration")
    }
	config.DSN = dsn
	config.DB = db

	// 连接到 MongoDB
	client, err := connectToMongoDB(ctx, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// 获取数据库实例
	database := client.Database(config.DB)
	collection := database.Collection("ads_nginx_log_var_20240827_0")

	query := struct {
		T string `bson:"t"`
	}{
		T: config.Query,
	}

	var foundDoc struct {
		T string `bson:"t"`
		V []byte `bson:"v"`
	}

	err = collection.FindOne(ctx, query).Decode(&foundDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %v", err)
	}

	// 解压 V 字段
	decompressedData, err := gzipDecompress(foundDoc.V)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	// 将解压后的数据转换为字符串
	decompressedUserData := string(decompressedData)

	// 解析逗号分隔的字符串为数组
	userDataArray := strings.Split(decompressedUserData, ",")

	varList := map[string]int{
		"var200001": 0,
		"var200002": 1,
		"var200003": 2,
		"var200004": 3,
		"var200005": 4,
		"var200006": 5,
	}
	result := ProcessUserData(varList, userDataArray)
	return result, nil
}

func connectToMongoDB(ctx context.Context, dsn string) (*mongo.Client, error) {
	if dsn == "" {
		return nil, fmt.Errorf("invalid DSN")
	}

	// 设置连接选项
	clientOptions := options.Client().ApplyURI(dsn)

	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 连接到 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接是否成功
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

func gzipDecompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input data")
	}

	buf := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var out bytes.Buffer
	_, err = out.ReadFrom(gz)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func ProcessUserData(varlist map[string]int, userDataArray []string) map[string]int {
	result := make(map[string]int)

	for key, index := range varlist {
		if index >= 0 && index < len(userDataArray) {
			value, err := strconv.Atoi(userDataArray[index])
			if err != nil {
				// 处理转换错误，可以选择跳过或记录日志
				fmt.Printf("Error converting string to int for key %s: %v\n", key, err)
				continue
			}
			result[key] = value
		}
	}

	return result
}
