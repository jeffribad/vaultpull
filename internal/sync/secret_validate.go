package sync

import (
	"fmt"
	"os"
	"strings"
)

// ValidationRule defines a rule type for secret validation.
type ValidationRule string

const (
	RuleRequired ValidationRule = "required"
	RuleNonempty  ValidationRule = "nonempty"
)

// ValidateConfig holds validation settings.
type ValidateConfig struct {
	Enabled      bool
	RequiredKeys []string
	NonemptyKeys []string
}

// ValidateConfigFromEnv reads validation config from environment variables.
// VAULTPULL_VALIDATE=true
// VAULTPULL_REQUIRED_KEYS=DB_HOST,DB_PORT
// VAULTPULL_NONEMPTY_KEYS=API_KEY,SECRET_TOKEN
func ValidateConfigFromEnv() ValidateConfig {
	enabled := strings.ToLower(os.Getenv("VAULTPULL_VALIDATE")) == "true" ||
		os.Getenv("VAULTPULL_VALIDATE") == "1"
	return ValidateConfig{
		Enabled:      enabled,
		RequiredKeys: splitTrimmed(os.Getenv("VAULTPULL_REQUIRED_KEYS"), ","),
		NonemptyKeys: splitTrimmed(os.Getenv("VAULTPULL_NONEMPTY_KEYS"), ","),
	}
}

// ValidationError holds all validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("secret validation failed: %s", strings.Join(e.Errors, "; "))
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ValidateSecrets checks secrets against the provided config rules.
func ValidateSecrets(secrets map[string]string, cfg ValidateConfig) error {
	if !cfg.Enabled {
		return nil
	}
	var errs []string
	for _, key := range cfg.RequiredKeys {
		if _, ok := secrets[key]; !ok {
			errs = append(errs, fmt.Sprintf("required key %q is missing", key))
		}
	}
	for _, key := range cfg.NonemptyKeys {
		val, ok := secrets[key]
		if !ok {
			errs = append(errs, fmt.Sprintf("nonempty key %q is missing", key))
		} else if strings.TrimSpace(val) == "" {
			errs = append(errs, fmt.Sprintf("nonempty key %q has empty value", key))
		}
	}
	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}
