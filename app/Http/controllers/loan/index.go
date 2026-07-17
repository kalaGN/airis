package loan

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	appconfig "github.com/kalaGN/airis/pkg/config"
	"github.com/kalaGN/airis/pkg/env"
	"github.com/kalaGN/airis/pkg/mongo"
	"github.com/kalaGN/airis/pkg/rescode"
	"github.com/kalaGN/airis/pkg/utils"
)

type CommonRes struct {
	Status int            `json:"status"`
	Msg    string         `json:"msg"`
	Sid    string         `json:"sid"`
	Data   map[string]int `json:"data"`
}

func Create(c *gin.Context) {
	// 生成会话 ID
	sid := utils.GenerateSID("100", 29)

	var body struct {
		Phone     interface{} `json:"phone"`
		Pcode     interface{} `json:"pcode"`
		Apikey    string      `json:"apikey"`
		Timestamp interface{} `json:"timestamp"`
		Sign      string      `json:"sign"`
	}

	// 绑定 JSON
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("Error reading JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidJSON, Msg: rescode.GetCodeMsg(rescode.ErrInvalidJSON), Sid: ""})
		return
	}

	// 类型断言以确保 phone 是字符串
	phone, ok := body.Phone.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidPhone, Msg: "phone must be a string", Sid: ""})
		return
	}
	if phone == "" {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrMissingParam, Msg: "phone is required", Sid: ""})
		return
	}

	// 验证 pcode（10001-99999 的数字）
	var pcode int
	switch v := body.Pcode.(type) {
	case float64:
		pcode = int(v)
	case int:
		pcode = v
	default:
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidPcode, Msg: "pcode must be a number", Sid: ""})
		return
	}
	if pcode < 10001 || pcode > 99999 {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidPcode, Msg: rescode.GetCodeMsg(rescode.ErrInvalidPcode), Sid: ""})
		return
	}

	// 验证 apikey
	if body.Apikey == "" {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidApikey, Msg: rescode.GetCodeMsg(rescode.ErrInvalidApikey), Sid: ""})
		return
	}

	// 验证 timestamp（毫秒时间戳）
	var timestamp int64
	switch v := body.Timestamp.(type) {
	case float64:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case int:
		timestamp = int64(v)
	default:
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidTimestamp, Msg: "timestamp must be a number", Sid: ""})
		return
	}
	if !utils.VerifyTimestamp(timestamp) {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidTimestamp, Msg: rescode.GetCodeMsg(rescode.ErrInvalidTimestamp), Sid: ""})
		return
	}

	// 验证签名
	if body.Sign == "" {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidSign, Msg: "sign is required", Sid: ""})
		return
	}

	// 构建签名参数
	signParams := map[string]interface{}{
		"phone":     phone,
		"pcode":     pcode,
		"apikey":    body.Apikey,
		"timestamp": timestamp,
	}

	if !utils.VerifySign(signParams, body.Sign, appconfig.GetSecretKey()) {
		c.JSON(http.StatusBadRequest, CommonRes{Status: rescode.ErrInvalidSign, Msg: rescode.GetCodeMsg(rescode.ErrInvalidSign), Sid: ""})
		return
	}

	// 准备配置
	config := mongo.Config{
		DSN:   "", // 这里通常是空的，因为实际 DSN 和 DB 会从 env.GetQa 获取
		DB:    "",
		Query: "some_query_value", // 根据实际情况设置查询条件
	}

	config.DSN, config.DB, config.Collection, _, _ = env.GetQa()
	// 使用 phone 作为查询条件，这样可以灵活查询
	config.Query = phone

	// 调用 GetMongo 方法，使用 Gin 的 Request Context
	result, err := mongo.GetMongo(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusOK, CommonRes{Status: rescode.ErrDataNotFound, Msg: err.Error(), Sid: "", Data: nil})
		return
	}

	fmt.Println("Successfully retrieved and processed data   ")

	// 构建响应体
	d1 := CommonRes{rescode.SuccessCode, rescode.GetCodeMsg(rescode.SuccessCode), sid, result}
	// 将响应体以 JSON 格式发送回客户端。
	c.JSON(http.StatusOK, d1)
}
