package middleware

import (
	"github.com/andywang2011/go-app-kit/core"

	libI18n "github.com/andywang2011/go-i18n"
	"github.com/gin-gonic/gin"
)

// InjectCoreDeps stores Wire-injected core dependencies (i18n engine, app env)
// on the Gin context so core.Context can resolve them per request.
func InjectCoreDeps(engine *libI18n.Engine, appEnv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		core.SetCoreDepsOnGin(c, engine, appEnv)
		c.Next()
	}
}
