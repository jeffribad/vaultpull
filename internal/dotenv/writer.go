package dotenv

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
	overwrite bool
}

// NewWriter creates a new Writer targeting the given file path.
func NewWriter(filePath string, overwrite bool) *Writer {
	return &Writer{
		filePath: filePath,
		overwrite: overwrite,
	}
}

// Write serializes the provided secrets map into a .env file.
// Existing keys are preserved unless overwrite is true.
func (w *Writer) Write(secrets map[string]string) error {
	existing := map[string]string{}

	if !w.overwrite {
		var err error
		existing, err = Parse(w.filePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reading existing env file: %w", err)
		}
	}

	merged := make(map[string]string, len(existing)+len(secrets))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range secrets {
		merged[k] = v
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, quoteValue(merged[k])))
	}

	if err := os.WriteFile(w.filePath, []byte(sb.String()), 0600); err != nil {
		return fmt.Errorf("writing env file %q: %w", w.filePath, err)
	}
	return nil
}

// quoteValue wraps the value in double quotes if it contains spaces or special chars.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
