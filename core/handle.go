package core

import (
	"github.com/andywang2011/go-app-kit/def"

	"github.com/andywang2011/go-logging"
	"github.com/gin-gonic/gin"
)

// Ctx is the Gin context key under which the per-request *Context is cached.
const Ctx = "__context__"

// HandlerFunc is a kit handler operating on the enriched *Context.
type HandlerFunc func(c *Context)

// Handle adapts a kit HandlerFunc into a gin.HandlerFunc, building (and caching)
// the per-request *Context on first use.
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := getContext(c)
		h(ctx)
	}
}

func getContext(c *gin.Context) *Context {
	ctx, ok := c.Get(Ctx)
	if ok {
		return ctx.(*Context)
	}

	lang := c.GetHeader(def.Language)
	if len(lang) == 0 {
		lang = def.ZhLanguage
	}

	logger := logging.New()
	logger.SetRequestID(c.GetString(def.XRequestID))
	logger.Setv1(c.GetString("v1"))

	context := &Context{
		Env:      appEnvFromGin(c),
		Context:  c,
		Logger:   logger,
		Language: lang,
	}

	if engine := i18nEngineFromGin(c); engine != nil {
		context.Localizer = engine.NewLocalizer(lang)
	}

	c.Set(Ctx, context)
	return context
}
