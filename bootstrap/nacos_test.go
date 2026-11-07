package bootstrap

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	libNacos "github.com/andywang2011/go-nacos"
)

type fakeNacos struct {
	content     string
	downloadErr error
	watchErr    error

	mu       sync.Mutex
	dlCtx    context.Context
	watchCtx context.Context
	onChange func([]byte) error
}

func (f *fakeNacos) Download(ctx context.Context) (string, error) {
	f.mu.Lock()
	f.dlCtx = ctx
	f.mu.Unlock()
	return f.content, f.downloadErr
}

func (f *fakeNacos) Watch(ctx context.Context, onChange func([]byte) error) error {
	f.mu.Lock()
	f.watchCtx = ctx
	f.onChange = onChange
	f.mu.Unlock()
	return f.watchErr
}

type ctxKey struct{}

func nacosEnv(t *testing.T) (metaPath, confPath string) {
	t.Helper()
	dir := t.TempDir()
	metaPath = filepath.Join(dir, "nacos.yaml")
	if err := os.WriteFile(metaPath, []byte("ns1:\n  data_id: app.yaml\n  group: g\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	confPath = filepath.Join(dir, "conf.yaml")
	t.Setenv(libNacos.EnvNamespace, "ns1")
	t.Setenv(libNacos.EnvServer, "127.0.0.1")
	t.Setenv(libNacos.EnvPort, "8848")
	t.Setenv(libNacos.EnvUsername, "u")
	t.Setenv(libNacos.EnvPassword, "p")
	t.Setenv(libNacos.EnvSkip, "")
	return metaPath, confPath
}

func TestSyncFromNacosContext_PassesContextToWatch(t *testing.T) {
	meta, conf := nacosEnv(t)
	fake := &fakeNacos{content: "app:\n  env: development\n"}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxKey{}, "v"))
	defer cancel()

	var changeErr error
	changed := false
	err := SyncFromNacosContext(ctx, NacosOptions{
		ConfigPath:    conf,
		MetaPaths:     []string{meta},
		ClientFactory: func(libNacos.Config) (NacosClient, error) { return fake, nil },
		OnNacosChange: func(err error) { changed = true; changeErr = err },
	})
	if err != nil {
		t.Fatal(err)
	}
	if fake.dlCtx.Value(ctxKey{}) != "v" || fake.watchCtx.Value(ctxKey{}) != "v" {
		t.Fatal("ctx not propagated to Download/Watch")
	}
	cancel()
	if fake.watchCtx.Err() == nil {
		t.Fatal("watch ctx must be cancellable")
	}
	got, err := os.ReadFile(conf)
	if err != nil || string(got) != fake.content {
		t.Fatalf("conf = %q err=%v", got, err)
	}

	if err := fake.onChange([]byte("new: 1\n")); err != nil {
		t.Fatal(err)
	}
	if !changed || changeErr != nil {
		t.Fatalf("changed=%v err=%v", changed, changeErr)
	}
	got, _ = os.ReadFile(conf)
	if string(got) != "new: 1\n" {
		t.Fatalf("conf after change = %q", got)
	}
}

func TestSyncFromNacos_BackgroundContext(t *testing.T) {
	meta, conf := nacosEnv(t)
	fake := &fakeNacos{content: "a: 1\n"}
	err := SyncFromNacos(NacosOptions{
		ConfigPath:    conf,
		MetaPaths:     []string{meta},
		ClientFactory: func(libNacos.Config) (NacosClient, error) { return fake, nil },
	})
	if err != nil || fake.watchCtx == nil {
		t.Fatalf("err=%v watchCtx=%v", err, fake.watchCtx)
	}
}

func TestSyncFromNacosContext_NilContextDefaults(t *testing.T) {
	meta, conf := nacosEnv(t)
	fake := &fakeNacos{content: "a: 1\n"}
	//nolint:staticcheck // deliberately nil
	err := SyncFromNacosContext(nil, NacosOptions{
		ConfigPath:    conf,
		MetaPaths:     []string{meta},
		ClientFactory: func(libNacos.Config) (NacosClient, error) { return fake, nil },
	})
	if err != nil || fake.watchCtx == nil {
		t.Fatalf("err=%v", err)
	}
}

