package fieldcrypt

import (
	"bytes"
	"strings"
	"testing"
)

func testKey(seed byte) []byte {
	k := make([]byte, 32)
	for i := range k {
		k[i] = seed + byte(i)
	}
	return k
}

func mustNew(t *testing.T, cfg Config) Cipher {
	t.Helper()
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNewValidation(t *testing.T) {
	base := Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	}
	if _, err := New(Config{}); err == nil {
		t.Error("expected error for empty config")
	}
	shortKey := make([]byte, 16)
	cfg := base
	cfg.Keys = map[string][]byte{"v1": shortKey}
	if _, err := New(cfg); err == nil {
		t.Error("expected error for short encryption key")
	}
	cfg = base
	cfg.BlindIndexKey = shortKey
	if _, err := New(cfg); err == nil {
		t.Error("expected error for short blind index key")
	}
	cfg = base
	cfg.ActiveKey = "v99"
	if _, err := New(cfg); err == nil {
		t.Error("expected error for missing ActiveKey in Keys")
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	c := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	plain := []byte("sensitive-user-data-123")
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(enc, "v1:") {
		t.Fatalf("ciphertext prefix = %q", enc[:3])
	}
	got, err := c.Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("roundtrip = %q want %q", got, plain)
	}
	enc2, err := c.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if enc == enc2 {
		t.Error("expected random nonce to produce distinct ciphertexts")
	}
	got2, err := c.Decrypt(enc2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got2, plain) {
		t.Fatalf("second roundtrip = %q", got2)
	}
}

func TestDecryptWrongKeyFails(t *testing.T) {
	keyV1 := testKey(1)
	cEnc := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": keyV1},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	enc, err := cEnc.Encrypt([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	wrongKey := testKey(2)
	cDec := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": wrongKey},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	if _, err := cDec.Decrypt(enc); err == nil {
		t.Fatal("expected decrypt failure with wrong key")
	}
}

func TestVersionRotation(t *testing.T) {
	keyV1 := testKey(1)
	keyV2 := testKey(2)
	cEnc := mustNew(t, Config{
		Keys: map[string][]byte{
			"v1": keyV1,
			"v2": keyV2,
		},
		ActiveKey:     "v2",
		BlindIndexKey: testKey(9),
	})
	plain := []byte("rotate-me")
	enc, err := cEnc.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(enc, "v2:") {
		t.Fatalf("prefix = %q", enc[:3])
	}
	cDec := mustNew(t, Config{
		Keys: map[string][]byte{
			"v1": keyV1,
			"v2": keyV2,
		},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	got, err := cDec.Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q want %q", got, plain)
	}
	oldEnc, err := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": keyV1, "v2": keyV2},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	}).Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	gotOld, err := cDec.Decrypt(oldEnc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotOld, plain) {
		t.Fatalf("v1 ciphertext decrypt failed: %q", gotOld)
	}
}

func TestDecryptUnknownVersion(t *testing.T) {
	c := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	enc, err := c.Encrypt([]byte("x"))
	if err != nil {
		t.Fatal(err)
	}
	enc = "v9" + enc[2:]
	c2 := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	if _, err := c2.Decrypt(enc); err == nil {
		t.Fatal("expected error for unknown version")
	}
}

func TestBlindIndexDeterministic(t *testing.T) {
	c := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(9),
	})
	norm := NormalizeEmail("  User@Example.COM ")
	a := c.BlindIndex(norm)
	b := c.BlindIndex(norm)
	if a != b {
		t.Fatalf("blind index not deterministic: %q vs %q", a, b)
	}
	if len(a) != 64 {
		t.Fatalf("blind index length = %d want 64 hex chars", len(a))
	}
	c2 := mustNew(t, Config{
		Keys:          map[string][]byte{"v1": testKey(1)},
		ActiveKey:     "v1",
		BlindIndexKey: testKey(10),
	})
	if c2.BlindIndex(norm) == a {
		t.Error("expected different blind index with different HMAC key")
	}
}
