package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnforceProtect_Disabled_ReturnsNil(t *testing.T) {
	cfg := ProtectConfig{Enabled: false, Keys: []string{"DB_PASS"}}
	existing := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}
	errs := EnforceProtect(cfg, existing, incoming)
	assert.Nil(t, errs)
}

func TestEnforceProtect_NoKeys_ReturnsNil(t *testing.T) {
	cfg := ProtectConfig{Enabled: true, Keys: []string{}}
	errs := EnforceProtect(cfg, map[string]string{"X": "1"}, map[string]string{"X": "2"})
	assert.Nil(t, errs)
}

func TestEnforceProtect_NoExistingValue_Allowed(t *testing.T) {
	cfg := ProtectConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	existing := map[string]string{}
	incoming := map[string]string{"DB_PASS": "new"}
	errs := EnforceProtect(cfg, existing, incoming)
	assert.Nil(t, errs)
}

func TestEnforceProtect_SameValue_Allowed(t *testing.T) {
	cfg := ProtectConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	existing := map[string]string{"DB_PASS": "same"}
	incoming := map[string]string{"DB_PASS": "same"}
	errs := EnforceProtect(cfg, existing, incoming)
	assert.Nil(t, errs)
}

func TestEnforceProtect_ChangedValue_ReturnsError(t *testing.T) {
	cfg := ProtectConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	existing := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}
	errs := EnforceProtect(cfg, existing, incoming)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "DB_PASS")
}

func TestEnforceProtect_MultipleViolations(t *testing.T) {
	cfg := ProtectConfig{Enabled: true, Keys: []string{"KEY_A", "KEY_B"}}
	existing := map[string]string{"KEY_A": "a", "KEY_B": "b"}
	incoming := map[string]string{"KEY_A": "a2", "KEY_B": "b2"}
	errs := EnforceProtect(cfg, existing, incoming)
	assert.Len(t, errs, 2)
}

func TestProtectConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_PROTECT_ENABLED", "")
	t.Setenv("VAULTPULL_PROTECT_KEYS", "")
	cfg := ProtectConfigFromEnv()
	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.Keys)
}

func TestProtectConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_PROTECT_ENABLED", "true")
	t.Setenv("VAULTPULL_PROTECT_KEYS", "SECRET_KEY, API_TOKEN ")
	cfg := ProtectConfigFromEnv()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, []string{"SECRET_KEY", "API_TOKEN"}, cfg.Keys)
}
