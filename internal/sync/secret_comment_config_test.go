package sync

import (
	"os"
	"testing"
)

func TestCommentConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_COMMENTS_ENABLED")
	os.Unsetenv("VAULTPULL_COMMENTS_PREFIX")
	cfg := CommentConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Prefix != "# -- %s --" {
		t.Errorf("unexpected default prefix: %s", cfg.Prefix)
	}
}

func TestCommentConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_COMMENTS_ENABLED", "true")
	cfg := CommentConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestCommentConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_COMMENTS_ENABLED", "1")
	cfg := CommentConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestCommentConfigFromEnv_CustomPrefix(t *testing.T) {
	t.Setenv("VAULTPULL_COMMENTS_ENABLED", "true")
	t.Setenv("VAULTPULL_COMMENTS_PREFIX", "## %s")
	cfg := CommentConfigFromEnv()
	if cfg.Prefix != "## %s" {
		t.Errorf("unexpected prefix: %s", cfg.Prefix)
	}
}
