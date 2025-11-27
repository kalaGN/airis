package utils

import (
	"strings"
	"testing"
)

func TestGenerateSID(t *testing.T) {
	// 测试带前缀的 SID
	sid := GenerateSID("615", 29)
	
	// 验证长度
	expectedLen := len("615") + 29
	if len(sid) != expectedLen {
		t.Errorf("Expected SID length %d, got %d", expectedLen, len(sid))
	}
	
	// 验证前缀
	if !strings.HasPrefix(sid, "615") {
		t.Errorf("Expected SID to start with '615', got %s", sid)
	}
	
	// 验证字符集
	for _, ch := range sid[3:] {
		if !strings.ContainsRune(charset, ch) {
			t.Errorf("Invalid character '%c' in SID", ch)
		}
	}
}

func TestGenerateRandomID(t *testing.T) {
	// 测试无前缀的随机 ID
	id := GenerateRandomID(20)
	
	// 验证长度
	if len(id) != 20 {
		t.Errorf("Expected ID length 20, got %d", len(id))
	}
	
	// 验证字符集
	for _, ch := range id {
		if !strings.ContainsRune(charset, ch) {
			t.Errorf("Invalid character '%c' in ID", ch)
		}
	}
}

func TestGenerateSIDUniqueness(t *testing.T) {
	// 测试生成的 ID 是否唯一
	ids := make(map[string]bool)
	count := 1000
	
	for i := 0; i < count; i++ {
		sid := GenerateSID("615", 29)
		if ids[sid] {
			t.Errorf("Duplicate SID generated: %s", sid)
		}
		ids[sid] = true
	}
	
	if len(ids) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

func BenchmarkGenerateSID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSID("615", 29)
	}
}
