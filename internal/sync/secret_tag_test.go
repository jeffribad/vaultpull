package sync

import (
	"testing"
)

func TestApplyTagFilter_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"KEY": "val", "_tags": "internal"}
	cfg := TagConfig{Enabled: false}
	out := ApplyTagFilter(secrets, cfg)
	if out["KEY"] != "val" {
		t.Error("expected KEY to be present")
	}
}

func TestApplyTagFilter_NoTagsKey_NoRequired(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	cfg := TagConfig{Enabled: true}
	out := ApplyTagFilter(secrets, cfg)
	if out["KEY"] != "val" {
		t.Error("expected KEY to pass through")
	}
}

func TestApplyTagFilter_NoTagsKey_WithRequired_Excluded(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	cfg := TagConfig{Enabled: true, RequiredTags: []string{"internal"}}
	out := ApplyTagFilter(secrets, cfg)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}

func TestApplyTagFilter_RequiredTag_Present(t *testing.T) {
	secrets := map[string]string{"KEY": "val", "_tags": "internal,prod"}
	cfg := TagConfig{Enabled: true, RequiredTags: []string{"internal"}}
	out := ApplyTagFilter(secrets, cfg)
	if out["KEY"] != "val" {
		t.Error("expected KEY to be present")
	}
	if _, ok := out["_tags"]; ok {
		t.Error("_tags should be stripped from output")
	}
}

func TestApplyTagFilter_RequiredTag_Missing(t *testing.T) {
	secrets := map[string]string{"KEY": "val", "_tags": "staging"}
	cfg := TagConfig{Enabled: true, RequiredTags: []string{"prod"}}
	out := ApplyTagFilter(secrets, cfg)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}

func TestApplyTagFilter_ExcludeTag_Matches(t *testing.T) {
	secrets := map[string]string{"KEY": "val", "_tags": "deprecated,internal"}
	cfg := TagConfig{Enabled: true, ExcludeTags: []string{"deprecated"}}
	out := ApplyTagFilter(secrets, cfg)
	if len(out) != 0 {
		t.Errorf("expected empty result due to excluded tag, got %v", out)
	}
}

func TestApplyTagFilter_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{"KEY": "val", "_tags": "Internal"}
	cfg := TagConfig{Enabled: true, RequiredTags: []string{"internal"}}
	out := ApplyTagFilter(secrets, cfg)
	if out["KEY"] != "val" {
		t.Error("expected case-insensitive match")
	}
}
