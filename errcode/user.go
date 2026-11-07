package errcode

// Shared user/account error codes (402xx) used by services built on the kit.
var (
	// WebUserLocked is returned when the account is locked.
	WebUserLocked = &ErrCode{Code: 40201, Key: "web.user_locked"}
	// WebPasswordExpire is returned when the password has expired.
	WebPasswordExpire = &ErrCode{Code: 40202, Key: "web.password_expire"}
	// WebInvalidPassword is returned when the password does not meet policy.
	WebInvalidPassword = &ErrCode{Code: 40203, Key: "web.invalid_password"}
	// WebOldPasswordInvalid is returned when the old password is wrong.
	WebOldPasswordInvalid = &ErrCode{Code: 40204, Key: "web.old_password_invalid"}
	// WebCaptchaInvalid is returned when a captcha is wrong or expired.
	WebCaptchaInvalid = &ErrCode{Code: 40205, Key: "web.captcha_invalid"}
	// WebUserNotExist is returned when the user does not exist.
	WebUserNotExist = &ErrCode{Code: 40206, Key: "web.user_not_exist"}
	// WebLoginFailed is returned for unknown user, bad password, and disabled
	// accounts to avoid account enumeration.
	WebLoginFailed = &ErrCode{Code: 40207, Key: "web.login_failed"}
)
