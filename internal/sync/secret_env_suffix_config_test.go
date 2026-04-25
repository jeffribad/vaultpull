package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuffixAddConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_ADD_ENABLED", "")
	t.Setenv("VAULTPULL_SUFFIX_ADD_SUFFIX", "")
	t.Setenv("VAULTPULL_SUFFIX_ADD_KEYS", "")
	cfg := SuffixAddConfigFromEnv()
	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.Suffix)
	assert.Empty(t, cfg.Keys)
}

func TestSuffixAddConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_SUFFIX_ADD_SUFFIX", "_LEGACY")
	t.Setenv("VAULTPULL_SUFFIX_ADD_KEYS", "")
	cfg := SuffixAddConfigFromEnv()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "_LEGACY", cfg.Suffix)
}

func TestSuffixAddConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_ADD_ENABLED", "1")
	t.Setenv("VAULTPULL_SUFFIX_ADD_SUFFIX", "_X")
	t.Setenv("VAULTPULL_SUFFIX_ADD_KEYS", "")
	cfg := SuffixAddConfigFromEnv()
	assert.True(t, cfg.Enabled)
}

func TestSuffixAddConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_SUFFIX_ADD_SUFFIX", "_V2")
	t.Setenv("VAULTPULL_SUFFIX_ADD_KEYS", "DB_HOST , DB_PORT ")
	cfg := SuffixAddConfigFromEnv()
	assert.Equal(t, []string{"DB_HOST", "DB_PORT"}, cfg.Keys)
}
