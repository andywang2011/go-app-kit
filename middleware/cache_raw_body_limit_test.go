package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCacheRawBody_OversizedBodyReturns413(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBodyWithLimit(16))
	called := false
	r.POST("/t", func(c *gin.Context) {
		called = true
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/t", strings.NewReader(strings.Repeat("a", 17)))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status=%d want 413", w.Code)
	}
	if called {
		t.Fatal("handler must not run for oversized body")
	}
}

func TestCacheRawBody_AtLimitAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBodyWithLimit(16))
	r.POST("/t", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/t", strings.NewReader(strings.Repeat("a", 16)))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d want 204", w.Code)
	}
}

func TestCacheRawBody_DefaultLimitAppliesWhenNonPositive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBodyWithLimit(0))
	r.POST("/t", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/t", strings.NewReader("small"))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d", w.Code)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, http.ErrBodyReadAfterClose }

func TestCacheRawBody_ReadErrorReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheRawBody())
	r.POST("/t", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/t", errReader{})
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want 400", w.Code)
	}
}
