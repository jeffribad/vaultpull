package sync_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/sync"
)

func TestResult_Summary_Written(t *testing.T) {
	var buf bytes.Buffer
	r := &sync.Result{Written: 3, FilePath: ".env"}
	r.Summary(&buf)
	if !strings.Contains(buf.String(), "3 secret(s)") {
		t.Errorf("unexpected summary: %q", buf.String())
	}
}

func TestResult_Summary_DryRun(t *testing.T) {
	var buf bytes.Buffer
	r := &sync.Result{Written: 0, Skipped: 5, FilePath: ".env"}
	r.Summary(&buf)
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in summary, got: %q", out)
	}
	if !strings.Contains(out, "5 secret(s)") {
		t.Errorf("expected count in summary, got: %q", out)
	}
}

func TestResult_IsEmpty_True(t *testing.T) {
	r := &sync.Result{}
	if !r.IsEmpty() {
		t.Error("expected IsEmpty() == true")
	}
}

func TestResult_IsEmpty_False(t *testing.T) {
	r := &sync.Result{Written: 1}
	if r.IsEmpty() {
		t.Error("expected IsEmpty() == false")
	}
}
