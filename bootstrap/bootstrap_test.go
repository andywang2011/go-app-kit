package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andywang2011/go-app-kit/def"

	"github.com/spf13/viper"
)

func TestValidateAppEnv(t *testing.T) {
	for _, env := range []string{def.EnvDevelopment, def.EnvTesting, def.EnvStaging, def.EnvProduction} {
		if err := ValidateAppEnv(env); err != nil {
			t.Errorf("ValidateAppEnv(%q) = %v", env, err)
		}
	}
	if err := ValidateAppEnv("bad"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSetupViperFromPath(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, "conf.yaml")
	validYaml := "app:\n  env: " + def.EnvDevelopment + "\n  port: 8080\n"
	if err := os.WriteFile(confPath, []byte(validYaml), 0644); err != nil {
		t.Fatal(err)
	}
	viper.Reset()
	defer viper.Reset()

	if err := SetupViperFromPath(dir); err != nil {
		t.Fatalf("SetupViperFromPath: %v", err)
	}
	if got := viper.GetString("app.env"); got != def.EnvDevelopment {
		t.Errorf("app.env = %q", got)
	}
}

func TestDefaultRelationPath(t *testing.T) {
	if DefaultRelationPath != "./config/conf.yaml" {
		t.Errorf("DefaultRelationPath = %q", DefaultRelationPath)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.ConfigDir == "" || len(opts.MetaPaths) == 0 {
		t.Fatalf("opts = %+v", opts)
	}
}
