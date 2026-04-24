// Package sync provides secret blacklisting support for vaultpull.
//
// The blacklist feature allows operators to explicitly block specific secret
// keys from being written to the output .env file, regardless of what Vault
// returns. This is useful for preventing sensitive internal keys (e.g.,
// INTERNAL_TOKEN, DEBUG_PASSWORD) from leaking into developer environments.
//
// Configuration is driven by environment variables:
//
//	VAULTPULL_BLACKLIST_ENABLED=true
//	VAULTPULL_BLACKLIST_KEYS=SECRET_KEY,INTERNAL_TOKEN,ADMIN_PASS
//
// Key matching is case-insensitive, so "secret_key" will block "SECRET_KEY".
//
// The blacklist is applied after all other transformations (rename, alias,
// transform, etc.) and before the final write step, ensuring that no
// pipeline stage can accidentally re-introduce a blocked key.
package sync
