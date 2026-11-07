// Package errcode defines the ErrCode value type carried across the HTTP kit
// and the framework/web error codes shared by every service built on the kit.
//
// ErrCode carries a numeric Code and an i18n Key but never the localized text
// itself; translation happens at the response boundary (see core.Context
// FailWithErrCode). Business/auth error codes live in the caller's own errcode
// package, which may embed or alias this type.
package errcode

import "errors"

// ErrCode is a transport-level error descriptor: a stable numeric code plus an
// i18n message key. It implements the error interface so it can flow through
// errors.Is / errors.As.
type ErrCode struct {
	Code   int                    `json:"code"`
	Msg    string                 `json:"msg"`
	Key    string                 `json:"-"` // i18n message key
	Detail string                 `json:"-"` // debug-only detail
	Args   map[string]interface{} `json:"-"` // i18n template args
}

// Error returns the most descriptive non-empty text: Detail, then Msg, then Key.
func (e *ErrCode) Error() string {
	if e == nil {
		return ""
	}
	if e.Detail != "" {
		return e.Detail
	}
	if e.Msg != "" {
		return e.Msg
	}
	return e.Key
}

// String implements fmt.Stringer and mirrors Error.
func (e *ErrCode) String() string {
	return e.Error()
}

// Is enables errors.Is(err, errcode.Xxx) by matching Code+Key.
func (e *ErrCode) Is(target error) bool {
	t, ok := target.(*ErrCode)
	if !ok || e == nil || t == nil {
		return false
	}
	return e.Code == t.Code && e.Key == t.Key
}

// AsErrCode extracts an *ErrCode from err (including wrapped errors).
func AsErrCode(err error) (*ErrCode, bool) {
	var ec *ErrCode
	if errors.As(err, &ec) {
		return ec, true
	}
	return nil, false
}

// WithDetail returns a copy with an added debug detail (i18n key is preserved).
func (e *ErrCode) WithDetail(detail string) *ErrCode {
	return &ErrCode{
		Code:   e.Code,
		Key:    e.Key,
		Msg:    e.Msg,
		Detail: detail,
	}
}

// WithArgs returns a copy carrying i18n template arguments.
func (e *ErrCode) WithArgs(args map[string]interface{}) *ErrCode {
	return &ErrCode{
		Code:   e.Code,
		Key:    e.Key,
		Msg:    e.Msg,
		Detail: e.Detail,
		Args:   args,
	}
}

// WithErr returns a copy whose Detail is set from err. It returns the receiver
// unchanged when err is nil.
func (e *ErrCode) WithErr(err error) *ErrCode {
	if err == nil {
		return e
	}
	return &ErrCode{
		Code:   e.Code,
		Key:    e.Key,
		Msg:    e.Msg,
		Detail: err.Error(),
	}
}
