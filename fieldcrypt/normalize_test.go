package fieldcrypt

import "testing"

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"  User@Example.COM  ", "user@example.com"},
		{"already@lowercase.com", "already@lowercase.com"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := NormalizeEmail(tt.in); got != tt.want {
			t.Errorf("NormalizeEmail(%q) = %q want %q", tt.in, got, tt.want)
		}
	}
}

func TestNormalizePhoneE164(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{" +86 138 0013 8000 ", "+8613800138000"},
		{"138-0013-8000", "13800138000"},
		{"+1 (415) 555-0100", "+14155550100"},
		{"", ""},
		{"  ", ""},
	}
	for _, tt := range tests {
		if got := NormalizePhoneE164(tt.in); got != tt.want {
			t.Errorf("NormalizePhoneE164(%q) = %q want %q", tt.in, got, tt.want)
		}
	}
}
