package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// Encryption key must be 32 bytes (base64 or raw).
// Recommended env: CONFIG_ENCRYPTION_KEY (base64-encoded 32 bytes).
func loadKey() ([]byte, error) {
	raw := os.Getenv("CONFIG_ENCRYPTION_KEY")
	if raw == "" {
		return nil, errors.New("CONFIG_ENCRYPTION_KEY is not set")
	}

	// Try base64 decode first
	if b, err := base64.StdEncoding.DecodeString(raw); err == nil && len(b) == 32 {
		return b, nil
	}

	// Fallback to raw bytes
	if len(raw) == 32 {
		return []byte(raw), nil
	}

	return nil, fmt.Errorf("CONFIG_ENCRYPTION_KEY must be 32 bytes raw or base64(32 bytes)")
}

func EncryptString(plain string) (string, error) {
	key, err := loadKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nil, nonce, []byte(plain), nil)
	payload := append(nonce, ct...)
	return "enc:" + base64.StdEncoding.EncodeToString(payload), nil
}

func DecryptString(enc string) (string, error) {
	if len(enc) < 4 || enc[:4] != "enc:" {
		return enc, nil // treat as plaintext for backward compatibility
	}
	key, err := loadKey()
	if err != nil {
		return "", err
	}
	b, err := base64.StdEncoding.DecodeString(enc[4:])
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(b) < gcm.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}
	nonce := b[:gcm.NonceSize()]
	ct := b[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}
