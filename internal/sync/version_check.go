package sync

import (
	"fmt"
	"time"
)

// VersionInfo holds information about a secret's version at sync time.
type VersionInfo struct {
	Path        string
	Version     int
	CreatedTime string
	Destroyed   bool
	CheckedAt   time.Time
}

// VersionChecker is capable of retrieving version metadata for a secret.
type VersionChecker interface {
	GetSecretVersion(path string, version int) (*SecretVersionResult, error)
}

// SecretVersionResult is a local representation of vault.SecretVersion
// to avoid coupling the sync package directly to the vault package.
type SecretVersionResult struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// CheckVersion fetches version metadata for the given path and returns
// a VersionInfo. Returns an error if the secret is destroyed.
func CheckVersion(checker VersionChecker, path string) (*VersionInfo, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	sv, err := checker.GetSecretVersion(path, 0)
	if err != nil {
		return nil, fmt.Errorf("version check failed for %q: %w", path, err)
	}

	if sv.Destroyed {
		return nil, fmt.Errorf("secret at %q (version %d) has been destroyed", path, sv.Version)
	}

	return &VersionInfo{
		Path:        path,
		Version:     sv.Version,
		CreatedTime: sv.CreatedTime,
		Destroyed:   sv.Destroyed,
		CheckedAt:   time.Now().UTC(),
	}, nil
}
