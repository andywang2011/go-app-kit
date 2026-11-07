package def_test

import (
	"testing"

	"github.com/andywang2011/go-app-kit/def"
)

func TestEnvironmentConstants(t *testing.T) {
	cases := map[string]string{
		"production":  def.EnvProduction,
		"development": def.EnvDevelopment,
		"testing":     def.EnvTesting,
		"staging":     def.EnvStaging,
	}
	for want, got := range cases {
		if got != want {
			t.Fatalf("env constant: want %q, got %q", want, got)
		}
	}
}

func TestLanguageConstants(t *testing.T) {
	if def.Language != "Language" {
		t.Fatalf("Language header: %q", def.Language)
	}
	for _, pair := range []struct {
		code string
		want string
	}{
		{def.ZhLanguage, "zh"},
		{def.EnLanguage, "en"},
		{def.JaLanguage, "ja"},
	} {
		if pair.code != pair.want {
			t.Fatalf("language code: want %q, got %q", pair.want, pair.code)
		}
	}
}

func TestHeaderConstants(t *testing.T) {
	if def.AuthorizationToken != "Authorization-Token" {
		t.Fatalf("AuthorizationToken: %q", def.AuthorizationToken)
	}
	if def.XRequestID != "X-Request-Id" {
		t.Fatalf("XRequestID: %q", def.XRequestID)
	}
}
