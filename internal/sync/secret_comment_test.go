package sync

import (
	"strings"
	"testing"
)

func TestInjectComments_Disabled_ReturnsNil(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := CommentConfig{Enabled: false, Prefix: "# -- %s --"}
	lines := InjectComments(secrets, cfg)
	if lines != nil {
		t.Errorf("expected nil, got %v", lines)
	}
}

func TestInjectComments_EmptyMap_ReturnsNil(t *testing.T) {
	cfg := CommentConfig{Enabled: true, Prefix: "# -- %s --"}
	lines := InjectComments(map[string]string{}, cfg)
	if lines != nil {
		t.Error("expected nil for empty map")
	}
}

func TestInjectComments_SingleGroup(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	cfg := CommentConfig{Enabled: true, Prefix: "# -- %s --"}
	lines := InjectComments(secrets, cfg)
	if len(lines) == 0 {
		t.Fatal("expected non-empty lines")
	}
	if lines[0] != "# -- DB --" {
		t.Errorf("expected comment header, got: %s", lines[0])
	}
}

func TestInjectComments_MultipleGroups_InsertsBlankLine(t *testing.T) {
	secrets := map[string]string{
		"APP_NAME": "vaultpull",
		"DB_HOST":  "localhost",
	}
	cfg := CommentConfig{Enabled: true, Prefix: "# -- %s --"}
	lines := InjectComments(secrets, cfg)
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "\n\n") {
		t.Error("expected blank line between groups")
	}
	if !strings.Contains(joined, "# -- APP --") {
		t.Error("expected APP group comment")
	}
	if !strings.Contains(joined, "# -- DB --") {
		t.Error("expected DB group comment")
	}
}

func TestInjectComments_CustomPrefix_NoPlaceholder(t *testing.T) {
	secrets := map[string]string{"REDIS_URL": "redis://localhost"}
	cfg := CommentConfig{Enabled: true, Prefix: "##"}
	lines := InjectComments(secrets, cfg)
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "## REDIS") {
		t.Errorf("unexpected first line: %v", lines)
	}
}
