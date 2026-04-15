package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheKey_Deterministic(t *testing.T) {
	k1 := CacheKey("secret/app", "backend")
	k2 := CacheKey("secret/app", "backend")
	if k1 != k2 {
		t.Fatalf("expected same key, got %q and %q", k1, k2)
	}
}

func TestCacheKey_DifferentInputs(t *testing.T) {
	k1 := CacheKey("secret/app", "backend")
	k2 := CacheKey("secret/app", "frontend")
	if k1 == k2 {
		t.Fatal("expected different keys for different roles")
	}
}

func TestCacheEntry_IsExpired(t *testing.T) {
	old := &CacheEntry{FetchedAt: time.Now().Add(-10 * time.Minute)}
	if !old.IsExpired(5 * time.Minute) {
		t.Error("expected entry to be expired")
	}

	fresh := &CacheEntry{FetchedAt: time.Now()}
	if fresh.IsExpired(5 * time.Minute) {
		t.Error("expected entry to be fresh")
	}
}

func TestNewSecretCache_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	_, err := NewSecretCache(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected cache dir to be created")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache, err := NewSecretCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	entry := &CacheEntry{
		Path:      "secret/app",
		Secrets:   map[string]string{"DB_URL": "postgres://localhost"},
		FetchedAt: time.Now().Truncate(time.Second),
		Version:   3,
	}
	key := CacheKey("secret/app", "backend")

	if err := cache.Set(key, entry); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := cache.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil entry")
	}
	if got.Version != 3 {
		t.Errorf("expected version 3, got %d", got.Version)
	}
	if got.Secrets["DB_URL"] != "postgres://localhost" {
		t.Errorf("unexpected secret value: %q", got.Secrets["DB_URL"])
	}
}

func TestCache_GetMissing_ReturnsNil(t *testing.T) {
	cache, _ := NewSecretCache(t.TempDir())
	got, err := cache.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatal("expected nil for missing key")
	}
}

func TestCache_Invalidate(t *testing.T) {
	cache, _ := NewSecretCache(t.TempDir())
	key := CacheKey("secret/app", "ops")

	_ = cache.Set(key, &CacheEntry{Path: "secret/app", FetchedAt: time.Now()})
	if err := cache.Invalidate(key); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}
	got, _ := cache.Get(key)
	if got != nil {
		t.Fatal("expected nil after invalidation")
	}
	// Invalidating non-existent key should not error
	if err := cache.Invalidate(key); err != nil {
		t.Fatalf("double invalidate should not error: %v", err)
	}
}
