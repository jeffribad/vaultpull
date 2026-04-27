// Package sync provides the ApplySplit transformation for vaultpull.
//
// # Secret Split
//
// ApplySplit reads a single "packed" secret value from Vault and expands it
// into multiple key/value pairs within the secrets map. This is useful when
// a Vault path stores a compact, delimited representation of several related
// configuration values.
//
// # Configuration (environment variables)
//
//	VAULTPULL_SPLIT_ENABLED   - enable the feature (true/1)
//	VAULTPULL_SPLIT_SOURCE    - name of the source secret key to split
//	VAULTPULL_SPLIT_DELIMITER - segment delimiter (default: ",")
//	VAULTPULL_SPLIT_SEPARATOR - key/value separator within each segment (default: "=")
//
// # Example
//
// Given a Vault secret:
//
//	PACKED_CONFIG = "DB_HOST=db.internal,DB_PORT=5432"
//
// With VAULTPULL_SPLIT_SOURCE=PACKED_CONFIG, ApplySplit will inject:
//
//	DB_HOST = "db.internal"
//	DB_PORT = "5432"
//
// The original source key is preserved in the output map.
package sync
