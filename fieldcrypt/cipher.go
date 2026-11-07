package fieldcrypt

// Cipher provides field encryption, decryption, and blind indexing.
type Cipher interface {
	// Encrypt encodes plaintext as a version-prefixed Base64 string.
	Encrypt(plain []byte) (string, error)
	// Decrypt decodes a string produced by Encrypt.
	Decrypt(enc string) ([]byte, error)
	// BlindIndex returns an HMAC-SHA256 hex blind index for a normalized value.
	BlindIndex(normalized string) string
}
