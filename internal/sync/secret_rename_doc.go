// Package sync provides secret synchronisation logic for vaultpull.
//
// # Secret Rename
//
// The rename feature allows callers to map Vault secret keys to different
// names in the output .env file.  This is useful when the naming convention
// used inside Vault differs from what an application expects.
//
// Configuration is driven by the VAULTPULL_RENAME_KEYS environment variable,
// which accepts a comma-separated list of OLD_KEY:NEW_KEY pairs:
//
//	VAULTPULL_RENAME_KEYS=DB_PASS:DATABASE_PASSWORD,API_TOKEN:APP_API_TOKEN
//
// The rename step should be applied after role-based filtering and secret
// key filtering so that include/exclude rules reference the original Vault
// key names.
//
// Example pipeline:
//
//	secrets  → FilterByRole → ApplySecretFilter → ApplyRenames → writer
package sync
