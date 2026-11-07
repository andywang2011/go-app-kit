package errcode

import (
	"fmt"
	"testing"
)

func TestWebInvalidParam(t *testing.T) {
	if WebInvalidParam == nil || WebInvalidParam.Code != 40100 {
		t.Fatalf("WebInvalidParam: %+v", WebInvalidParam)
	}
}

func TestWebNetworkBusy(t *testing.T) {
	if WebNetworkBusy == nil || WebNetworkBusy.Key == "" {
		t.Fatalf("WebNetworkBusy: %+v", WebNetworkBusy)
	}
}

func TestErrorPrefersDetailThenMsgThenKey(t *testing.T) {
	e := &ErrCode{Code: 1, Key: "k"}
	if e.Error() != "k" {
		t.Fatalf("want key, got %q", e.Error())
	}
	e.Msg = "m"
	if e.Error() != "m" {
		t.Fatalf("want msg, got %q", e.Error())
	}
	e.Detail = "d"
	if e.Error() != "d" {
		t.Fatalf("want detail, got %q", e.Error())
	}
	var nilEC *ErrCode
	if nilEC.Error() != "" {
		t.Fatalf("nil ErrCode should be empty")
	}
}

func TestWithDetailArgsErr(t *testing.T) {
	base := WebTokenExpired
	d := base.WithDetail("expired")
	if d.Detail != "expired" || d.Code != base.Code || d.Key != base.Key {
		t.Fatalf("WithDetail: %+v", d)
	}
	a := base.WithArgs(map[string]interface{}{"n": 1})
	if a.Args["n"] != 1 {
		t.Fatalf("WithArgs: %+v", a)
	}
	if base.WithErr(nil) != base {
		t.Fatal("WithErr(nil) should return receiver")
	}
	w := base.WithErr(fmt.Errorf("boom"))
	if w.Detail != "boom" {
		t.Fatalf("WithErr: %+v", w)
	}
}

func TestIsAndAsErrCode(t *testing.T) {
	err := WebTokenExpired.WithDetail("x")
	if !WebTokenExpired.Is(err) {
		t.Fatal("Is should match same code+key")
	}
	if WebInvalidParam.Is(err) {
		t.Fatal("Is should not match different code")
	}
	got, ok := AsErrCode(fmt.Errorf("wrap: %w", err))
	if !ok || got.Code != WebTokenExpired.Code {
		t.Fatalf("AsErrCode: %+v ok=%v", got, ok)
	}
	if _, ok := AsErrCode(fmt.Errorf("plain")); ok {
		t.Fatal("AsErrCode should fail on plain error")
	}
}
