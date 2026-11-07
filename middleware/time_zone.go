package middleware

import (
	"time"

	"github.com/andywang2011/go-app-kit/core"
	"github.com/andywang2011/go-app-kit/errcode"
)

// SetTimeZone sets the request-scoped time zone used by handlers for localized
// time formatting. It fails the request if the zone database is unavailable.
func SetTimeZone(c *core.Context) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		c.FailWithErrCode(errcode.WebNetworkBusy, err.Error())
		return
	}

	c.TimeLocal = loc
	c.Next()
}
