package sync

import (
	"bytes"
	"os"
	"testing"
)

func TestHookRunner_RunPre_ExecutesCommand(t *testing.T) {
	hooks := []Hook{
		{Type: HookPreSync, Command: "echo pre-hook"},
	}
	r, w, _ := os.Pipe()
	runner := NewHookRunner(hooks, w)

	if err := runner.RunPre(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() == 0 {
		t.Error("expected output from pre hook")
	}
}

func TestHookRunner_RunPost_ExecutesCommand(t *testing.T) {
	hooks := []Hook{
		{Type: HookPostSync, Command: "echo post-hook"},
	}
	r, w, _ := os.Pipe()
	runner := NewHookRunner(hooks, w)

	if err := runner.RunPost(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() == 0 {
		t.Error("expected output from post hook")
	}
}

func TestHookRunner_SkipsWrongType(t *testing.T) {
	hooks := []Hook{
		{Type: HookPostSync, Command: "false"},
	}
	runner := NewHookRunner(hooks, os.Stdout)
	if err := runner.RunPre(); err != nil {
		t.Fatalf("pre-sync should not run post-sync hooks, got %v", err)
	}
}

func TestHookRunner_FailingCommand_ReturnsError(t *testing.T) {
	hooks := []Hook{
		{Type: HookPreSync, Command: "false"},
	}
	runner := NewHookRunner(hooks, os.Stdout)
	if err := runner.RunPre(); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestHookRunner_EmptyCommand_ReturnsError(t *testing.T) {
	hooks := []Hook{
		{Type: HookPreSync, Command: ""},
	}
	runner := NewHookRunner(hooks, os.Stdout)
	if err := runner.RunPre(); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestNewHookRunner_NilOut_DefaultsToStdout(t *testing.T) {
	runner := NewHookRunner(nil, nil)
	if runner.out == nil {
		t.Error("expected non-nil output writer")
	}
}

func TestHookRunner_MultipleHooks_AllExecuted(t *testing.T) {
	r, w, _ := os.Pipe()
	hooks := []Hook{
		{Type: HookPreSync, Command: "echo first"},
		{Type: HookPreSync, Command: "echo second"},
	}
	runner := NewHookRunner(hooks, w)

	if err := runner.RunPre(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("first")) {
		t.Error("expected output from first hook")
	}
	if !bytes.Contains([]byte(output), []byte("second")) {
		t.Error("expected output from second hook")
	}
}
