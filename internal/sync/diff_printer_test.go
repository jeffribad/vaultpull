package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	d := Diff(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	PrintDiff(&buf, d)

	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestPrintDiff_ShowsAddedKeys(t *testing.T) {
	var buf bytes.Buffer
	d := Diff(map[string]string{}, map[string]string{"NEW_KEY": "secret"})
	PrintDiff(&buf, d)

	if !strings.Contains(buf.String(), "+ NEW_KEY") {
		t.Errorf("expected '+ NEW_KEY' in output, got: %s", buf.String())
	}
}

func TestPrintDiff_ShowsUpdatedKeys(t *testing.T) {
	var buf bytes.Buffer
	d := Diff(map[string]string{"FOO": "old"}, map[string]string{"FOO": "new"})
	PrintDiff(&buf, d)

	if !strings.Contains(buf.String(), "~ FOO") {
		t.Errorf("expected '~ FOO' in output, got: %s", buf.String())
	}
}

func TestPrintDiff_ShowsRemovedKeys(t *testing.T) {
	var buf bytes.Buffer
	d := Diff(map[string]string{"GONE": "bye"}, map[string]string{})
	PrintDiff(&buf, d)

	if !strings.Contains(buf.String(), "- GONE") {
		t.Errorf("expected '- GONE' in output, got: %s", buf.String())
	}
}

func TestPrintDiff_DoesNotLeakValues(t *testing.T) {
	var buf bytes.Buffer
	d := Diff(map[string]string{}, map[string]string{"SECRET_KEY": "super-secret-value"})
	PrintDiff(&buf, d)

	if strings.Contains(buf.String(), "super-secret-value") {
		t.Error("diff output must not contain secret values")
	}
}
