package bootstrap

import (
	"testing"

	"github.com/andywang2011/go-app-kit/def"

	"github.com/spf13/viper"
)

func TestIsPlaceholderSecret(t *testing.T) {
	for _, v := range []string{"", "xxx", "123456", "local-dev-aes-16", "local"} {
		if !IsPlaceholderSecret(v) {
			t.Errorf("IsPlaceholderSecret(%q) = false", v)
		}
	}
	if IsPlaceholderSecret("real-secret-value-here") {
		t.Fatal("expected real secret to pass")
	}
}

func TestValidateSecretsSkipsDevelopment(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	if err := ValidateSecrets(def.EnvDevelopment); err != nil {
		t.Fatalf("development: %v", err)
	}
}

func TestValidateSecretsTestingRejectsLocalDevPlaceholders(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set("mysql.default.password", "local")
	viper.Set("jwt.secret", "ok-jwt-secret-value")
	viper.Set("app.aes.disabled", true)
	if err := ValidateSecrets(def.EnvTesting); err == nil {
		t.Fatal("expected error for local mysql password in testing")
	}
}

func TestValidateSecretsProdOK(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set("mysql.default.password", "ok-mysql-password")
	viper.Set("jwt.secret", "ok-jwt-secret-value")
	viper.Set("security.blind_index_key", "ok-blind-index-key-32bytes-long!!")
	viper.Set("app.aes.disabled", true)
	if err := ValidateSecrets(def.EnvProduction); err != nil {
		t.Fatalf("expected ok when aes disabled: %v", err)
	}
}

func TestValidateSecretsCustomKeys(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set("mysql.default.password", "ok-mysql-password")
	viper.Set("jwt.secret", "ok-jwt-secret-value")
	viper.Set("grpc.auth.token", "ok-grpc-token-value")
	keys := []string{"mysql.default.password", "jwt.secret", "grpc.auth.token"}
	if err := ValidateSecrets(def.EnvProduction, keys...); err != nil {
		t.Fatalf("custom keys: %v", err)
	}
	viper.Set("grpc.auth.token", "xxx")
	if err := ValidateSecrets(def.EnvProduction, keys...); err == nil {
		t.Fatal("expected placeholder grpc.auth.token to fail")
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	t.Setenv("APP_JWT_SECRET", "from-env")
	ApplyEnvOverrides()
	if got := viper.GetString("jwt.secret"); got != "from-env" {
		t.Fatalf("jwt.secret = %q, want from-env", got)
	}
}
