package middleware

import (
	"bytes"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/andywang2011/go-app-kit/core"
	"github.com/andywang2011/go-app-kit/def"
	"github.com/andywang2011/go-app-kit/errcode"

	"github.com/andywang2011/go-logging"
	"github.com/gin-gonic/gin"
)

const bodyLimitKB = 5000

// maxLogRespBodyBytes caps how much of a response body is buffered for logging.
const maxLogRespBodyBytes = 64 << 10

const (
	ctxKeyLogDisablePrintBody = "disable_print_body" // whether to print body
	ctxKeyLogEncryptStatus    = "encrypt_status"     // log encryption status
)

var defaultLogRedactor logging.Redactor = logging.NewDefaultRedactor()

// defaultSkipBodyPaths are endpoints whose (large/binary) request bodies are not
// logged. Callers may override via SetSkipBodyPaths at startup.
var defaultSkipBodyPaths = []string{
	"/admin/upload/image",
	"/admin/upload/file",
	"/common/upload/img",
	"/common/upload/base64img",
	"/api/common/upload/img",
	"/api/common/upload/base64img",
}

var (
	skipBodyMu    sync.RWMutex
	skipBodyPaths = buildSkipSet(defaultSkipBodyPaths)
)

func buildSkipSet(paths []string) map[string]struct{} {
	out := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		out[p] = struct{}{}
	}
	return out
}

// SetSkipBodyPaths overrides the set of request paths whose bodies are not
// logged. Passing an empty slice restores the built-in defaults. This keeps the
// kit free of viper/global config: the caller decides the list at startup.
func SetSkipBodyPaths(paths []string) {
	skipBodyMu.Lock()
	defer skipBodyMu.Unlock()
	if len(paths) == 0 {
		skipBodyPaths = buildSkipSet(defaultSkipBodyPaths)
		return
	}
	skipBodyPaths = buildSkipSet(paths)
}

func shouldSkipBody(path string) bool {
	skipBodyMu.RLock()
	defer skipBodyMu.RUnlock()
	_, ok := skipBodyPaths[path]
	return ok
}

// BodyLogWriter captures the response body for logging while still writing it
// to the client. At most limit bytes are buffered; the client always receives
// the full response.
type BodyLogWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	limit     int
	truncated bool
}

func (w *BodyLogWriter) capture(b []byte) {
	if w.limit <= 0 {
		w.limit = maxLogRespBodyBytes
	}
	remain := w.limit - w.body.Len()
	if remain <= 0 {
		w.truncated = w.truncated || len(b) > 0
		return
	}
	if len(b) > remain {
		w.body.Write(b[:remain])
		w.truncated = true
		return
	}
	w.body.Write(b)
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.capture(b)
	return w.ResponseWriter.Write(b)
}

// WriteString captures and forwards string writes (e.g. gin's c.String).
func (w *BodyLogWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

// redactQueryForLog flattens query values and redacts sensitive keys with the
// default redactor.
func redactQueryForLog(q url.Values) map[string]any {
	return defaultLogRedactor.RedactMap(logging.SpreadMaps(q))
}

// responseLogBody returns the redacted response body to log. Truncated bodies
// are never logged because partial JSON cannot be reliably redacted.
func responseLogBody(w *BodyLogWriter) []byte {
	if w.truncated {
		return []byte("[response body truncated for logging]")
	}
	respStr := w.body.String()
	if strings.TrimSpace(respStr) == "" {
		return w.body.Bytes()
	}
	return []byte(defaultLogRedactor.RedactJSONBody(respStr))
}

// Logging prints the request/response (with field redaction). AES encryption of
// the log payload can still be disabled via the logging options.
func Logging(c *core.Context) {
	c.Set(ctxKeyLogEncryptStatus, logging.EncryptStatusDefault)
	logRequest(c)
}

// LoggerDisableEncrypt disables log payload encryption; combine with Logging.
func LoggerDisableEncrypt(c *core.Context) {
	c.Set(ctxKeyLogEncryptStatus, logging.EncryptStatusDisable)
}

// LoggerEnableEncrypt enables log payload encryption; combine with Logging.
func LoggerEnableEncrypt(c *core.Context) {
	c.Set(ctxKeyLogEncryptStatus, logging.EncryptStatusEnable)
}

// LoggerDisablePrintBody disables body printing; combine with Logging.
func LoggerDisablePrintBody(c *core.Context) {
	c.Set(ctxKeyLogDisablePrintBody, true)
}

func logRequest(c *core.Context) {
	start := time.Now()
	raw, _ := c.Get(gin.BodyBytesKey)
	rawBytes, _ := raw.([]byte)
	if rawBytes == nil {
		rawBytes = []byte{}
	}

	w := &BodyLogWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil), limit: maxLogRespBodyBytes}
	c.Writer = w

	if c.Logger == nil {
		c.Logger = logging.New()
	}
	c.Logger.SetRequestID(c.GetString(def.XRequestID))
	c.Logger.Setv1(c.GetString("v1"))

	rawKB := len(rawBytes) / 1024
	skipBody := shouldSkipBody(c.Request.URL.Path)
	if rawKB > bodyLimitKB && !skipBody {
		preview := string(rawBytes[:min(1024, len(rawBytes))])
		c.Info("API request and response", defaultLogRedactor.RedactJSONBody(preview), nil)
		c.FailWithErrCode(errcode.LibBodyTooBig, nil)
		return
	}

	c.Next()

	bodyStr := string(rawBytes)
	if rawKB > bodyLimitKB {
		bodyStr = string(rawBytes[:min(1024, len(rawBytes))])
	}
	if skipBody || c.GetBool(ctxKeyLogDisablePrintBody) {
		bodyStr = ""
	} else {
		bodyStr = defaultLogRedactor.RedactJSONBody(bodyStr)
	}

	request := map[string]any{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"query":  redactQueryForLog(c.Request.URL.Query()),
		"header": defaultLogRedactor.RedactHeaders(c.Request.Header),
		"body":   bodyStr,
	}

	respBody := responseLogBody(w)

	opts := logging.AttrOption{
		StartTime:        &start,
		EncryptStatus:    c.GetInt(ctxKeyLogEncryptStatus),
		DisablePrintBody: c.GetBool(ctxKeyLogDisablePrintBody),
	}
	c.Info("API request and response", request, respBody, opts)
}
