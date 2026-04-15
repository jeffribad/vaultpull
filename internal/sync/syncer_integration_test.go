package sync_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/dotenv"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestSyncer_Integration_FullPipeline(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":    "db.internal",
		"DB_PASS":    "s3cr3t",
		"STRIPE_KEY": "sk_live_xxx",
	}
	policies := vault.Policies{
		"backend":  {"DB_HOST", "DB_PASS"},
		"payments": {"STRIPE_KEY"},
	}

	client, err := vault.NewClient("http://127.0.0.1:8200", "tok")
	if err != nil {
		t.Fatal(err)
	}
	client.SetFakeSecrets(secrets)

	var logBuf bytes.Buffer
	syncer := sync.New(client, audit.NewLogger(&logBuf))

	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	result, err := syncer.Run("secret/data/app", policies, sync.Options{
		Role:       "backend",
		OutputPath: outPath,
		Overwrite:  true,
	})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result.Written != 2 {
		t.Errorf("expected 2 written, got %d", result.Written)
	}

	parsed, err := dotenv.Parse(outPath)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if parsed["DB_HOST"] != "db.internal" {
		t.Errorf("DB_HOST mismatch: %q", parsed["DB_HOST"])
	}
	if _, ok := parsed["STRIPE_KEY"]; ok {
		t.Error("STRIPE_KEY should not be present for backend role")
	}

	logData, _ := os.ReadFile(outPath)
	if strings.Contains(string(logData), "STRIPE_KEY") {
		t.Error(".env file must not contain STRIPE_KEY")
	}
}
