package middleware

import (
	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-libs/random"

	"github.com/gin-gonic/gin"
)

// SetV ensures every request carries an X-Request-Id (generating one when
// absent) and records the request path under the "v1" key for logging.
func SetV(c *gin.Context) {
	xRequestID := c.GetHeader(def.XRequestID)
	if xRequestID == "" {
		xRequestID = random.MustToken(32)
	}
	c.Set(def.XRequestID, xRequestID)
	c.Set("v1", c.Request.URL.Path)
	c.Next()
}
