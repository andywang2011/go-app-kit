package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/core"
	"github.com/andywang2011/go-app-kit/def"

	"github.com/gin-gonic/gin"
)

func newWhitelistRouter(t *testing.T, env string, cfg AccessWhitelistConfig, trusted []string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := r.SetTrustedProxies(trusted); err != nil {
		t.Fatal(err)
	}
	r.Use(core.Handle(NewAccessWhitelist(env, cfg).Handle))
	r.GET("/t", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	return r
}

func doWhitelistReq(r *gin.Engine, remoteAddr string, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	r.ServeHTTP(w, req)
	return w
}

func TestAccessWhitelist_ForgedXFFCannotBypass(t *testing.T) {
	cfg := AccessWhitelistConfig{Enabled: true, CIDRs: []string{"10.0.0.1"}}
	// No trusted proxies: forwarded headers must be ignored.
	r := newWhitelistRouter(t, def.EnvTesting, cfg, nil)

	w := doWhitelistReq(r, "203.0.113.9:1234", map[string]string{
		"X-Forwarded-For":          "10.0.0.1",
		"X-Original-Forwarded-For": "10.0.0.1",
	})
	if w.Code == http.StatusNoContent {
		t.Fatal("forged X-Forwarded-For bypassed whitelist")
	}
}

func TestAccessWhitelist_AllowsWhitelistedRemoteIP(t *testing.T) {
	cfg := AccessWhitelistConfig{Enabled: true, CIDRs: []string{"10.0.0.1"}}
	r := newWhitelistRouter(t, def.EnvTesting, cfg, nil)

	w := doWhitelistReq(r, "10.0.0.1:1234", nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d want 204", w.Code)
	}
}

func TestAccessWhitelist_TrustedProxyXFFHonored(t *testing.T) {
	cfg := AccessWhitelistConfig{Enabled: true, CIDRs: []string{"10.0.0.1"}}
	r := newWhitelistRouter(t, def.EnvTesting, cfg, []string{"192.0.2.1"})

	w := doWhitelistReq(r, "192.0.2.1:1234", map[string]string{"X-Forwarded-For": "10.0.0.1"})
	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d want 204 via trusted proxy", w.Code)
	}
}

func TestAccessWhitelist_UntrustedProxyXFFIgnored(t *testing.T) {
	cfg := AccessWhitelistConfig{Enabled: true, CIDRs: []string{"10.0.0.1"}}
	r := newWhitelistRouter(t, def.EnvTesting, cfg, []string{"192.0.2.1"})

	w := doWhitelistReq(r, "198.51.100.7:1234", map[string]string{"X-Forwarded-For": "10.0.0.1"})
	if w.Code == http.StatusNoContent {
		t.Fatal("XFF from untrusted peer bypassed whitelist")
	}
}

func TestAccessWhitelist_ProductionBypassesAndDisabledBypasses(t *testing.T) {
	cfg := AccessWhitelistConfig{Enabled: true, CIDRs: []string{"10.0.0.1"}}
	if w := doWhitelistReq(newWhitelistRouter(t, def.EnvProduction, cfg, nil), "203.0.113.9:1", nil); w.Code != http.StatusNoContent {
		t.Fatalf("production status=%d", w.Code)
	}
	off := AccessWhitelistConfig{Enabled: false}
	if w := doWhitelistReq(newWhitelistRouter(t, "dev", off, nil), "203.0.113.9:1", nil); w.Code != http.StatusNoContent {
		t.Fatalf("disabled status=%d", w.Code)
	}
}
