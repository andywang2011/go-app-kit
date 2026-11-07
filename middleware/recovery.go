// Package middleware provides reusable Gin middleware for services built on the
// kit: panic recovery, request logging, CORS, security headers, rate limiting,
// IP whitelisting, raw-body caching, request id/time-zone setup, Prometheus RED
// metrics and core dependency injection.
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-logging"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// maxRecoveryBodyLogBytes caps the request body included in panic logs.
const maxRecoveryBodyLogBytes = 4 << 10

// recoveryLogFields builds the panic log payload with headers and body redacted
// by the default redactor. Non-JSON or invalid JSON bodies are omitted because
// they cannot be redacted reliably.
func recoveryLogFields(c *gin.Context) map[string]interface{} {
	bodyStr := ""
	if raw, ok := c.Get(gin.BodyBytesKey); ok {
		if byts, _ := raw.([]byte); len(byts) > 0 {
			trim := strings.TrimSpace(string(byts))
			if trim != "" && (trim[0] == '{' || trim[0] == '[') && json.Valid(byts) {
				bodyStr = defaultLogRedactor.RedactJSONBody(string(byts))
				if len(bodyStr) > maxRecoveryBodyLogBytes {
					bodyStr = bodyStr[:maxRecoveryBodyLogBytes]
				}
			} else {
				bodyStr = "[non-JSON body omitted]"
			}
		}
	}
	return map[string]interface{}{
		"body":   bodyStr,
		"url":    c.Request.URL.Path,
		"header": defaultLogRedactor.RedactHeaders(c.Request.Header),
	}
}

// Recovery recovers from panics, logs the stack with the request id (headers
// and body redacted), and responds 500.
func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var errMsg string
			switch e := err.(type) {
			case error:
				errMsg = fmt.Sprintf("%+v", errors.WithStack(e))
			default:
				errMsg = fmt.Sprintf("%v", err)
			}
			log := recoveryLogFields(c)

			xRequestID := c.GetString(def.XRequestID)
			logger := logging.New()
			logger.SetRequestID(xRequestID)
			logger.Fatal("internal server error", log, errMsg)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}
