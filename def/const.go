// Package def holds framework-level constants shared by the HTTP kit:
// environment names, language codes and well-known request header keys.
//
// Business constants (status enums, password policies, page sizes, ...) do NOT
// belong here; callers keep those in their own def package.
package def

// Application environment names. app.env must be one of these values.
const (
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvTesting     = "testing"
	EnvStaging     = "staging"
)

// Language header value and supported language codes.
const (
	// Language is the HTTP header carrying the client language preference.
	Language = "Language"

	// ZhLanguage is the Simplified Chinese language code (default).
	ZhLanguage = "zh"
	// EnLanguage is the English language code.
	EnLanguage = "en"
	// JaLanguage is the Japanese language code.
	JaLanguage = "ja"
)

// Well-known request header keys used by the kit middleware and context.
const (
	// AuthorizationToken carries the session token for authenticated requests.
	AuthorizationToken = "Authorization-Token"
	// XRequestID carries/propagates the per-request trace id.
	XRequestID = "X-Request-Id"
)
