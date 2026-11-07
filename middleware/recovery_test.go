package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func recoveryCtx(t *testing.T, body []byte, headers map[string]string) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/p", nil)
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	if body != nil {
		c.Set(gin.BodyBytesKey, body)
	}
	return c
}

func TestRecoveryLogFields_RedactsHeadersAndBody(t *testing.T) {
	c := recoveryCtx(t,
		[]byte(`{"password":"hunter2","user":"bob","nested":{"token":"abc123"}}`),
		map[string]string{"Authorization": "Bearer top-secret", "Cookie": "sid=xyz", "X-Ok": "1"},
	)
	f := recoveryLogFields(c)

	hdr, _ := f["header"].(map[string]any)
	if hdr["Authorization"] != "***" || hdr["Cookie"] != "***" {
		t.Fatalf("headers not redacted: %v", hdr)
	}
	if hdr["X-Ok"] != "1" {
		t.Fatalf("innocent header changed: %v", hdr)
	}
	body, _ := f["body"].(string)
	for _, leak := range []string{"hunter2", "abc123"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body leaked %q: %s", leak, body)
		}
	}
	if !strings.Contains(body, "bob") {
		t.Fatalf("innocent field dropped: %s", body)
	}
	if f["url"] != "/p" {
		t.Fatalf("url=%v", f["url"])
	}
}

func TestRecoveryLogFields_NonJSONBodyOmitted(t *testing.T) {
	c := recoveryCtx(t, []byte("password=hunter2&user=bob"), nil)
	body, _ := recoveryLogFields(c)["body"].(string)
	if strings.Contains(body, "hunter2") {
		t.Fatalf("non-JSON body leaked: %s", body)
	}
	if body == "" {
		t.Fatal("expected omission marker")
	}
}

func TestRecoveryLogFields_InvalidJSONOmitted(t *testing.T) {
	c := recoveryCtx(t, []byte(`{"password":"hunter2"`), nil)
	body, _ := recoveryLogFields(c)["body"].(string)
	if strings.Contains(body, "hunter2") {
		t.Fatalf("invalid JSON body leaked: %s", body)
	}
}

func TestRecoveryLogFields_NoBodyAndTruncation(t *testing.T) {
	if body := recoveryLogFields(recoveryCtx(t, nil, nil))["body"]; body != "" {
		t.Fatalf("body=%v want empty", body)
	}
	big := `{"data":"` + strings.Repeat("x", 2*maxRecoveryBodyLogBytes) + `"}`
	body, _ := recoveryLogFields(recoveryCtx(t, []byte(big), nil))["body"].(string)
	if len(body) > maxRecoveryBodyLogBytes {
		t.Fatalf("len=%d exceeds cap", len(body))
	}
}

func TestRecovery_PanicReturns500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery)
	r.POST("/boom", func(c *gin.Context) { panic("boom") })
	r.POST("/err", func(c *gin.Context) { panic(http.ErrAbortHandler) })
	r.POST("/ok", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	for path, want := range map[string]int{
		"/boom": http.StatusInternalServerError,
		"/err":  http.StatusInternalServerError,
		"/ok":   http.StatusNoContent,
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, nil))
		if w.Code != want {
			t.Fatalf("%s status=%d want %d", path, w.Code, want)
		}
	}
}
