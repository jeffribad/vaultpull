package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatExport ExportFormat = "export"
	FormatJSON   ExportFormat = "json"
)

// ExportSecrets writes secrets to w in the given format.
func ExportSecrets(w io.Writer, secrets map[string]string, format ExportFormat) error {
	if w == nil {
		w = os.Stdout
	}
	keys := sortedSecretKeys(secrets)
	switch format {
	case FormatExport:
		for _, k := range keys {
			fmt.Fprintf(w, "export %s=%q\n", k, secrets[k])
		}
	case FormatJSON:
		fmt.Fprintln(w, "{")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(w, "  %q: %q%s\n", k, secrets[k], comma)
		}
		fmt.Fprintln(w, "}")
	case FormatDotenv:
		fallthrough
	default:
		for _, k := range keys {
			v := secrets[k]
			if strings.ContainsAny(v, " \t\n") {
				fmt.Fprintf(w, "%s=%q\n", k, v)
			} else {
				fmt.Fprintf(w, "%s=%s\n", k, v)
			}
		}
	}
	return nil
}

func sortedSecretKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
