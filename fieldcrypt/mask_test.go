package fieldcrypt

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"user@example.com", "u***@example.com"},
		{"  alice@corp.io ", "a***@corp.io"},
		{"invalid", "***"},
		{"@nodomain", "***"},
		{"", "***"},
	}
	for _, tt := range tests {
		if got := MaskEmail(tt.in); got != tt.want {
			t.Errorf("MaskEmail(%q) = %q want %q", tt.in, got, tt.want)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"+8613800138000", "+*********8000"},
		{"13800138000", "*******8000"},
		{"1234", "****"},
		{"12", "**"},
		{"", ""},
		{"abc", "***"},
	}
	for _, tt := range tests {
		if got := MaskPhone(tt.in); got != tt.want {
			t.Errorf("MaskPhone(%q) = %q want %q", tt.in, got, tt.want)
		}
	}
}
