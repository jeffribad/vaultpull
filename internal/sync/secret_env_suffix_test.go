package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddKeySuffix_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := SuffixAddConfig{Enabled: false, Suffix: "_NEW"}
	out := AddKeySuffix(secrets, cfg)
	assert.Equal(t, secrets, out)
}

func TestAddKeySuffix_EmptySuffix_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := SuffixAddConfig{Enabled: true, Suffix: ""}
	out := AddKeySuffix(secrets, cfg)
	assert.Equal(t, secrets, out)
}

func TestAddKeySuffix_AllKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	cfg := SuffixAddConfig{Enabled: true, Suffix: "_V2"}
	out := AddKeySuffix(secrets, cfg)
	assert.Equal(t, "localhost", out["DB_HOST_V2"])
	assert.Equal(t, "5432", out["DB_PORT_V2"])
	assert.Len(t, out, 2)
}

func TestAddKeySuffix_SpecificKeys_OnlySuffixesListed(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
	cfg := SuffixAddConfig{Enabled: true, Suffix: "_PROD", Keys: []string{"DB_HOST"}}
	out := AddKeySuffix(secrets, cfg)
	assert.Equal(t, "localhost", out["DB_HOST_PROD"])
	assert.Equal(t, "5432", out["DB_PORT"])
	assert.Equal(t, "secret", out["API_KEY"])
	assert.Len(t, out, 3)
}

func TestAddKeySuffix_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc123"}
	cfg := SuffixAddConfig{Enabled: true, Suffix: "_COPY"}
	_ = AddKeySuffix(secrets, cfg)
	_, original := secrets["TOKEN"]
	assert.True(t, original, "original map should be unchanged")
	_, mutated := secrets["TOKEN_COPY"]
	assert.False(t, mutated, "original map must not be mutated")
}

func TestAddKeySuffix_CaseInsensitiveKeyMatch(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "API_KEY": "xyz"}
	cfg := SuffixAddConfig{Enabled: true, Suffix: "_NEW", Keys: []string{"DB_HOST"}}
	out := AddKeySuffix(secrets, cfg)
	assert.Equal(t, "localhost", out["db_host_NEW"])
	assert.Equal(t, "xyz", out["API_KEY"])
}
