package fieldcrypt

import (
	"strings"
	"unicode"
)

// NormalizeEmail trims space and lowercases an email address.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizePhoneE164 keeps digits and an optional leading '+' for E.164-style phones.
func NormalizePhoneE164(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range phone {
		if r == '+' && i == 0 {
			b.WriteRune('+')
			continue
		}
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
