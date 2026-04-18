package sync

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func randomKey(t *testing.T) string {
	t.Helper()
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func TestEncryptSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	cfg := EncryptConfig{Enabled: false}
	out, err := EncryptSecrets(secrets, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected original value, got %q", out["FOO"])
	}
}

func TestEncryptSecrets_MissingKey_ReturnsError(t *testing.T) {
	cfg := EncryptConfig{Enabled: true, Key: ""}
	_, err := EncryptSecrets(map[string]string{"A": "b"}, cfg)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestEncryptSecrets_InvalidKey_ReturnsError(t *testing.T) {
	cfg := EncryptConfig{Enabled: true, Key: "not-valid-base64!!!"}
	_, err := EncryptSecrets(map[string]string{"A": "b"}, cfg)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestEncryptSecrets_EncryptsValues(t *testing.T) {
	key := randomKey(t)
	cfg := EncryptConfig{Enabled: true, Key: key}
	secrets := map[string]string{"DB_PASS": "supersecret", "API_KEY": "abc123"}
	out, err := EncryptSecrets(secrets, cfg)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range out {
		if v == secrets[k] {
			t.Errorf("key %s: expected encrypted value, got plaintext", k)
		}
		if _, err := base64.StdEncoding.DecodeString(v); err != nil {
			t.Errorf("key %s: encrypted value is not valid base64", k)
		}
	}
}

func TestEncryptSecrets_DifferentNonceEachCall(t *testing.T) {
	key := randomKey(t)
	cfg := EncryptConfig{Enabled: true, Key: key}
	secrets := map[string]string{"X": "value"}
	out1, _ := EncryptSecrets(secrets, cfg)
	out2, _ := EncryptSecrets(secrets, cfg)
	if out1["X"] == out2["X"] {
		t.Error("expected different ciphertext due to random nonce")
	}
}
