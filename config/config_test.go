package config

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestLoadServerTimeoutsDefaults(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	if got := LoadServerTimeouts(); got != DefaultServerTimeouts() {
		t.Fatalf("got %+v", got)
	}
}

func TestLoadServerTimeoutsOverride(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader("app:\n  http:\n    read_timeout: 30s\n")); err != nil {
		t.Fatal(err)
	}
	if got := LoadServerTimeouts(); got.ReadTimeout != 30*time.Second {
		t.Fatalf("ReadTimeout = %v", got.ReadTimeout)
	}
}

func TestLoadAppConfigDefaults(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	cfg := LoadAppConfig()
	if cfg.Port != 8080 {
		t.Fatalf("Port = %d", cfg.Port)
	}
	if cfg.RateLimit.Rate != 100 {
		t.Fatalf("Rate = %d", cfg.RateLimit.Rate)
	}
}

func TestLoadAccessWhitelistConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("access_whitelist.enabled", true)
	viper.Set("access_whitelist.ip_list", []string{"127.0.0.1"})
	cfg := LoadAccessWhitelistConfig()
	if !cfg.Enabled || len(cfg.CIDRs) != 1 || cfg.CIDRs[0] != "127.0.0.1" {
		t.Fatalf("cfg = %+v", cfg)
	}
}
