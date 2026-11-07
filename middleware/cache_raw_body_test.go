package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestCacheRawBody_RestoresBody verifies the body is both cached under
// gin.BodyBytesKey AND restored on Request.Body so downstream readers (e.g.
// ShouldBindJSON) can still consume it.
func TestCacheRawBody_RestoresBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBody())

	const payload = `{"name":"kit"}`
	var (
		cached  string
		drained string
	)
	r.POST("/t", func(c *gin.Context) {
		if raw, ok := c.Get(gin.BodyBytesKey); ok {
			if b, ok := raw.([]byte); ok {
				cached = string(b)
			}
		}
		body, _ := io.ReadAll(c.Request.Body)
		drained = string(body)
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/t", strings.NewReader(payload))
	r.ServeHTTP(w, req)

	if cached != payload {
		t.Fatalf("cached body = %q, want %q", cached, payload)
	}
	if drained != payload {
		t.Fatalf("restored body = %q, want %q (body was not restored)", drained, payload)
	}
}

func TestCacheRawBody_SkipsGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBody())
	r.GET("/t", func(c *gin.Context) {
		if _, ok := c.Get(gin.BodyBytesKey); ok {
			t.Error("GET should not cache body")
		}
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d", w.Code)
	}
}
