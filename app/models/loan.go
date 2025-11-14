package models

import "time"

// LoanRequest 贷款请求模型
type LoanRequest struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	SessionID     string    `json:"sid" bson:"sid"`
	TestNewFormat string    `json:"testNewFormat" bson:"test_new_format"`
	Data          map[string]int `json:"data" bson:"data"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// CommonResponse 通用响应结构
type CommonResponse struct {
	Status int            `json:"status"`
	Msg    string         `json:"msg"`
	Sid    string         `json:"sid,omitempty"`
	Data   interface{}    `json:"data,omitempty"`
}

// NewCommonResponse 创建通用响应
func NewCommonResponse(status int, msg string, sid string, data interface{}) *CommonResponse {
	return &CommonResponse{
		Status: status,
		Msg:    msg,
		Sid:    sid,
		Data:   data,
	}
}
