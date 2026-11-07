package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/core"

	"github.com/gin-gonic/gin"
)

func TestCORS_OptionsAbortsWithNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodOptions, "/ping", nil)
	c.Request.Header.Set("Origin", "http://localhost:3000")

	core.Handle(CORS)(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status: %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("Allow-Origin=%q", got)
	}
}

func TestCORS_RejectsUnknownOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ping", nil)
	c.Request.Header.Set("Origin", "https://evil.example")

	core.Handle(NewCORS(DefaultCORSConfig()))(c)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected Allow-Origin=%q", got)
	}
}

func TestCORS_NeverStar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ping", nil)
	c.Request.Header.Set("Origin", "*")

	core.Handle(NewCORS(CORSConfig{AllowOrigins: []string{"*"}}))(c)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got == "*" {
		t.Fatal("wildcard origin must not be reflected")
	}
}
