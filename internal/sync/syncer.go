package sync

import (
	"context"
	"fmt"

	"github.com/user/vaultpull/internal/audit"
	"github.com/user/vaultpull/internal/dotenv"
	"github.com/user/vaultpull/internal/vault"
)

// VaultClient is the interface required by Syncer to read secrets.
type VaultClient interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	client   VaultClient
	policies map[string][]string
	logger   *audit.Logger
}

// New creates a new Syncer with the provided client, role policies, and audit logger.
func New(client VaultClient, policies map[string][]string, logger *audit.Logger) *Syncer {
	return &Syncer{
		client:   client,
		policies: policies,
		logger:   logger,
	}
}

// SyncOptions controls the behaviour of a sync operation.
type SyncOptions struct {
	SecretPath string
	OutputPath string
	Role       string
	DryRun     bool
	Overwrite  bool
	Backup     bool
}

// Run executes the sync: fetch secrets, filter by role, optionally backup, then write.
func (s *Syncer) Run(ctx context.Context, opts SyncOptions) (*Result, error) {
	if opts.Role != "" {
		if _, ok := s.policies[opts.Role]; !ok {
			return nil, fmt.Errorf("unknown role %q: no policy defined", opts.Role)
		}
	}

	secrets, err := s.client.ReadSecrets(ctx, opts.SecretPath)
	if err != nil {
		s.logger.LogError(opts.SecretPath, err)
		return nil, fmt.Errorf("read secrets: %w", err)
	}

	filtered := vault.FilterByRole(secrets, opts.Role, s.policies)

	if opts.DryRun {
		s.logger.LogSync(opts.SecretPath, opts.OutputPath, opts.Role, len(filtered), true)
		return &Result{Written: 0, DryRun: true, Keys: keys(filtered)}, nil
	}

	var bak *Backup
	if opts.Backup {
		bak, err = CreateBackup(opts.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("create backup: %w", err)
		}
	}

	w, err := dotenv.NewWriter(opts.OutputPath, opts.Overwrite)
	if err != nil {
		_ = bak.Restore()
		return nil, fmt.Errorf("open writer: %w", err)
	}

	if err := w.Write(filtered); err != nil {
		_ = bak.Restore()
		return nil, fmt.Errorf("write secrets: %w", err)
	}

	_ = bak.Discard()
	s.logger.LogSync(opts.SecretPath, opts.OutputPath, opts.Role, len(filtered), false)
	return &Result{Written: len(filtered), DryRun: false, Keys: keys(filtered)}, nil
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
