package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"DB_PORT":     "5432",
	"APP_SECRET":  "topsecret",
	"APP_DEBUG":   "false",
	"REDIS_URL":   "redis://localhost",
}

var testPolicies = []RolePolicy{
	{
		Role:    "backend",
		Allowed: []string{"DB", "APP_SECRET"},
	},
	{
		Role:    "frontend",
		Allowed: []string{"APP_DEBUG"},
	},
	{
		Role:    "ops",
		Allowed: []string{"DB", "REDIS_URL", "APP"},
	},
}

func TestFilterByRole_Backend(t *testing.T) {
	result := FilterByRole(testSecrets, "backend", testPolicies)
	assert.Equal(t, "localhost", result["DB_HOST"])
	assert.Equal(t, "secret", result["DB_PASSWORD"])
	assert.Equal(t, "5432", result["DB_PORT"])
	assert.Equal(t, "topsecret", result["APP_SECRET"])
	assert.NotContains(t, result, "APP_DEBUG")
	assert.NotContains(t, result, "REDIS_URL")
}

func TestFilterByRole_Frontend(t *testing.T) {
	result := FilterByRole(testSecrets, "frontend", testPolicies)
	assert.Len(t, result, 1)
	assert.Equal(t, "false", result["APP_DEBUG"])
}

func TestFilterByRole_Ops(t *testing.T) {
	result := FilterByRole(testSecrets, "ops", testPolicies)
	assert.Contains(t, result, "DB_HOST")
	assert.Contains(t, result, "REDIS_URL")
	assert.Contains(t, result, "APP_SECRET")
	assert.Contains(t, result, "APP_DEBUG")
}

func TestFilterByRole_UnknownRole(t *testing.T) {
	// Unknown role should return all secrets
	result := FilterByRole(testSecrets, "unknown", testPolicies)
	assert.Len(t, result, len(testSecrets))
}

func TestFilterByRole_CaseInsensitive(t *testing.T) {
	result := FilterByRole(testSecrets, "BACKEND", testPolicies)
	assert.Contains(t, result, "DB_HOST")
}

func TestFilterByRole_EmptyPolicies(t *testing.T) {
	result := FilterByRole(testSecrets, "backend", []RolePolicy{})
	assert.Len(t, result, len(testSecrets))
}

func TestMatchesAny_ExactMatch(t *testing.T) {
	assert.True(t, matchesAny("APP_SECRET", []string{"APP_SECRET"}))
	assert.False(t, matchesAny("APP_SECRET", []string{"APP_DEBUG"}))
}

func TestMatchesAny_PrefixMatch(t *testing.T) {
	assert.True(t, matchesAny("DB_HOST", []string{"DB"}))
	assert.True(t, matchesAny("DB_PASSWORD", []string{"DB"}))
	assert.False(t, matchesAny("DATABASE_URL", []string{"DB"}))
}
