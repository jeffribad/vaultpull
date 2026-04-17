package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/vaultpull/internal/dotenv"
)

// WriteGroups writes each group of secrets to a separate .env file inside outDir.
// Files are named <group>.env.
func WriteGroups(groups map[string]map[string]string, outDir string, overwrite bool) error {
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	for groupName, secrets := range groups {
		filePath := filepath.Join(outDir, groupName+".env")
		w, err := dotenv.NewWriter(filePath, overwrite)
		if err != nil {
			return fmt.Errorf("creating writer for group %s: %w", groupName, err)
		}
		if err := w.Write(secrets); err != nil {
			return fmt.Errorf("writing group %s: %w", groupName, err)
		}
	}
	return nil
}

// GroupFileNames returns the expected file paths for each group.
func GroupFileNames(groups map[string]map[string]string, outDir string) []string {
	paths := make([]string, 0, len(groups))
	for name := range groups {
		paths = append(paths, filepath.Join(outDir, name+".env"))
	}
	return paths
}
