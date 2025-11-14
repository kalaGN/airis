package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
)

// RateLimiter 简单的内存限流器
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每个时间窗口允许的请求数
	window   time.Duration // 时间窗口
}

type Visitor struct {
	lastSeen time.Time
	count    int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// 定期清理过期访客
	go limiter.cleanupVisitors()

	return limiter
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() iris.Handler {
	return func(ctx iris.Context) {
		ip := ctx.RemoteAddr()

		if !rl.allow(ip) {
			ctx.StatusCode(iris.StatusTooManyRequests)
			ctx.JSON(iris.Map{
				"status": 429,
				"msg":    fmt.Sprintf("Rate limit exceeded. Max %d requests per %v", rl.rate, rl.window),
			})
			return
		}

		ctx.Next()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	visitor, exists := rl.visitors[ip]

	if !exists || now.Sub(visitor.lastSeen) > rl.window {
		// 新访客或时间窗口已过期
		rl.visitors[ip] = &Visitor{
			lastSeen: now,
			count:    1,
		}
		return true
	}

	if visitor.count >= rl.rate {
		return false
	}

	visitor.count++
	visitor.lastSeen = now
	return true
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, visitor := range rl.visitors {
			if now.Sub(visitor.lastSeen) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
