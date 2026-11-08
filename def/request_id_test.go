package def

import (
	"context"
	"testing"
)

func TestRequestIDRoundTrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "rid-1")
	if got := RequestIDFromContext(ctx); got != "rid-1" {
		t.Fatalf("got %q", got)
	}
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Fatalf("empty ctx = %q", got)
	}
}
