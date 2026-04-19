package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnforceRequired_Disabled_ReturnsNil(t *testing.T) {
	cfg := RequiredConfig{Enabled: false, Keys: []string{"DB_PASSWORD"}}
	secrets := map[string]string{}
	violations := EnforceRequired(cfg, secrets)
	assert.Nil(t, violations)
}

func TestEnforceRequired_NoKeys_ReturnsNil(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{}}
	secrets := map[string]string{"DB_PASSWORD": "secret"}
	violations := EnforceRequired(cfg, secrets)
	assert.Nil(t, violations)
}

func TestEnforceRequired_AllPresent_ReturnsNil(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{"DB_HOST", "DB_PORT"}}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	violations := EnforceRequired(cfg, secrets)
	assert.Empty(t, violations)
}

func TestEnforceRequired_MissingKey_ReturnsViolation(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{"DB_HOST", "DB_PASS"}}
	secrets := map[string]string{"DB_HOST": "localhost"}
	violations := EnforceRequired(cfg, secrets)
	assert.Len(t, violations, 1)
	assert.Contains(t, violations[0], "DB_PASS")
	assert.Contains(t, violations[0], "missing")
}

func TestEnforceRequired_EmptyValue_ReturnsViolation(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{"API_KEY"}}
	secrets := map[string]string{"API_KEY": "   "}
	violations := EnforceRequired(cfg, secrets)
	assert.Len(t, violations, 1)
	assert.Contains(t, violations[0], "empty")
}

func TestEnforceRequired_FailFast_StopsAtFirst(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{"KEY_A", "KEY_B", "KEY_C"}, FailFast: true}
	secrets := map[string]string{}
	violations := EnforceRequired(cfg, secrets)
	assert.Len(t, violations, 1)
}

func TestEnforceRequired_MultipleViolations_AllReported(t *testing.T) {
	cfg := RequiredConfig{Enabled: true, Keys: []string{"KEY_A", "KEY_B"}, FailFast: false}
	secrets := map[string]string{}
	violations := EnforceRequired(cfg, secrets)
	assert.Len(t, violations, 2)
}
