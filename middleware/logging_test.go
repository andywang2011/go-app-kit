package middleware

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
)

func TestRedactQueryForLog_RedactsSensitiveKeys(t *testing.T) {
	q := url.Values{"token": {"abc"}, "password": {"p"}, "page": {"1"}, "tag": {"a", "b"}}
	got := redactQueryForLog(q)
	if got["token"] != "***" || got["password"] != "***" {
		t.Fatalf("sensitive query not redacted: %v", got)
	}
	if got["page"] != "1" {
		t.Fatalf("page=%v", got["page"])
	}
	if tags, ok := got["tag"].([]string); !ok || len(tags) != 2 {
		t.Fatalf("tag=%v", got["tag"])
	}
}

func TestBodyLogWriter_CapsCapturedBody(t *testing.T) {
	w := &BodyLogWriter{body: bytes.NewBuffer(nil), limit: 8}
	w.capture([]byte("12345"))
	if w.truncated {
		t.Fatal("unexpected truncation")
	}
	w.capture([]byte("6789abc"))
	if !w.truncated || w.body.String() != "12345678" {
		t.Fatalf("body=%q truncated=%v", w.body.String(), w.truncated)
	}
	w.capture([]byte("more"))
	if w.body.Len() != 8 {
		t.Fatalf("len=%d, want capped at 8", w.body.Len())
	}
}

func TestBodyLogWriter_DefaultLimit(t *testing.T) {
	w := &BodyLogWriter{body: bytes.NewBuffer(nil)}
	w.capture(bytes.Repeat([]byte("a"), maxLogRespBodyBytes+10))
	if w.body.Len() != maxLogRespBodyBytes || !w.truncated {
		t.Fatalf("len=%d truncated=%v", w.body.Len(), w.truncated)
	}
}

func TestResponseLogBody(t *testing.T) {
	w := &BodyLogWriter{body: bytes.NewBufferString(`{"token":"abc","ok":1}`), limit: 100}
	got := string(responseLogBody(w))
	if strings.Contains(got, "abc") {
		t.Fatalf("response token leaked: %s", got)
	}

	trunc := &BodyLogWriter{body: bytes.NewBufferString(`{"token":"ab`), limit: 12, truncated: true}
	got = string(responseLogBody(trunc))
	if strings.Contains(got, "ab") && strings.Contains(got, "token") {
		t.Fatalf("truncated body content leaked: %s", got)
	}

	empty := &BodyLogWriter{body: bytes.NewBuffer(nil)}
	if len(responseLogBody(empty)) != 0 {
		t.Fatal("expected empty")
	}
}
