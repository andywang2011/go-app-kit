package appkit

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

// fakeModule verifies the Module interface is implementable.
type fakeModule struct {
	name       string
	enabled    bool
	registered bool
}

func (m *fakeModule) Name() string              { return m.name }
func (m *fakeModule) Enabled() bool             { return m.enabled }
func (m *fakeModule) Register(*gin.RouterGroup) { m.registered = true }

// fakeCloser verifies the Closer interface is implementable.
type fakeCloser struct{ err error }

func (c fakeCloser) Close() error { return c.err }

// compile-time assertions.
var (
	_ Module = (*fakeModule)(nil)
	_ Closer = fakeCloser{}
)

func TestModuleContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := &fakeModule{name: "demo", enabled: true}
	if m.Name() != "demo" || !m.Enabled() {
		t.Fatalf("module metadata: %+v", m)
	}
	g := gin.New().Group("/x")
	m.Register(g)
	if !m.registered {
		t.Fatal("Register not invoked")
	}
}

func TestCloserContract(t *testing.T) {
	if err := (fakeCloser{}).Close(); err != nil {
		t.Fatalf("nil close: %v", err)
	}
	want := errors.New("boom")
	if err := (fakeCloser{err: want}).Close(); !errors.Is(err, want) {
		t.Fatalf("close err = %v", err)
	}
}
