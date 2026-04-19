package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_READONLY_ENABLED", "")
	t.Setenv("VAULTPULL_READONLY_KEYS", "")
	t.Setenv("VAULTPULL_READONLY_FAIL_FAST", "")

	cfg := ReadOnlyConfigFromEnv()
	assert.False(t, cfg.Enabled)
	assert.False(t, cfg.FailFast)
	assert.Empty(t, cfg.Keys)
}

func TestReadOnlyConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_READONLY_ENABLED", "true")
	t.Setenv("VAULTPULL_READONLY_KEYS", "DB_PASSWORD, API_SECRET ")
	t.Setenv("VAULTPULL_READONLY_FAIL_FAST", "1")

	cfg := ReadOnlyConfigFromEnv()
	assert.True(t, cfg.Enabled)
	assert.True(t, cfg.FailFast)
	require.Len(t, cfg.Keys, 2)
	assert.Equal(t, "DB_PASSWORD", cfg.Keys[0])
	assert.Equal(t, "API_SECRET", cfg.Keys[1])
}

func TestEnforceReadOnly_Disabled_ReturnsNil(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: false, Keys: []string{"SECRET"}}
	secrets := map[string]string{"SECRET": "value"}
	errs := EnforceReadOnly(cfg, secrets)
	assert.Nil(t, errs)
}

func TestEnforceReadOnly_NoKeys_ReturnsNil(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: true, Keys: []string{}}
	secrets := map[string]string{"SECRET": "value"}
	errs := EnforceReadOnly(cfg, secrets)
	assert.Nil(t, errs)
}

func TestEnforceReadOnly_NoViolations_ReturnsNil(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: true, Keys: []string{"LOCKED_KEY"}}
	secrets := map[string]string{"OTHER_KEY": "value"}
	errs := EnforceReadOnly(cfg, secrets)
	assert.Nil(t, errs)
}

func TestEnforceReadOnly_DetectsViolation(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: true, Keys: []string{"DB_PASSWORD"}}
	secrets := map[string]string{"DB_PASSWORD": "secret", "APP_NAME": "myapp"}
	errs := EnforceReadOnly(cfg, secrets)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "DB_PASSWORD")
}

func TestEnforceReadOnly_FailFast_StopsEarly(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: true, FailFast: true, Keys: []string{"KEY_A", "KEY_B"}}
	secrets := map[string]string{"KEY_A": "v1", "KEY_B": "v2"}
	errs := EnforceReadOnly(cfg, secrets)
	assert.Len(t, errs, 1)
}

func TestEnforceReadOnly_CaseInsensitiveMatch(t *testing.T) {
	cfg := ReadOnlyConfig{Enabled: true, Keys: []string{"DB_PASSWORD"}}
	secrets := map[string]string{"db_password": "secret"}
	errs := EnforceReadOnly(cfg, secrets)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "db_password")
}
