package utils

import (
	"math/rand"
	"strings"
	"time"
)

var (
	// charset 用于生成随机字符串的字符集
	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	// rnd 本地随机数生成器，避免并发竞争
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateSID 生成会话 ID
// prefix: ID 前缀
// length: 随机部分的长度
func GenerateSID(prefix string, length int) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rnd.Intn(len(charset))])
	}
	
	return sb.String()
}

// GenerateRandomID 生成纯随机 ID（无前缀）
func GenerateRandomID(length int) string {
	return GenerateSID("", length)
}
