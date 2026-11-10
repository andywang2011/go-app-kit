package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-logging"

	"github.com/gin-gonic/gin"
)

func TestSetV_PropagatesRequestIDToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var gotDef, gotLog string
	r.Use(SetV)
	r.GET("/t", func(c *gin.Context) {
		gotDef = def.RequestIDFromContext(c.Request.Context())
		gotLog = logging.FromContext(c.Request.Context()).RequestID
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.Header.Set(def.XRequestID, "req-from-client")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if gotDef != "req-from-client" {
		t.Fatalf("def request id = %q", gotDef)
	}
	if gotLog != "req-from-client" {
		t.Fatalf("logging request id = %q", gotLog)
	}
}

func TestSetV_GeneratesWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var gotDef, gotLog, gotV1 string
	r.Use(SetV)
	r.GET("/t", func(c *gin.Context) {
		ctx := c.Request.Context()
		gotDef = def.RequestIDFromContext(ctx)
		l := logging.FromContext(ctx)
		gotLog = l.RequestID
		gotV1 = l.V1
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if gotDef == "" || gotLog == "" {
		t.Fatal("expected generated request id")
	}
	if gotDef != gotLog {
		t.Fatalf("def=%q logging=%q", gotDef, gotLog)
	}
	if gotV1 != "/t" {
		t.Fatalf("v1 = %q", gotV1)
	}
}
