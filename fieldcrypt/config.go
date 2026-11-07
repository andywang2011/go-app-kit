package fieldcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	aesKeySize = 32
	nonceSize  = 12
	versionSep = ":"
)

// Config holds field-encryption settings.
type Config struct {
	// Keys maps version label to 32-byte AES-256 key material.
	Keys map[string][]byte
	// ActiveKey is the version label used for encryption, e.g. "v1".
	ActiveKey string
	// BlindIndexKey is the 32-byte HMAC key for blind indexes.
	BlindIndexKey []byte
}

type aesCipher struct {
	keys          map[string][]byte
	activeKey     string
	blindIndexKey []byte
}

// New validates cfg and returns a Cipher implementation.
func New(cfg Config) (Cipher, error) {
	if cfg.ActiveKey == "" {
		return nil, errors.New("fieldcrypt: ActiveKey is required")
	}
	if len(cfg.Keys) == 0 {
		return nil, errors.New("fieldcrypt: Keys must not be empty")
	}
	active, ok := cfg.Keys[cfg.ActiveKey]
	if !ok {
		return nil, fmt.Errorf("fieldcrypt: ActiveKey %q not found in Keys", cfg.ActiveKey)
	}
	if len(active) != aesKeySize {
		return nil, fmt.Errorf("fieldcrypt: key %q must be %d bytes", cfg.ActiveKey, aesKeySize)
	}
	keys := make(map[string][]byte, len(cfg.Keys))
	for ver, key := range cfg.Keys {
		if len(key) != aesKeySize {
			return nil, fmt.Errorf("fieldcrypt: key %q must be %d bytes", ver, aesKeySize)
		}
		k := make([]byte, aesKeySize)
		copy(k, key)
		keys[ver] = k
	}
	if len(cfg.BlindIndexKey) != aesKeySize {
		return nil, fmt.Errorf("fieldcrypt: BlindIndexKey must be %d bytes", aesKeySize)
	}
	bik := make([]byte, aesKeySize)
	copy(bik, cfg.BlindIndexKey)

	return &aesCipher{
		keys:          keys,
		activeKey:     cfg.ActiveKey,
		blindIndexKey: bik,
	}, nil
}

// Encrypt implements Cipher.
func (c *aesCipher) Encrypt(plain []byte) (string, error) {
	key, ok := c.keys[c.activeKey]
	if !ok {
		return "", fmt.Errorf("fieldcrypt: active key %q not found", c.activeKey)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, plain, nil)
	payload := append(nonce, ciphertext...)
	return c.activeKey + versionSep + base64.RawStdEncoding.EncodeToString(payload), nil
}

func parseVersionAndPayload(enc string) (version string, payload []byte, err error) {
	idx := strings.Index(enc, versionSep)
	if idx <= 0 {
		return "", nil, errors.New("fieldcrypt: invalid ciphertext format")
	}
	version = enc[:idx]
	b64 := enc[idx+1:]
	payload, err = base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		return "", nil, fmt.Errorf("fieldcrypt: decode ciphertext: %w", err)
	}
	if len(payload) < nonceSize+1 {
		return "", nil, errors.New("fieldcrypt: ciphertext too short")
	}
	return version, payload, nil
}

// Decrypt implements Cipher.
func (c *aesCipher) Decrypt(enc string) ([]byte, error) {
	version, payload, err := parseVersionAndPayload(enc)
	if err != nil {
		return nil, err
	}
	key, ok := c.keys[version]
	if !ok {
		return nil, fmt.Errorf("fieldcrypt: unknown key version %q", version)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("fieldcrypt: decrypt: %w", err)
	}
	return plain, nil
}

// BlindIndex implements Cipher.
func (c *aesCipher) BlindIndex(normalized string) string {
	mac := hmac.New(sha256.New, c.blindIndexKey)
	_, _ = mac.Write([]byte(normalized))
	return hex.EncodeToString(mac.Sum(nil))
}
