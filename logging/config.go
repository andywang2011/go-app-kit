package logging

import (
	goLog "github.com/andywang2011/go-logging"
	"github.com/spf13/viper"
)

// NewConfigFromViper builds go-logging config from viper keys app.log and app.aes.
func NewConfigFromViper() goLog.Config {
	maxSize := viper.GetInt("app.log.file_max_size")
	if maxSize <= 0 {
		maxSize = 500
	}
	backupCount := viper.GetInt("app.log.file_backup_count")
	if backupCount <= 0 {
		backupCount = 2
	}
	return goLog.Config{
		Output:          viper.GetString("app.log.output"),
		AESDisabled:     viper.GetBool("app.aes.disabled"),
		// AESKey must stay empty: go-logging v0.2+ rejects AES field masking.
		AESFields:       viper.GetStringSlice("app.aes.fields"),
		FileMaxSize:     maxSize,
		FileBackupCount: backupCount,
		LogTopic:        viper.GetString("app.log.topic"),
	}
}
