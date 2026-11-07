package logging

import (
	"errors"
	"testing"

	goLog "github.com/andywang2011/go-logging"
)

func resetSetup(t *testing.T, fn func(goLog.Config) error) {
	t.Helper()
	oldFn := setupFn
	setupMu.Lock()
	setupDone = false
	setupMu.Unlock()
	setupFn = fn
	t.Cleanup(func() {
		setupFn = oldFn
		setupMu.Lock()
		setupDone = false
		setupMu.Unlock()
	})
}

func TestSetup_ReturnsErrorInsteadOfFatal(t *testing.T) {
	boom := errors.New("bad logging config")
	resetSetup(t, func(goLog.Config) error { return boom })

	err := Setup()
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want wrapped %v", err, boom)
	}
}

func TestSetup_RetryAfterFailureThenIdempotent(t *testing.T) {
	calls := 0
	resetSetup(t, func(goLog.Config) error {
		calls++
		if calls == 1 {
			return errors.New("first fails")
		}
		return nil
	})

	if err := Setup(); err == nil {
		t.Fatal("expected first error")
	}
	if err := Setup(); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if err := Setup(); err != nil {
		t.Fatalf("third: %v", err)
	}
	if calls != 2 {
		t.Fatalf("setup calls = %d, want 2 (no re-init after success)", calls)
	}
}

func TestClose(t *testing.T) {
	if err := Close(); err != nil {
		t.Logf("Close: %v", err) // go-logging may report nothing-to-close; must not panic
	}
}
