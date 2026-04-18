package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChecksumConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_CHECKSUM_ENABLED")
	os.Unsetenv("VAULTPULL_CHECKSUM_PATH")
	cfg := ChecksumConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.OutputPath != ".env.sha256" {
		t.Errorf("unexpected default path: %s", cfg.OutputPath)
	}
}

func TestChecksumConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_CHECKSUM_ENABLED", "true")
	t.Setenv("VAULTPULL_CHECKSUM_PATH", "/tmp/my.sha256")
	cfg := ChecksumConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.OutputPath != "/tmp/my.sha256" {
		t.Errorf("unexpected path: %s", cfg.OutputPath)
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	a := ComputeChecksum(secrets)
	b := ComputeChecksum(secrets)
	if a != b {
		t.Error("checksum not deterministic")
	}
	if len(a) != 64 {
		t.Errorf("expected 64-char hex, got %d", len(a))
	}
}

func TestComputeChecksum_DifferentInputs(t *testing.T) {
	a := ComputeChecksum(map[string]string{"KEY": "val1"})
	b := ComputeChecksum(map[string]string{"KEY": "val2"})
	if a == b {
		t.Error("expected different checksums for different values")
	}
}

func TestWriteChecksum_Disabled_ReturnsEmpty(t *testing.T) {
	cfg := ChecksumConfig{Enabled: false}
	sum, err := WriteChecksum(cfg, map[string]string{"A": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != "" {
		t.Errorf("expected empty sum, got %s", sum)
	}
}

func TestWriteChecksum_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.sha256")
	cfg := ChecksumConfig{Enabled: true, OutputPath: path}
	secrets := map[string]string{"DB_PASS": "secret"}

	sum, err := WriteChecksum(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	contents := strings.TrimSpace(string(data))
	if contents != sum {
		t.Errorf("file contents %q != returned sum %q", contents, sum)
	}
}
