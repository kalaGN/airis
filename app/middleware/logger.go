package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

// Logger 请求日志中间件
func Logger() iris.Handler {
	return func(ctx iris.Context) {
		start := time.Now()
		requestTime := start.Format("2006-01-02 15:04:05")

		// 读取请求体
		var requestBody string
		if ctx.Method() == "POST" || ctx.Method() == "PUT" || ctx.Method() == "PATCH" {
			body, err := io.ReadAll(ctx.Request().Body)
			if err == nil {
				requestBody = string(body)
				// 重新设置请求体，供后续处理使用
				ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		// 记录请求开始
		logrus.WithFields(logrus.Fields{
			"time":         requestTime,
			"method":       ctx.Method(),
			"path":         ctx.Path(),
			"ip":           ctx.RemoteAddr(),
			"request_body": requestBody,
		}).Info("Request started")

		// 捕获响应数据
		recorder := ctx.Recorder()
		ctx.Next()

		// 获取响应体
		responseBody := string(recorder.Body())

		// 解析响应状态字段
		var responseStatus interface{}
		var responseData map[string]interface{}
		if err := json.Unmarshal(recorder.Body(), &responseData); err == nil {
			if status, ok := responseData["status"]; ok {
				responseStatus = status
			}
		}

		// 记录请求完成
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"time":            requestTime,
			"method":          ctx.Method(),
			"path":            ctx.Path(),
			"http_status":     ctx.GetStatusCode(),
			"response_status": responseStatus,
			"duration_ms":     duration.Milliseconds(),
			"request_body":    requestBody,
			"response_body":   responseBody,
		}).Info("Request completed")
	}
}
