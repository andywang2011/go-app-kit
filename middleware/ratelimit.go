package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/andywang2011/go-app-kit/errcode"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// ErrInvalidRate is returned (wrapped) by RateLimit when rate is not positive.
var ErrInvalidRate = errors.New("middleware: invalid rate limit")

// RateLimit returns a token-bucket rate limiter allowing up to rate requests
// per second (burst = rate). It returns an error wrapping ErrInvalidRate when
// rate <= 0.
func RateLimit(rate int64) (gin.HandlerFunc, error) {
	if rate <= 0 {
		return nil, fmt.Errorf("%w: rate must be positive, got %d", ErrInvalidRate, rate)
	}
	interval := time.Duration(int64(time.Second) / rate)
	if interval <= 0 {
		interval = time.Nanosecond // rate above 1e9/s: ratelimit panics on a zero interval
	}
	bucket := ratelimit.NewBucket(interval, rate)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": errcode.LibRateLimit.Code,
				"msg":  errcode.LibRateLimit.String(),
				"data": nil,
			})
			return
		}
		c.Next()
	}, nil
}
