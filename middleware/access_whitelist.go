package middleware

import (
	"github.com/andywang2011/go-app-kit/core"
	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-app-kit/errcode"
	"github.com/andywang2011/go-libs/types"

	"github.com/gin-gonic/gin"
)

// AccessWhitelistConfig configures the non-production IP whitelist. It is a
// kit-owned struct so the middleware never imports a caller's wire/config.
type AccessWhitelistConfig struct {
	// Enabled turns the whitelist on in non-production, non-testing envs.
	Enabled bool
	// CIDRs is the set of allowed client IPs (matched exactly).
	CIDRs []string
}

// AccessWhitelist restricts access by client IP in non-production environments.
//
// The client IP comes from gin's Context.ClientIP, which honors forwarded
// headers only for proxies listed in the engine's TrustedProxies. Callers behind
// a reverse proxy must configure gin.Engine.SetTrustedProxies accordingly;
// otherwise forged X-Forwarded-For headers cannot influence the result.
type AccessWhitelist struct {
	env string
	cfg AccessWhitelistConfig
}

// NewAccessWhitelist creates IP whitelist middleware from the given env and
// config. Production always bypasses; testing always enforces.
func NewAccessWhitelist(env string, cfg AccessWhitelistConfig) *AccessWhitelist {
	return &AccessWhitelist{env: env, cfg: cfg}
}

// Handle is the kit handler enforcing the whitelist.
func (m *AccessWhitelist) Handle(c *core.Context) {
	if m.env == def.EnvProduction {
		c.Next()
		return
	}
	if m.env != def.EnvTesting && !m.cfg.Enabled {
		c.Next()
		return
	}

	userIP := c.ClientIP()
	if userIP != "" && types.SliceContain(m.cfg.CIDRs, userIP) {
		c.Next()
		return
	}

	c.FailWithErrCode(errcode.LibNotInWhitelist, gin.H{
		"ip":      userIP,
		"enabled": m.cfg.Enabled,
	})
	c.Abort()
}
