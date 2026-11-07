// Package core provides the per-request Context used by handlers built on the
// kit, the Handle decorator that adapts handlers to Gin, and the JSON response
// envelope (Success / Fail / FailWithErrCode).
//
// Context intentionally carries NO concrete business identity fields (no
// Member/Admin/UserInfo). Authenticated identity is exposed only through the
// Subject interface via SetSubject/Subject, so the kit stays free of any
// caller-specific session types. Callers store their concrete session on the
// embedded *gin.Context keys and adapt it to Subject.
package core

import (
	"net/http"
	"time"

	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-app-kit/errcode"

	"github.com/andywang2011/go-logging"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Subject is the authenticated principal of a request. It is a narrow,
// business-agnostic view: callers adapt their concrete session (member/admin)
// to this interface and attach it with Context.SetSubject.
type Subject interface {
	// ID returns the principal's numeric identifier.
	ID() int
	// Type returns the principal kind, e.g. "member" or "admin".
	Type() string
	// ClientType returns the client channel the principal authenticated from.
	ClientType() string
}

// Context is the per-request context passed to kit handlers. It embeds
// *gin.Context (full Gin access) and *logging.Logger (request-scoped logging).
type Context struct {
	// Env is the current application environment (see def.Env*).
	Env string
	// Language is the resolved client language code.
	Language string
	// TimeLocal is the request time zone, set by the SetTimeZone middleware.
	TimeLocal *time.Location
	// Localizer translates i18n keys for the resolved Language.
	Localizer *i18n.Localizer

	*gin.Context
	*logging.Logger

	subject Subject
}

// SetSubject attaches the authenticated principal to the context.
func (c *Context) SetSubject(s Subject) {
	c.subject = s
}

// Subject returns the authenticated principal, or nil when unauthenticated.
func (c *Context) Subject() Subject {
	return c.subject
}

// Success writes a 200 success envelope with the given data payload.
func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

// Fail writes an error envelope. The debug detail is omitted in production.
func (c *Context) Fail(code int, msg string, detail, data interface{}) {
	if c.Env == def.EnvProduction {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":   code,
			"msg":    msg,
			"detail": detail,
			"data":   data,
		})
	}
}

// FailWithErrCode writes an error envelope from an ErrCode, localizing the
// message through Localizer when available.
func (c *Context) FailWithErrCode(err *errcode.ErrCode, data interface{}) {
	if err == nil {
		c.Fail(500, "Internal error", nil, data)
		return
	}

	msg := err.Msg
	if c.Localizer == nil {
		if len(msg) == 0 {
			msg = err.Key
		}
		c.Fail(err.Code, msg, err.Detail, data)
		return
	}

	if len(err.Key) > 0 {
		if translated, e := c.Localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    err.Key,
			TemplateData: err.Args,
		}); e == nil && len(translated) > 0 {
			msg = translated
		}
	}

	if len(msg) == 0 {
		msg = err.Key
	}

	c.Fail(err.Code, msg, err.Detail, data)
}
