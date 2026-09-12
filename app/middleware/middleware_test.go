package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCORS_HandlesPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.OPTIONS("/resource", func(c *gin.Context) {
		t.Fatal("OPTIONS handler should not run after CORS abort")
	})

	request := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want *", got)
	}
}

func TestRateLimiter_RejectsRequestsOverLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     2,
		window:   time.Minute,
	}
	router := gin.New()
	router.Use(limiter.RateLimit())
	router.GET("/resource", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	for requestNumber := 1; requestNumber <= 3; requestNumber++ {
		request := httptest.NewRequest(http.MethodGet, "/resource", nil)
		request.RemoteAddr = "192.0.2.1:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		want := http.StatusNoContent
		if requestNumber == 3 {
			want = http.StatusTooManyRequests
		}
		if response.Code != want {
			t.Fatalf("request %d HTTP status = %d, want %d", requestNumber, response.Code, want)
		}
	}
}
