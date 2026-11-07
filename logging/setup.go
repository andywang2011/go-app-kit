package logging

import (
	"fmt"
	"sync"

	goLog "github.com/andywang2011/go-logging"
)

var (
	setupMu   sync.Mutex
	setupDone bool
	// setupFn initializes go-logging; replaceable in tests.
	setupFn = goLog.SetupFromConfig
)

// Setup reads logging config from viper and initializes go-logging; call once from main.
// It returns an error instead of exiting the process, so callers decide how to react.
// After a successful call further calls are no-ops; after a failure Setup may be retried.
func Setup() error {
	setupMu.Lock()
	defer setupMu.Unlock()
	if setupDone {
		return nil
	}
	if err := setupFn(NewConfigFromViper()); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	setupDone = true
	return nil
}

// Close 排空异步队列并关闭输出，进程退出前调用（可重复调用）
func Close() error {
	return goLog.Close()
}
