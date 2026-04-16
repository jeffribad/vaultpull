package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressReporter_Advance_PrintsKey(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporterWithWriter(&buf, 3, false)
	p.Advance("DB_PASSWORD")
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "1/3") {
		t.Errorf("expected counter in output, got: %s", out)
	}
}

func TestProgressReporter_Quiet_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporterWithWriter(&buf, 3, true)
	p.Advance("SECRET_KEY")
	p.Done(2, 1)
	if buf.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %s", buf.String())
	}
}

func TestProgressReporter_Done_ShowsSummary(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporterWithWriter(&buf, 2, false)
	p.Done(2, 0)
	out := buf.String()
	if !strings.Contains(out, "2 written") {
		t.Errorf("expected written count, got: %s", out)
	}
}

func TestProgressConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_PROGRESS", "")
	t.Setenv("VAULTPULL_QUIET", "")
	cfg := ProgressConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected progress enabled by default")
	}
	if cfg.Quiet {
		t.Error("expected quiet false by default")
	}
}

func TestProgressConfigFromEnv_QuietMode(t *testing.T) {
	t.Setenv("VAULTPULL_QUIET", "true")
	cfg := ProgressConfigFromEnv()
	if !cfg.Quiet {
		t.Error("expected quiet to be true")
	}
}

func TestProgressConfigFromEnv_DisabledProgress(t *testing.T) {
	t.Setenv("VAULTPULL_PROGRESS", "0")
	cfg := ProgressConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected progress to be disabled")
	}
}
