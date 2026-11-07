package core

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-app-kit/errcode"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func newTestContext() (*Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	gc, _ := gin.CreateTestContext(w)
	gc.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	return &Context{Context: gc}, w
}

func decode(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("unmarshal %q: %v", w.Body.String(), err)
	}
	return m
}

func TestContextSuccess(t *testing.T) {
	c, w := newTestContext()
	c.Success(gin.H{"ok": true})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	m := decode(t, w)
	if m["code"].(float64) != 200 || m["msg"] != "success" {
		t.Fatalf("body = %v", m)
	}
}

func TestContextFailHidesDetailInProduction(t *testing.T) {
	c, w := newTestContext()
	c.Env = def.EnvProduction
	c.Fail(40100, "bad", "secret detail", nil)
	m := decode(t, w)
	if _, ok := m["detail"]; ok {
		t.Fatalf("detail must be hidden in production: %v", m)
	}

	c2, w2 := newTestContext()
	c2.Env = def.EnvDevelopment
	c2.Fail(40100, "bad", "shown detail", nil)
	m2 := decode(t, w2)
	if m2["detail"] != "shown detail" {
		t.Fatalf("detail must be shown in dev: %v", m2)
	}
}

func TestContextFailWithErrCode(t *testing.T) {
	c, w := newTestContext()
	c.FailWithErrCode(errcode.WebTokenExpired, nil)
	m := decode(t, w)
	if int(m["code"].(float64)) != errcode.WebTokenExpired.Code {
		t.Fatalf("code = %v", m["code"])
	}

	c2, w2 := newTestContext()
	c2.FailWithErrCode(nil, nil)
	m2 := decode(t, w2)
	if int(m2["code"].(float64)) != 500 {
		t.Fatalf("nil errcode should map to 500: %v", m2)
	}
}

type fakeSubject struct {
	id  int
	typ string
	ct  string
}

func (f fakeSubject) ID() int            { return f.id }
func (f fakeSubject) Type() string       { return f.typ }
func (f fakeSubject) ClientType() string { return f.ct }

func TestContextSubject(t *testing.T) {
	c, _ := newTestContext()
	if c.Subject() != nil {
		t.Fatal("subject should start nil")
	}
	c.SetSubject(fakeSubject{id: 7, typ: "member", ct: "web"})
	s := c.Subject()
	if s == nil || s.ID() != 7 || s.Type() != "member" || s.ClientType() != "web" {
		t.Fatalf("subject = %+v", s)
	}
}

// compile-time check that fakeSubject satisfies Subject.
var _ Subject = fakeSubject{}
