package def

import "context"

type requestIDKey struct{}

// WithRequestID stores the per-request id on ctx for outbound clients.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestIDFromContext returns the id set by WithRequestID, or "".
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
