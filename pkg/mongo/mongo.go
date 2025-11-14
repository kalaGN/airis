package mongo

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"sync"
	"github.com/kalaGN/airis/pkg/config"
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
	DSN        string
	DB         string
	Collection string // 集合名称
	Query      string
}

var (
	client *mongo.Client
	once   sync.Once
	clientErr error
)

// GetClient 获取全局 MongoDB 客户端（单例模式 + 连接池）
func GetClient(ctx context.Context) (*mongo.Client, error) {
	once.Do(func() {
		cfg := config.GetMongoConfig()
		if cfg.DSN == "" {
			clientErr = fmt.Errorf("MongoDB DSN is empty")
			return
		}

		// 配置连接池参数
		clientOptions := options.Client().
			ApplyURI(cfg.DSN).
			SetMaxPoolSize(uint64(cfg.MaxPool)).
			SetMinPoolSize(uint64(cfg.MinPool)).
			SetMaxConnIdleTime(30 * time.Second).
			SetConnectTimeout(5 * time.Second).
			SetServerSelectionTimeout(5 * time.Second)

		// 连接 MongoDB
		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientErr = fmt.Errorf("failed to connect to MongoDB: %v", err)
			return
		}

		// 测试连接
		err = client.Ping(ctx, nil)
		if err != nil {
			clientErr = fmt.Errorf("failed to ping MongoDB: %v", err)
			return
		}

		log.Println("MongoDB connected with connection pool")
	})

	return client, clientErr
}

// Close 关闭 MongoDB 连接
func Close(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}

func GetMongo(ctx context.Context, config Config) (map[string]int, error) {
	dsn, db, collectionName, _, _ := env.GetQa()
	if dsn == "" || db == "" {
		return nil, fmt.Errorf("invalid DSN or DB configuration")
	}
	config.DSN = dsn
	config.DB = db
	if config.Collection == "" {
		config.Collection = collectionName
	}

	// 使用连接池客户端
	client, err := GetClient(ctx)
	if err != nil {
		return nil, err
	}

	// 获取数据库实例
	database := client.Database(config.DB)
	
	// 使用配置的集合名称，如果未指定则使用默认值
	colName := config.Collection
	if colName == "" {
		colName = "data_20251101_0"
	}
	collection := database.Collection(colName)

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
		"var100001": 0,
		"var100002": 1,
		"var100003": 2,
		"var100004": 3,
		"var100005": 4,
		"var100006": 5,
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
