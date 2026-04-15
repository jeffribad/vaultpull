package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry holds a cached snapshot of secrets fetched from Vault.
type CacheEntry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	Version   int               `json:"version"`
}

// IsExpired reports whether the cache entry is older than ttl.
func (e *CacheEntry) IsExpired(ttl time.Duration) bool {
	return time.Since(e.FetchedAt) > ttl
}

// CacheKey returns a stable filename-safe key derived from the secret path and role.
func CacheKey(secretPath, role string) string {
	h := sha256.Sum256([]byte(secretPath + "|" + role))
	return hex.EncodeToString(h[:])[:16]
}

// SecretCache manages on-disk caching of Vault secrets.
type SecretCache struct {
	dir string
}

// NewSecretCache creates a SecretCache that stores entries under dir.
func NewSecretCache(dir string) (*SecretCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &SecretCache{dir: dir}, nil
}

// Get retrieves a cache entry for the given key. Returns nil if not found.
func (c *SecretCache) Get(key string) (*CacheEntry, error) {
	data, err := os.ReadFile(c.entryPath(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return &entry, nil
}

// Set writes a cache entry for the given key.
func (c *SecretCache) Set(key string, entry *CacheEntry) error {
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	if err := os.WriteFile(c.entryPath(key), data, 0600); err != nil {
		return fmt.Errorf("cache: write: %w", err)
	}
	return nil
}

// Invalidate removes the cache entry for the given key.
func (c *SecretCache) Invalidate(key string) error {
	err := os.Remove(c.entryPath(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *SecretCache) entryPath(key string) string {
	return filepath.Join(c.dir, key+".json")
}
