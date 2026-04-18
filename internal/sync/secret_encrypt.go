package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strconv"
)

type EncryptConfig struct {
	Enabled bool
	Key     string // 32-byte base64-encoded AES-256 key
}

func EncryptConfigFromEnv() EncryptConfig {
	enabled := false
	if v := os.Getenv("VAULTPULL_ENCRYPT_OUTPUT"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}
	return EncryptConfig{
		Enabled: enabled,
		Key:     os.Getenv("VAULTPULL_ENCRYPT_KEY"),
	}
}

// EncryptSecrets encrypts each secret value using AES-256-GCM.
// Returns a new map with base64-encoded ciphertext values.
func EncryptSecrets(secrets map[string]string, cfg EncryptConfig) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}
	if cfg.Key == "" {
		return nil, errors.New("encrypt: VAULTPULL_ENCRYPT_KEY is required when encryption is enabled")
	}
	keyBytes, err := base64.StdEncoding.DecodeString(cfg.Key)
	if err != nil || len(keyBytes) != 32 {
		return nil, errors.New("encrypt: key must be a base64-encoded 32-byte AES-256 key")
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		result[k] = base64.StdEncoding.EncodeToString(ciphertext)
	}
	return result, nil
}
