package helloworld

import (
	"testing"
)

func TestHelloWorldProto(t *testing.T) {
	// 测试 HelloRequest
	req := &HelloRequest{Name: "Airis"}
	if req.GetName() != "Airis" {
		t.Errorf("Expected name 'Airis', got '%s'", req.GetName())
	}

	// 测试 HelloReply
	rep := &HelloReply{Message: "Hello Airis"}
	if rep.GetMessage() != "Hello Airis" {
		t.Errorf("Expected message 'Hello Airis', got '%s'", rep.GetMessage())
	}

	t.Log("✅ gRPC protobuf 代码生成成功！")
	t.Log("✅ HelloRequest 和 HelloReply 类型定义正常！")
}
