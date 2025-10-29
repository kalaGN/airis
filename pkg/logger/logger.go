// Package logger 日志工具包
package logger

import (
	"fmt"
	"os"
)

// ErrorString 打印错误信息到标准错误输出
func ErrorString(module, function, message string) {
	_, _ = fmt.Fprintf(os.Stderr, "Error in %s.%s: %s\n", module, function, message)
}

// LogIf 如果 err 不为 nil，则打印错误信息到标准错误输出
func LogIf(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
}
