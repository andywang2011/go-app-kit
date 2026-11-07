package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/core"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders_SetsBaseline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	core.Handle(SecurityHeaders)(c)

	for _, key := range []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"Referrer-Policy",
		"Content-Security-Policy",
	} {
		if w.Header().Get(key) == "" {
			t.Fatalf("missing header %s", key)
		}
	}
}

func hstsFor(t *testing.T, cfg SecurityHeadersConfig, tlsConn bool) string {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if tlsConn {
		c.Request.TLS = &tls.ConnectionState{}
	}
	core.Handle(NewSecurityHeaders(cfg))(c)
	return w.Header().Get("Strict-Transport-Security")
}

func TestSecurityHeaders_HSTS_NotSetOnPlaintext(t *testing.T) {
	if got := hstsFor(t, DefaultSecurityHeadersConfig(), false); got != "" {
		t.Fatalf("HSTS on plaintext = %q, want empty", got)
	}
}

func TestSecurityHeaders_HSTS_SetOnTLS(t *testing.T) {
	if got := hstsFor(t, DefaultSecurityHeadersConfig(), true); got == "" {
		t.Fatal("expected HSTS on TLS connection")
	}
}

func TestSecurityHeaders_HSTS_SetWhenTLSTermination(t *testing.T) {
	cfg := DefaultSecurityHeadersConfig()
	cfg.TLSTermination = true
	if got := hstsFor(t, cfg, false); got == "" {
		t.Fatal("expected HSTS when TLSTermination declared")
	}
}

func TestSecurityHeaders_HSTS_DisabledEvenOnTLS(t *testing.T) {
	cfg := DefaultSecurityHeadersConfig()
	cfg.HSTS = false
	cfg.TLSTermination = true
	if got := hstsFor(t, cfg, true); got != "" {
		t.Fatalf("HSTS = %q, want empty when HSTS=false", got)
	}
}

func TestSecurityHeaders_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	core.Handle(NewSecurityHeaders(SecurityHeadersConfig{Enabled: false}))(c)

	if w.Header().Get("X-Frame-Options") != "" {
		t.Fatal("expected no headers when disabled")
	}
}
