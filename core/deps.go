package core

import (
	libI18n "github.com/andywang2011/go-i18n"
	"github.com/gin-gonic/gin"
)

const (
	ginKeyI18nEngine = "core.i18n_engine"
	ginKeyAppEnv     = "core.app_env"
)

// SetCoreDepsOnGin stores the i18n engine and app env on the Gin context so the
// per-request Context can resolve them without global state.
func SetCoreDepsOnGin(c *gin.Context, engine *libI18n.Engine, appEnv string) {
	c.Set(ginKeyI18nEngine, engine)
	c.Set(ginKeyAppEnv, appEnv)
}

func i18nEngineFromGin(c *gin.Context) *libI18n.Engine {
	if v, ok := c.Get(ginKeyI18nEngine); ok {
		if engine, ok := v.(*libI18n.Engine); ok {
			return engine
		}
	}
	return nil
}

func appEnvFromGin(c *gin.Context) string {
	if v, ok := c.Get(ginKeyAppEnv); ok {
		if env, ok := v.(string); ok {
			return env
		}
	}
	return ""
}
