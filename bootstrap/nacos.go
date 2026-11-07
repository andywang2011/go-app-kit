package bootstrap

import (
	"context"
	"fmt"
	"os"

	libNacos "github.com/andywang2011/go-nacos"
)

// NacosClient is the subset of the go-nacos client used by bootstrap (injectable for tests).
type NacosClient interface {
	Download(ctx context.Context) (string, error)
	Watch(ctx context.Context, onChange func(content []byte) error) error
}

// NacosClientFactory builds a NacosClient from Config.
type NacosClientFactory func(cfg libNacos.Config) (NacosClient, error)

func defaultNacosClientFactory(cfg libNacos.Config) (NacosClient, error) {
	return libNacos.NewClient(cfg)
}

// SyncFromNacos downloads and watches remote config when not skipped. The watch
// lives for the whole process; use SyncFromNacosContext for a cancellable watch.
func SyncFromNacos(opts NacosOptions) error {
	return SyncFromNacosContext(context.Background(), opts)
}

// SyncFromNacosContext downloads remote config and registers a watch that is
// cancelled when ctx is done.
func SyncFromNacosContext(ctx context.Context, opts NacosOptions) error {
	if libNacos.ShouldSkip(opts.SkipNacosFlag) {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	namespace := libNacos.NamespaceFromEnv()
	if namespace == "" {
		return fmt.Errorf("bootstrap: %s is required", libNacos.EnvNamespace)
	}
	metaPath, err := firstExistingFile(opts.MetaPaths)
	if err != nil {
		return err
	}
	cfg, err := libNacos.LoadConfigFromEnvAndMeta(metaPath, namespace)
	if err != nil {
		return fmt.Errorf("bootstrap: load nacos config: %w", err)
	}
	factory := opts.ClientFactory
	if factory == nil {
		factory = defaultNacosClientFactory
	}
	client, err := factory(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap: nacos client: %w", err)
	}
	content, err := client.Download(ctx)
	if err != nil {
		return fmt.Errorf("bootstrap: nacos download: %w", err)
	}
	if err := libNacos.WriteConfigFile(opts.ConfigPath, []byte(content)); err != nil {
		return fmt.Errorf("bootstrap: write config: %w", err)
	}
	return client.Watch(ctx, func(data []byte) error {
		err := libNacos.WriteConfigFile(opts.ConfigPath, data)
		if opts.OnNacosChange != nil {
			opts.OnNacosChange(err)
		}
		return err
	})
}

func firstExistingFile(paths []string) (string, error) {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("bootstrap: nacos meta file not found in %v", paths)
}
