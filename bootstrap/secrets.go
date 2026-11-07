package bootstrap

import (
	"fmt"
	"strings"

	"github.com/andywang2011/go-app-kit/def"

	"github.com/spf13/viper"
)

// DefaultProdSecretKeys are the keys that must be non-placeholder outside
// development (values should come from Nacos or APP_* env).
var DefaultProdSecretKeys = []string{
	"mysql.default.password",
	"jwt.secret",
	"security.blind_index_key",
}

var placeholderSecrets = map[string]struct{}{
	"":                                  {},
	"xxx":                               {},
	"xxxx":                              {},
	"xxxxx":                             {},
	"xxxxxx":                            {},
	"123456":                            {},
	"0123456789abcdef":                  {},
	"changeme":                          {},
	"password":                          {},
	"secret":                            {},
	"token_test111":                     {},
	"sk-xxxx":                           {},
	"local":                             {},
	"guest":                             {},
	"ftppass":                           {},
	"local-dev-aes-16":                  {},
	"local-dev-jwt-secret-not-for-prod": {},
	"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef": {},
	"fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210": {},
}

// ApplyEnvOverrides enables APP_* environment overrides for nested keys
// (e.g. APP_MYSQL_DEFAULT_PASSWORD → mysql.default.password).
func ApplyEnvOverrides() {
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// IsPlaceholderSecret reports whether v is empty or a known unsafe placeholder.
func IsPlaceholderSecret(v string) bool {
	_, ok := placeholderSecrets[strings.TrimSpace(v)]
	return ok
}

// requiresSecretValidation validates secrets for testing/staging/production
// (expect Nacos); development is skipped for local -skip_nacos work.
func requiresSecretValidation(env string) bool {
	switch env {
	case def.EnvTesting, def.EnvStaging, def.EnvProduction:
		return true
	default:
		return false
	}
}

// ValidateSecrets rejects empty/placeholder secrets for non-development
// environments. Logging AES masking was removed in go-logging v0.2; app.aes.key
// is no longer required (mobile encrypt uses aes.backend.key / aes.frontend.key).
func ValidateSecrets(env string, keys ...string) error {
	if !requiresSecretValidation(env) {
		return nil
	}
	if len(keys) == 0 {
		keys = DefaultProdSecretKeys
	}
	var missing []string
	for _, key := range keys {
		if key == "app.aes.key" {
			continue
		}
		if IsPlaceholderSecret(viper.GetString(key)) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s secrets missing or placeholder: %s (set via Nacos or APP_* env)", env, strings.Join(missing, ", "))
	}
	return nil
}