func TestSyncFromNacosContext_Errors(t *testing.T) {
	ctx := context.Background()
	t.Run("skip", func(t *testing.T) {
		t.Setenv(libNacos.EnvSkip, "Y")
		if err := SyncFromNacosContext(ctx, NacosOptions{}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("missing namespace", func(t *testing.T) {
		t.Setenv(libNacos.EnvSkip, "")
		t.Setenv(libNacos.EnvNamespace, "")
		if err := SyncFromNacosContext(ctx, NacosOptions{}); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("missing meta", func(t *testing.T) {
		nacosEnv(t)
		if err := SyncFromNacosContext(ctx, NacosOptions{MetaPaths: []string{"/nonexistent/x.yaml"}}); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("bad env config", func(t *testing.T) {
		meta, _ := nacosEnv(t)
		t.Setenv(libNacos.EnvPort, "")
		if err := SyncFromNacosContext(ctx, NacosOptions{MetaPaths: []string{meta}}); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("factory error", func(t *testing.T) {
		meta, conf := nacosEnv(t)
		boom := errors.New("factory boom")
		err := SyncFromNacosContext(ctx, NacosOptions{
			ConfigPath: conf, MetaPaths: []string{meta},
			ClientFactory: func(libNacos.Config) (NacosClient, error) { return nil, boom },
		})
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("download error", func(t *testing.T) {
		meta, conf := nacosEnv(t)
		boom := errors.New("dl boom")
		err := SyncFromNacosContext(ctx, NacosOptions{
			ConfigPath: conf, MetaPaths: []string{meta},
			ClientFactory: func(libNacos.Config) (NacosClient, error) { return &fakeNacos{downloadErr: boom}, nil },
		})
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("write error", func(t *testing.T) {
		meta, _ := nacosEnv(t)
		err := SyncFromNacosContext(ctx, NacosOptions{
			ConfigPath: "/nonexistent-dir/conf.yaml", MetaPaths: []string{meta},
			ClientFactory: func(libNacos.Config) (NacosClient, error) { return &fakeNacos{content: "a: 1"}, nil },
		})
		if err == nil {
			t.Fatal("expected write error")
		}
	})
	t.Run("watch error", func(t *testing.T) {
		meta, conf := nacosEnv(t)
		boom := errors.New("watch boom")
		err := SyncFromNacosContext(ctx, NacosOptions{
			ConfigPath: conf, MetaPaths: []string{meta},
			ClientFactory: func(libNacos.Config) (NacosClient, error) {
				return &fakeNacos{content: "a: 1", watchErr: boom}, nil
			},
		})
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestRunContext_SkipNacosLoadsViper(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "conf.yaml"), []byte("app:\n  env: development\n  port: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RunContext(context.Background(), Options{SkipNacos: true, ConfigDir: dir}); err != nil {
		t.Fatal(err)
	}
	if err := Run(Options{SkipNacos: true, ConfigDir: dir}); err != nil {
		t.Fatal(err)
	}
}

func TestRunContext_PropagatesFactoryAndContext(t *testing.T) {
	meta, conf := nacosEnv(t)
	dir := filepath.Dir(conf)
	fake := &fakeNacos{content: "app:\n  env: development\n  port: 1\n"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := RunContext(ctx, Options{
		ConfigDir:     dir,
		ConfigPath:    conf,
		MetaPaths:     []string{meta},
		ClientFactory: func(libNacos.Config) (NacosClient, error) { return fake, nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if fake.watchCtx == nil || fake.watchCtx.Err() == nil {
		t.Fatal("RunContext must pass its ctx to Watch")
	}
}

func TestRunContext_NacosErrorStopsRun(t *testing.T) {
	nacosEnv(t)
	err := RunContext(context.Background(), Options{MetaPaths: []string{"/nonexistent/nacos.yaml"}})
	if err == nil {
		t.Fatal("expected nacos sync error")
	}
}
