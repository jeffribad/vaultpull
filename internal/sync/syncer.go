package sync

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/dotenv"
	"github.com/your-org/vaultpull/internal/vault"
)

// Options holds configuration for a sync operation.
type Options struct {
	Role       string
	OutputPath string
	Overwrite  bool
	DryRun     bool
}

// Result summarises what happened during a sync.
type Result struct {
	Written  int
	Skipped  int
	FilePath string
}

// Syncer orchestrates fetching secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	logger *audit.Logger
}

// New creates a Syncer with the given Vault client and audit logger.
func New(client *vault.Client, logger *audit.Logger) *Syncer {
	return &Syncer{client: client, logger: logger}
}

// Run executes the sync pipeline: read → filter → write.
func (s *Syncer) Run(secretPath string, policies vault.Policies, opts Options) (*Result, error) {
	secrets, err := s.client.ReadSecrets(secretPath)
	if err != nil {
		s.logger.LogError(fmt.Sprintf("vault read failed: %v", err))
		return nil, fmt.Errorf("reading secrets: %w", err)
	}

	filtered := vault.FilterByRole(secrets, opts.Role, policies)
	if len(filtered) == 0 {
		s.logger.LogError(fmt.Sprintf("no secrets matched role %q", opts.Role))
		return nil, fmt.Errorf("no secrets available for role %q", opts.Role)
	}

	if opts.DryRun {
		s.logger.LogSync(opts.Role, secretPath, len(filtered), true)
		return &Result{Written: 0, Skipped: len(filtered), FilePath: opts.OutputPath}, nil
	}

	w, err := dotenv.NewWriter(opts.OutputPath, opts.Overwrite)
	if err != nil {
		return nil, fmt.Errorf("creating writer: %w", err)
	}

	n, err := w.Write(filtered)
	if err != nil {
		s.logger.LogError(fmt.Sprintf("write failed: %v", err))
		return nil, fmt.Errorf("writing .env file: %w", err)
	}

	s.logger.LogSync(opts.Role, secretPath, n, false)
	return &Result{Written: n, FilePath: opts.OutputPath}, nil
}
