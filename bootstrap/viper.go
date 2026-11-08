package bootstrap

import (
	"log"

	"github.com/andywang2011/go-app-kit/def"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// DefaultRelationPath is the local application config file updated from Nacos.
const DefaultRelationPath = "./config/conf.yaml"

// ValidateAppEnv checks app.env against allowed values.
func ValidateAppEnv(env string) error {
	switch env {
	case def.EnvDevelopment, def.EnvTesting, def.EnvStaging, def.EnvProduction:
		return nil
	default:
		return &envError{env: env}
	}
}

type envError struct{ env string }

func (e *envError) Error() string { return "app.env is not correct: " + e.env }

// SetupViperFromPath loads conf.yaml, watches changes, validates app.env,
// applies APP_* overrides, and rejects placeholder secrets for
// testing/staging/production (which expect Nacos-provided config).
// Optional secretKeys are forwarded to ValidateSecrets; empty uses DefaultProdSecretKeys.
func SetupViperFromPath(configDir string, secretKeys ...string) error {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	ApplyEnvOverrides()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	env := viper.GetString("app.env")
	if err := ValidateAppEnv(env); err != nil {
		return err
	}
	return ValidateSecrets(env, secretKeys...)
}
