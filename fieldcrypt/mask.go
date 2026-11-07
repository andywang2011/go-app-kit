package fieldcrypt

import (
	"strings"
	"unicode/utf8"
)

// MaskEmail masks an email address for API responses.
func MaskEmail(email string) string {
	email = strings.TrimSpace(email)
	at := strings.LastIndex(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return "***"
	}
	local := email[:at]
	domain := email[at+1:]
	if local == "" {
		return "***@" + domain
	}
	first, _ := utf8.DecodeRuneInString(local)
	return string(first) + "***@" + domain
}

// MaskPhone masks a phone number for API responses.
func MaskPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	prefix := ""
	digits := phone
	if strings.HasPrefix(phone, "+") {
		prefix = "+"
		digits = phone[1:]
	}
	var digitRunes []rune
	for _, r := range digits {
		if r >= '0' && r <= '9' {
			digitRunes = append(digitRunes, r)
		}
	}
	if len(digitRunes) == 0 {
		return "***"
	}
	if len(digitRunes) <= 4 {
		return prefix + strings.Repeat("*", len(digitRunes))
	}
	tail := string(digitRunes[len(digitRunes)-4:])
	masked := strings.Repeat("*", len(digitRunes)-4) + tail
	return prefix + masked
}
