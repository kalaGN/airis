package rescode

import "fmt"

// CodeDefine 定义了一个错误代码及其对应的错误消息。
type CodeDefine struct {
	code int    // code 是错误代码。
	msg  string // msg 是错误消息。
}

// NewCodeDefine 创建一个新的 CodeDefine 实例。
func NewCodeDefine(code int, msg string) *CodeDefine {
	return &CodeDefine{
		code: code,
		msg:  msg,
	}
}

// 定义不同的错误代码及其对应的错误消息。
const (
	SuccessCode = 0
	Err1110     = 1110
	Err1101     = 1101
)

var (
	SuccessMsg = "success"
	Err1110Msg = "invalid parameter cid"
	Err1101Msg = "illegal parameter apikey"
)

// GetCodeMsg 根据错误代码返回对应的错误消息。
func GetCodeMsg(code int) string {
	switch code {
	case SuccessCode:
		return SuccessMsg
	case Err1110:
		return Err1110Msg
	case Err1101:
		return Err1101Msg
	default:
		return fmt.Sprintf("unknown error code: %d", code)
	}
}
