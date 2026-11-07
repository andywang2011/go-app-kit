package appkit

import "github.com/gin-gonic/gin"

// Module is a self-contained feature that registers its routes onto a Gin
// router group. A process composes several Modules behind its router.
//
// Implementations must be safe to call Register once during startup.
type Module interface {
	// Name returns a stable identifier used for logging and diagnostics.
	Name() string
	// Enabled reports whether the module should be registered. A disabled
	// module is skipped entirely (no routes, no background work).
	Enabled() bool
	// Register mounts the module's routes onto the provided router group.
	Register(r *gin.RouterGroup)
}

// Closer is implemented by components that hold resources (connections, async
// writers, background goroutines) which must be released on shutdown.
//
// Close must be safe to call once; implementations should make it idempotent
// where practical. The process shutdown sequence calls Close on every Closer.
type Closer interface {
	// Close releases the component's resources, returning the first error.
	Close() error
}
