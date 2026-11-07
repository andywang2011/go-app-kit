package logging

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestNewConfigFromViper(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, "conf.yaml")
	content := `
app:
  log:
    output: std
    topic: test-topic
    file_max_size: 100
    file_backup_count: 3
  aes:
    disabled: true
    key: ""
    fields: []
`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	viper.Reset()
	defer viper.Reset()
	viper.SetConfigFile(confPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	cfg := NewConfigFromViper()
	if cfg.Output != "std" {
		t.Errorf("Output = %q", cfg.Output)
	}
	if !cfg.AESDisabled {
		t.Error("AESDisabled should be true")
	}
	if cfg.FileMaxSize != 100 || cfg.FileBackupCount != 3 {
		t.Errorf("FileMaxSize=%d FileBackupCount=%d", cfg.FileMaxSize, cfg.FileBackupCount)
	}
	if cfg.LogTopic != "test-topic" {
		t.Errorf("LogTopic = %q", cfg.LogTopic)
	}
}
