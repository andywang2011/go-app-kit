package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func rateLimitRouter(t *testing.T, rate int64) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	h, err := RateLimit(rate)
	if err != nil {
		t.Fatalf("RateLimit(%d): %v", rate, err)
	}
	r := gin.New()
	r.Use(h)
	r.GET("/t", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	return r
}

func TestRateLimit_HighRateAllowsRequest(t *testing.T) {
	r := rateLimitRouter(t, 1_000_000)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/t", nil))

	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestRateLimit_InvalidRate(t *testing.T) {
	for _, rate := range []int64{0, -1, -1000} {
		h, err := RateLimit(rate)
		if err == nil {
			t.Fatalf("RateLimit(%d) expected error", rate)
		}
		if h != nil {
			t.Fatalf("RateLimit(%d) handler must be nil on error", rate)
		}
	}
}

func TestRateLimit_ExhaustedBucketReturns429(t *testing.T) {
	r := rateLimitRouter(t, 1)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/t", nil))
	if first.Code != http.StatusNoContent {
		t.Fatalf("first status=%d", first.Code)
	}
	second := httptest.NewRecorder()
	r.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/t", nil))
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status=%d body=%s", second.Code, second.Body.String())
	}
}
