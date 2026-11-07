package middleware

import (
	"net/http"
	"strings"

	"github.com/andywang2011/go-app-kit/core"
)

// CORSConfig controls browser CORS behavior. Empty AllowOrigins denies cross-origin
// requests (no Access-Control-Allow-Origin header). Use explicit origins in production.
type CORSConfig struct {
	AllowOrigins     []string
	AllowCredentials bool
	AllowHeaders     string
	AllowMethods     string
}

// DefaultCORSConfig is safe for local development with common localhost origins.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		},
		AllowCredentials: true,
		AllowHeaders:     "Authorization, Authorization-Token, Authorization-AppKey, Content-Type, Language, X-Request-Id",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE, PATCH",
	}
}

// NewCORS returns a CORS middleware bound to cfg.
func NewCORS(cfg CORSConfig) func(c *core.Context) {
	allowed := make(map[string]struct{}, len(cfg.AllowOrigins))
	for _, o := range cfg.AllowOrigins {
		o = strings.TrimSpace(o)
		if o != "" && o != "*" {
			allowed[o] = struct{}{}
		}
	}
	headers := cfg.AllowHeaders
	if headers == "" {
		headers = DefaultCORSConfig().AllowHeaders
	}
	methods := cfg.AllowMethods
	if methods == "" {
		methods = DefaultCORSConfig().AllowMethods
	}

	return func(c *core.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				if cfg.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			}
		}
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Allow-Methods", methods)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// CORS is the default development CORS middleware (localhost whitelist, not *).
func CORS(c *core.Context) {
	NewCORS(DefaultCORSConfig())(c)
}
