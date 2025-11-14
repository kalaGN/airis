package loan

import (
    "fmt"
    "github.com/kalaGN/airis/pkg/env"
    "github.com/kalaGN/airis/pkg/mongo"
    "github.com/kalaGN/airis/pkg/rescode"
    "github.com/kataras/iris/v12"
    "math/rand"
    "strings"
    "time"
)

type CommonRes struct {
	Status int            `json:"status"`
	Msg    string         `json:"msg"`
	Sid    string         `json:"sid"`
	Data   map[string]int `json:"data"`
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func Create(ctx iris.Context) {
	// 使用字母 (a-z, A-Z) 和数字 (0-9) 来生成随机字符串。
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sid strings.Builder
	for i := 0; i < 29; i++ {
		sid.WriteByte(charset[rand.Intn(len(charset))])
	}
	var body struct {
		TestNewFormat interface{} `json:"testNewFormat"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		fmt.Printf("Error reading JSON: %v\n", err)
		ctx.JSON(CommonRes{Status: 1, Msg: "Invalid JSON format", Sid: ""})
		return
	}

	// 类型断言以确保 testNewFormat 是字符串
	testNewFormat, ok := body.TestNewFormat.(string)
	if !ok {
		ctx.JSON(CommonRes{Status: 1, Msg: "testNewFormat must be a string", Sid: ""})
		return
	}
	if testNewFormat == "" {
		ctx.JSON(CommonRes{Status: 1, Msg: "testNewFormat is required", Sid: ""})
		return
	}
	// 准备配置
	config := mongo.Config{
		DSN:   "", // 这里通常是空的，因为实际 DSN 和 DB 会从 env.GetQa 获取
		DB:    "",
		Query: "some_query_value", // 根据实际情况设置查询条件
	}

	config.DSN, config.DB, _, _ = env.GetQa()
	config.Query = "0a62b59dabfc07d58bd3"
	// 调用 GetMongo 方法
    result, err := mongo.GetMongo(ctx, config)
    if err != nil {
        ctx.JSON(CommonRes{Status: rescode.Err1101, Msg: err.Error(), Sid: "", Data: nil})
        return
    }

	fmt.Println("Successfully retrieved and processed data   ")

	// 构建一个 CommonRes 类型的响应体，其中包含一个状态码 "0" 和一个包含随机数的字符串。
    d1 := CommonRes{rescode.SuccessCode, rescode.GetCodeMsg(rescode.SuccessCode), "615" + sid.String(), result}
	// 将响应体以 JSON 格式发送回客户端。
	ctx.JSON(d1)
}
