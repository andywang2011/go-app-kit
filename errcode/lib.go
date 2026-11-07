package errcode

// Framework component error codes (emitted by kit middleware).
var (
	// LibNoRoute is returned when no route matches the request.
	LibNoRoute = &ErrCode{Code: 1001, Key: "web.lib_no_route"}
	// LibRateLimit is returned when the rate limiter rejects a request.
	LibRateLimit = &ErrCode{Code: 1002, Key: "web.lib_rate_limit"}
	// LibNotInWhitelist is returned when the client IP is not whitelisted.
	LibNotInWhitelist = &ErrCode{Code: 1003, Key: "web.lib_not_in_whitelist"}
	// LibBodyTooBig is returned when the request body exceeds the log limit.
	LibBodyTooBig = &ErrCode{Code: 1005, Key: "web.lib_body_too_big"}
)
