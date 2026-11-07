package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefaultMaxRawBodyBytes is the request body size limit (10 MiB) applied by
// CacheRawBody.
const DefaultMaxRawBodyBytes int64 = 10 << 20

// CacheRawBody caches the raw request body under gin.BodyBytesKey (for
// ShouldBindBodyWith) AND restores Request.Body so downstream readers such as
// ShouldBindJSON can still consume the original stream.
//
// The body is capped at DefaultMaxRawBodyBytes; larger requests are rejected
// with 413 Request Entity Too Large.
//
// Restoring the body fixes a bug in the naive implementation that only cached
// the bytes but left Request.Body drained, breaking ShouldBindJSON.
func CacheRawBody() gin.HandlerFunc {
	return CacheRawBodyWithLimit(DefaultMaxRawBodyBytes)
}

// CacheRawBodyWithLimit is like CacheRawBody with a caller-chosen limit in
// bytes. A non-positive maxBytes falls back to DefaultMaxRawBodyBytes.
func CacheRawBodyWithLimit(maxBytes int64) gin.HandlerFunc {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxRawBodyBytes
	}
	return func(c *gin.Context) {
		// Only handle methods that carry a request body.
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "request body too large",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		// Cache for ShouldBindBodyWith.
		c.Set(gin.BodyBytesKey, raw)
		// Restore Request.Body so ShouldBindJSON and friends still read it.
		c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

		c.Next()
	}
}
