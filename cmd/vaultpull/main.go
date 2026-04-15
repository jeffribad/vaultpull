package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/dotenv"
	"vaultpull/internal/vault"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"

	// flags
	cfgFile    string
	outputFile string
	role       string
	overwrite  bool
	dryRun     bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull reads secrets from a HashiCorp Vault KV store and writes
them into a local .env file, optionally filtering keys by role-based policies.`,
	Version: Version,
	RunE:    run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: .vaultpull.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", ".env", "path to the output .env file")
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "role to filter secrets by (e.g. backend, frontend, ops)")
	rootCmd.PersistentFlags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in the .env file")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print secrets to stdout without writing to disk")
}

func run(cmd *cobra.Command, args []string) error {
	// Load configuration from env vars or config file.
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Override output path if provided via flag.
	if cmd.Flags().Changed("output") {
		cfg.OutputFile = outputFile
	} else if cfg.OutputFile == "" {
		cfg.OutputFile = outputFile
	}

	// Override role if provided via flag.
	if cmd.Flags().Changed("role") {
		cfg.Role = role
	}

	// Build the Vault client.
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	// Fetch secrets from Vault.
	secrets, err := client.ReadSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("reading secrets from vault path %q: %w", cfg.SecretPath, err)
	}

	// Apply role-based filtering when a role is specified.
	if cfg.Role != "" {
		policies, pErr := vault.DefaultPoliciesFromEnv()
		if pErr != nil {
			log.Printf("warning: could not load role policies: %v — skipping filter", pErr)
		} else {
			secrets, err = vault.FilterByRole(secrets, cfg.Role, policies)
			if err != nil {
				return fmt.Errorf("filtering secrets by role %q: %w", cfg.Role, err)
			}
		}
	}

	if len(secrets) == 0 {
		log.Println("no secrets matched — nothing to write")
		return nil
	}

	// Dry-run: print to stdout and exit.
	if dryRun k, v := range secrets {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	}

	// Write secrets to the .env file.
	writer, err := dotenv.NewWriter(cfg.OutputFile, overwrite)
	if err != nil {
		return fmt.Errorf("creating env writer: %w", err)	written, skipped, err := writer.Write(secrets)
	if err != nil {
		return fmt.Errorf("writing .env file: %w", err)
	}

	log.Printf("done: %d key(s) written, %d key(s) skipped → %s", written, skipped, cfg.OutputFile)
	return nil
}
