package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/def"

	"github.com/gin-gonic/gin"
)

func TestSetV_PropagatesRequestIDToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var got string
	r.Use(SetV)
	r.GET("/t", func(c *gin.Context) {
		got = def.RequestIDFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.Header.Set(def.XRequestID, "req-from-client")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got != "req-from-client" {
		t.Fatalf("context request id = %q", got)
	}
}

func TestSetV_GeneratesWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var got string
	r.Use(SetV)
	r.GET("/t", func(c *gin.Context) {
		got = def.RequestIDFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got == "" {
		t.Fatal("expected generated request id")
	}
}
