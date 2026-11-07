// Package config holds pure, reusable application configuration structs and
// their viper Load helpers. It deliberately contains no business/service
// dependencies: callers keep service-specific config (e.g. AES keys bound to a
// service type) in their own wire/config package.
//
// Like bootstrap, this package uses viper only for the optional Load* helpers;
// the structs themselves are plain data and can be built without viper.
package config

import (
	"time"

	"github.com/andywang2011/go-app-kit/middleware"

	"github.com/spf13/viper"
)

// ServerTimeouts configures net/http.Server deadline fields for slow-client
// protection.
type ServerTimeouts struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// DefaultServerTimeouts are safe defaults when YAML omits app.http.
func DefaultServerTimeouts() ServerTimeouts {
	return ServerTimeouts{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// LoadServerTimeouts reads app.http.* duration keys; zero/missing values fall
// back to defaults.
func LoadServerTimeouts() ServerTimeouts {
	d := DefaultServerTimeouts()
	if v := viper.GetDuration("app.http.read_header_timeout"); v > 0 {
		d.ReadHeaderTimeout = v
	}
	if v := viper.GetDuration("app.http.read_timeout"); v > 0 {
		d.ReadTimeout = v
	}
	if v := viper.GetDuration("app.http.write_timeout"); v > 0 {
		d.WriteTimeout = v
	}
	if v := viper.GetDuration("app.http.idle_timeout"); v > 0 {
		d.IdleTimeout = v
	}
	return d
}

// RateLimitConfig controls optional API rate limiting.
type RateLimitConfig struct {
	Enabled bool
	Rate    int64
}

// LoadRateLimitConfig reads app.ratelimit.* from viper (defaults Rate=100).
func LoadRateLimitConfig() RateLimitConfig {
	rate := viper.GetInt64("app.ratelimit.rate")
	if rate <= 0 {
		rate = 100
	}
	return RateLimitConfig{
		Enabled: viper.GetBool("app.ratelimit.enabled"),
		Rate:    rate,
	}
}

// ReadinessConfig controls /ready dependency probes.
type ReadinessConfig struct {
	Enabled    bool
	CheckMySQL bool
	CheckRedis bool
}

// LoadReadinessConfig reads app.readiness.* from viper.
func LoadReadinessConfig() ReadinessConfig {
	return ReadinessConfig{
		Enabled:    viper.GetBool("app.readiness.enabled"),
		CheckMySQL: viper.GetBool("app.readiness.mysql"),
		CheckRedis: viper.GetBool("app.readiness.redis"),
	}
}

// AppConfig holds application-level settings injected via Wire.
type AppConfig struct {
	Port           int
	Env            string
	HTTP           ServerTimeouts
	MetricsEnabled bool
	PprofEnabled   bool
	Readiness      ReadinessConfig
	RateLimit      RateLimitConfig
}

// LoadAppConfig assembles AppConfig from viper, applying HTTP/port defaults.
func LoadAppConfig() AppConfig {
	cfg := AppConfig{
		Port:           viper.GetInt("app.port"),
		Env:            viper.GetString("app.env"),
		HTTP:           LoadServerTimeouts(),
		MetricsEnabled: viper.GetBool("app.metrics.enabled"),
		PprofEnabled:   viper.GetBool("app.pprof.enabled"),
		Readiness:      LoadReadinessConfig(),
		RateLimit:      LoadRateLimitConfig(),
	}
	if cfg.Port <= 0 {
		cfg.Port = 8080
	}
	return cfg
}

// LoadAccessWhitelistConfig reads access_whitelist.* from viper into the
// kit middleware config struct.
func LoadAccessWhitelistConfig() middleware.AccessWhitelistConfig {
	return middleware.AccessWhitelistConfig{
		Enabled: viper.GetBool("access_whitelist.enabled"),
		CIDRs:   viper.GetStringSlice("access_whitelist.ip_list"),
	}
}
