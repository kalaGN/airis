package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// VerifyTimestamp 验证时间戳是否在有效期内（默认5分钟）
func VerifyTimestamp(timestamp int64) bool {
	now := time.Now().UnixMilli()
	diff := now - timestamp

	// 允许 5 分钟的时间差
	maxDiff := int64(5 * 60 * 1000) // 5分钟的毫秒数

	return diff >= 0 && diff <= maxDiff
}

// GenerateSign 生成签名
// 签名规则：将所有参数按字母顺序排序，拼接成字符串，加上密钥后做 MD5
func GenerateSign(params map[string]interface{}, secretKey string) string {
	// 提取所有 key 并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除 sign 字段本身
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 按顺序拼接 key=value
	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprintf("%v", params[k]))
		builder.WriteString("&")
	}

	// 加上密钥
	builder.WriteString("key=")
	builder.WriteString(secretKey)

	// MD5 加密
	signStr := builder.String()
	hash := md5.Sum([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// VerifySign 验证签名
func VerifySign(params map[string]interface{}, sign string, secretKey string) bool {
	expectedSign := GenerateSign(params, secretKey)
	return expectedSign == strings.ToUpper(sign)
}
