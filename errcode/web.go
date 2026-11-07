package errcode

// Shared web error codes used by both front-end (miniapp) and admin services.
var (
	// WebNetworkBusy is a generic transient failure code.
	WebNetworkBusy = &ErrCode{Code: 40000, Key: "web.network_busy"}
	// WebAuthInfoFailed is returned when authentication info cannot be resolved.
	WebAuthInfoFailed = &ErrCode{Code: 40001, Key: "web.auth_info_failed"}
	// WebCacheTokenFailed is returned when a token cannot be cached.
	WebCacheTokenFailed = &ErrCode{Code: 40002, Key: "web.cache_token_failed"}
	// WebTokenExpired is returned when the session token is missing/expired.
	WebTokenExpired = &ErrCode{Code: 40003, Key: "web.token_expired"}
	// WebUserInfoFailed is returned when user info cannot be resolved.
	WebUserInfoFailed = &ErrCode{Code: 40004, Key: "web.user_info_failed"}
	// WebMobileFailed is returned when a mobile number cannot be resolved.
	WebMobileFailed = &ErrCode{Code: 40006, Key: "web.mobile_failed"}
	// WebInvalidParam is returned when request parameters fail validation.
	WebInvalidParam = &ErrCode{Code: 40100, Key: "web.invalid_param"}
	// WebInvalidSign is returned when a request signature is invalid.
	WebInvalidSign = &ErrCode{Code: 40101, Key: "web.invalid_sign"}
	// TextGreenCheckFail is returned when text content moderation fails.
	TextGreenCheckFail = &ErrCode{Code: 40103, Key: "web.text_green_check_fail"}
	// ImageGreenCheckFail is returned when image content moderation fails.
	ImageGreenCheckFail = &ErrCode{Code: 40104, Key: "web.image_green_check_fail"}
	// FileTypeError is returned when an uploaded file type is not allowed.
	FileTypeError = &ErrCode{Code: 40105, Key: "web.file_type_error"}
	// UploadLimitError is returned when an upload exceeds its size limit.
	UploadLimitError = &ErrCode{Code: 40106, Key: "web.upload_limit_error"}
	// WebRoleNoPermission is returned when the role lacks permission.
	WebRoleNoPermission = &ErrCode{Code: 40107, Key: "web.role_no_permission"}
	// WebUserForbidden is returned when the user is forbidden.
	WebUserForbidden = &ErrCode{Code: 40108, Key: "web.user_forbidden"}
	// WebNoRecord is returned when no matching record exists.
	WebNoRecord = &ErrCode{Code: 40111, Key: "web.no_record"}
)
