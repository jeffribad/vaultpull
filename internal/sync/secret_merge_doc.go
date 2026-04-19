// Package sync provides secret_merge functionality for vaultpull.
//
// # Secret Merge
//
// The merge feature allows combining secrets fetched from Vault with a local
// base map (e.g. previously parsed .env file). Three strategies are supported:
//
//   - vault-wins (default): Vault values overwrite any conflicting local values.
//     New local-only keys are preserved.
//
//   - local-wins: Local values take precedence over Vault for conflicting keys.
//     Vault-only keys are still included.
//
//   - union: Local values win for existing keys; Vault fills in keys that are
//     absent from the local base.
//
// # Configuration
//
// Set the following environment variables to enable and configure merging:
//
//	VAULTPULL_MERGE_ENABLED=true
//	VAULTPULL_MERGE_STRATEGY=vault-wins   # or local-wins, union
package sync
