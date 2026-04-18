package sync

import (
	"encoding/base64"
	"testing"
)

func TestEncryptSecrets_Integration_PipelineCompatible(t *testing.T) {
	key := randomKey(t)
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_TOKEN":   "tok_abc",
		"PLAIN_VAR":   "hello",
	}

	// Apply transforms first, then encrypt
	transformCfg := TransformConfig{Prefix: "APP_"}
	transformed := ApplyTransforms(secrets, transformCfg)

	encCfg := EncryptConfig{Enabled: true, Key: key}
	encrypted, err := EncryptSecrets(transformed, encCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(encrypted) != len(transformed) {
		t.Errorf("expected %d keys, got %d", len(transformed), len(encrypted))
	}
	for k, v := range encrypted {
		raw, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			t.Errorf("key %s: not valid base64: %v", k, err)
		}
		if len(raw) == 0 {
			t.Errorf("key %s: empty ciphertext", k)
		}
	}
}

func TestEncryptSecrets_Integration_DisabledPassthrough(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	cfg := EncryptConfig{Enabled: false}
	out, err := EncryptSecrets(secrets, cfg)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range secrets {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}
