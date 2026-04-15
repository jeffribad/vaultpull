package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' and blank lines are ignored.
// Values may optionally be wrapped in double quotes.
func Parse(filePath string) (map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) && len(val) >= 2 {
			val = val[1 : len(val)-1]
			val = strings.ReplaceAll(val, `\"`, `"`)
		}

		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return result, nil
}
