// Package bootstrap performs process startup before Wire: it optionally syncs
// remote Nacos configuration, loads the local conf.yaml into viper, applies
// APP_* environment overrides and validates app.env and production secrets.
//
// This package intentionally uses viper and reads environment variables: it is
// the documented startup exception (analogous to go-nacos). Request-path kit
// packages (core, errcode, middleware) never depend on it.
package bootstrap

import (
	"context"

	libNacos "github.com/andywang2011/go-nacos"
)

const (
	// DefaultConfigDir is the directory containing conf.yaml.
	DefaultConfigDir = "./config"
	// DefaultConfigPath is the local application config file path.
	DefaultConfigPath = "./config/conf.yaml"
)

// Options controls process bootstrap before Wire.
type Options struct {
	SkipNacos     bool
	SkipNacosFlag string
	ConfigDir     string
	ConfigPath    string
	MetaPaths     []string
	OnNacosChange func(error)
	// ClientFactory overrides how the Nacos client is built (tests); nil uses go-nacos.
	ClientFactory NacosClientFactory
}

// DefaultOptions returns common bootstrap defaults.
func DefaultOptions() Options {
	return Options{
		ConfigDir:  DefaultConfigDir,
		ConfigPath: DefaultConfigPath,
		MetaPaths:  []string{"./config/nacos.yaml", "./nacos.yaml"},
	}
}

// Run syncs remote Nacos config (optional) and loads local YAML into viper.
// It is RunContext with context.Background (the Nacos watch is never cancelled).
func Run(opts Options) error {
	return RunContext(context.Background(), opts)
}

// RunContext is Run with a caller-controlled ctx: cancelling ctx stops the
// Nacos config watch and aborts an in-flight download.
func RunContext(ctx context.Context, opts Options) error {
	if opts.ConfigDir == "" {
		opts.ConfigDir = DefaultConfigDir
	}
	if opts.ConfigPath == "" {
		opts.ConfigPath = DefaultConfigPath
	}
	if opts.MetaPaths == nil {
		opts.MetaPaths = DefaultOptions().MetaPaths
	}
	if !opts.SkipNacos && !libNacos.ShouldSkip(opts.SkipNacosFlag) {
		if err := SyncFromNacosContext(ctx, NacosOptions{
			SkipNacosFlag: opts.SkipNacosFlag,
			ConfigPath:    opts.ConfigPath,
			MetaPaths:     opts.MetaPaths,
			OnNacosChange: opts.OnNacosChange,
			ClientFactory: opts.ClientFactory,
		}); err != nil {
			return err
		}
	}
	return SetupViperFromPath(opts.ConfigDir)
}

// NacosOptions configures remote Nacos sync.
type NacosOptions struct {
	SkipNacosFlag string
	ConfigPath    string
	MetaPaths     []string
	OnNacosChange func(error)
	// ClientFactory overrides how the Nacos client is built (tests); nil uses go-nacos.
	ClientFactory NacosClientFactory
}
