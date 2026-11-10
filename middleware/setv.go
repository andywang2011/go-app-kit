package middleware

import (
	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-libs/random"
	"github.com/andywang2011/go-logging"

	"github.com/gin-gonic/gin"
)

// SetV ensures every request carries an X-Request-Id (generating one when
// absent), stores it on the Gin context and Request.Context for outbound
// clients, and records the request path under the "v1" key for logging.
func SetV(c *gin.Context) {
	xRequestID := c.GetHeader(def.XRequestID)
	if xRequestID == "" {
		xRequestID = random.MustToken(32)
	}
	path := c.Request.URL.Path
	c.Set(def.XRequestID, xRequestID)
	c.Set("v1", path)
	ctx := def.WithRequestID(c.Request.Context(), xRequestID)
	ctx = logging.Attach(ctx, xRequestID, path, "", "")
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
