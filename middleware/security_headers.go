package middleware

import "github.com/andywang2011/go-app-kit/core"

// SecurityHeadersConfig toggles baseline browser security response headers.
type SecurityHeadersConfig struct {
	Enabled bool
	// HSTS enables Strict-Transport-Security, but only on TLS connections or
	// when TLSTermination is set.
	HSTS bool
	// TLSTermination declares that TLS is terminated by a trusted upstream
	// (load balancer / ingress), so HSTS is also sent on plaintext hops.
	TLSTermination        bool
	ContentSecurityPolicy string
}

// DefaultSecurityHeadersConfig enables common hardening headers.
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		Enabled:               true,
		HSTS:                  true,
		ContentSecurityPolicy: "default-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
	}
}

// NewSecurityHeaders returns middleware that sets nosniff, frame deny, referrer, CSP, and
// HSTS (only over TLS or when cfg.TLSTermination is set).
func NewSecurityHeaders(cfg SecurityHeadersConfig) func(c *core.Context) {
	return func(c *core.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}
		if cfg.HSTS && (c.Request.TLS != nil || cfg.TLSTermination) {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

// SecurityHeaders is the default security headers middleware.
func SecurityHeaders(c *core.Context) {
	NewSecurityHeaders(DefaultSecurityHeadersConfig())(c)
}
