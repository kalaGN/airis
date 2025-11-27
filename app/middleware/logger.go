package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestTime := start.Format("2006-01-02 15:04:05")

		// 读取请求体
		var requestBody string
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(body)
				// 重新设置请求体，供后续处理使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		// 记录请求开始
		logrus.WithFields(logrus.Fields{
			"time":         requestTime,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"ip":           c.ClientIP(),
			"request_body": requestBody,
		}).Info("Request started")

		// 使用自定义 ResponseWriter 捕获响应数据
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 获取响应体
		responseBody := blw.body.String()

		// 解析响应状态字段
		var responseStatus interface{}
		var responseData map[string]interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &responseData); err == nil {
			if status, ok := responseData["status"]; ok {
				responseStatus = status
			}
		}

		// 记录请求完成
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"time":            requestTime,
			"method":          c.Request.Method,
			"path":            c.Request.URL.Path,
			"http_status":     c.Writer.Status(),
			"response_status": responseStatus,
			"duration_ms":     duration.Milliseconds(),
			"request_body":    requestBody,
			"response_body":   responseBody,
		}).Info("Request completed")
	}
}
