package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadDotEnv(files ...string) error {
	if len(files) == 0 {
		files = []string{".env"}
	}

	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("open %s: %w", path, err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if idx := strings.Index(line, "#"); idx >= 0 {
				line = strings.TrimSpace(line[:idx])
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"`)
			_ = os.Setenv(key, val)
		}
		if err := scanner.Err(); err != nil {
			_ = f.Close()
			return fmt.Errorf("scan %s: %w", path, err)
		}
		_ = f.Close()
	}
	return nil
}
