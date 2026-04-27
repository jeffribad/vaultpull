package sync

import (
	"testing"
)

func TestApplySplit_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"MULTI": "A=1,B=2"}
	cfg := SplitConfig{Enabled: false, SourceKey: "MULTI", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplySplit_NoSourceKey_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"MULTI": "A=1,B=2"}
	cfg := SplitConfig{Enabled: true, SourceKey: "", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplySplit_MissingSourceKey_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"OTHER": "value"}
	cfg := SplitConfig{Enabled: true, SourceKey: "MULTI", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OTHER"]; !ok {
		t.Error("expected OTHER key to be preserved")
	}
}

func TestApplySplit_SplitsIntoMultipleKeys(t *testing.T) {
	secrets := map[string]string{"MULTI": "DB_HOST=localhost,DB_PORT=5432"}
	cfg := SplitConfig{Enabled: true, SourceKey: "MULTI", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
	if out["MULTI"] != "DB_HOST=localhost,DB_PORT=5432" {
		t.Error("expected source key to be preserved")
	}
}

func TestApplySplit_CaseInsensitiveSourceKey(t *testing.T) {
	secrets := map[string]string{"multi_secret": "X=10,Y=20"}
	cfg := SplitConfig{Enabled: true, SourceKey: "MULTI_SECRET", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "10" {
		t.Errorf("expected X=10, got %q", out["X"])
	}
}

func TestApplySplit_SkipsMalformedSegments(t *testing.T) {
	secrets := map[string]string{"MULTI": "GOOD=ok,badsegment,ALSO=fine"}
	cfg := SplitConfig{Enabled: true, SourceKey: "MULTI", Delimiter: ",", Separator: "="}
	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GOOD"] != "ok" {
		t.Errorf("expected GOOD=ok")
	}
	if out["ALSO"] != "fine" {
		t.Errorf("expected ALSO=fine")
	}
	if _, ok := out["badsegment"]; ok {
		t.Error("malformed segment should not produce a key")
	}
}

func TestApplySplit_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"MULTI": "K=v"}
	cfg := SplitConfig{Enabled: true, SourceKey: "MULTI", Delimiter: ",", Separator: "="}
	_, _ = ApplySplit(secrets, cfg)
	if len(secrets) != 1 {
		t.Error("input map was mutated")
	}
}
