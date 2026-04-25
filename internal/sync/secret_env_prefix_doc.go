// Package sync provides the secret_env_prefix feature for vaultpull.
//
// # Key Prefix Addition
//
// AddKeyPrefix prepends a configurable string prefix to secret keys pulled
// from Vault before they are written to a .env file. This is useful when
// multiple services share a single Vault path but require namespaced keys
// in their respective environment files.
//
// # Configuration (environment variables)
//
//	VAULTPULL_PREFIX_ADD_ENABLED  — enable the feature (true/1)
//	VAULTPULL_PREFIX_ADD_VALUE    — prefix string to prepend (e.g. "APP_")
//	VAULTPULL_PREFIX_ADD_KEYS     — optional comma-separated list of keys to
//	                                prefix; if empty, all keys are prefixed
//
// # Example
//
// Given secrets {"HOST": "db", "PORT": "5432"} and prefix "DB_", the result
// is {"DB_HOST": "db", "DB_PORT": "5432"}.
//
// When VAULTPULL_PREFIX_ADD_KEYS is set to "HOST", only HOST is prefixed:
// {"DB_HOST": "db", "PORT": "5432"}.
package sync
