package telemetry

import (
	"context"
	"testing"

	"github.com/spf13/viper"
)

func TestSetupDisabled(t *testing.T) {
	shutdown, err := Setup(context.Background(), Config{Enabled: false})
	if err != nil {
		t.Fatalf("Setup disabled: %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestConfigDefault(t *testing.T) {
	cfg := ConfigDefault()
	if cfg.Enabled {
		t.Fatal("expected disabled by default")
	}
	if cfg.SampleRatio != 1.0 {
		t.Fatalf("sample ratio=%v", cfg.SampleRatio)
	}
}

func TestConfigDefault_TLSByDefault(t *testing.T) {
	if ConfigDefault().Insecure {
		t.Fatal("expected Insecure=false by default (TLS)")
	}
}

func TestExporterOptions(t *testing.T) {
	cases := []struct {
		name    string
		cfg     Config
		wantErr bool
		wantLen int
	}{
		{"bare host TLS default", Config{Endpoint: "collector:4318"}, false, 1},
		{"bare host explicit insecure", Config{Endpoint: "collector:4318", Insecure: true}, false, 2},
		{"https URL", Config{Endpoint: "https://collector:4318"}, false, 1},
		{"http URL without Insecure", Config{Endpoint: "http://collector:4318"}, true, 0},
		{"http URL with Insecure", Config{Endpoint: "http://collector:4318", Insecure: true}, false, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := exporterOptions(tc.cfg)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
			}
			if len(opts) != tc.wantLen {
				t.Fatalf("len(opts)=%d want %d", len(opts), tc.wantLen)
			}
		})
	}
}

func TestSetup_PlaintextRequiresInsecure(t *testing.T) {
	_, err := Setup(context.Background(), Config{
		Enabled: true, Endpoint: "http://localhost:4318", ServiceName: "svc",
	})
	if err == nil {
		t.Fatal("expected error for plaintext endpoint without Insecure")
	}
}

func TestSetup_InsecureExplicitOK(t *testing.T) {
	shutdown, err := Setup(context.Background(), Config{
		Enabled: true, Endpoint: "localhost:4318", ServiceName: "svc", Insecure: true, SampleRatio: 2,
	})
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	_ = shutdown(context.Background())
}

func TestSetup_RequiresServiceName(t *testing.T) {
	if _, err := Setup(context.Background(), Config{Enabled: true}); err == nil {
		t.Fatal("expected service_name error")
	}
}

func TestLoadConfig_Insecure(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("app.otel.enabled", true)
	viper.Set("app.otel.insecure", true)
	viper.Set("app.otel.sample_ratio", 0.5)
	cfg := LoadConfig()
	if !cfg.Enabled || !cfg.Insecure || cfg.SampleRatio != 0.5 || cfg.Endpoint != "localhost:4318" {
		t.Fatalf("cfg=%+v", cfg)
	}
	viper.Reset()
	if LoadConfig().Insecure {
		t.Fatal("Insecure should default to false")
	}
}
