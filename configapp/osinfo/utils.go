package osinfo

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// parseKeyValueFile reads shell-style KEY=VALUE files (supports quoted values)
func parseKeyValueFile(r io.Reader) map[string]string {
	props := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove surrounding quotes
		val = strings.Trim(val, `"'`)

		props[key] = val
	}
	return props
}

// Checks if a file/glob pattern exists
func FileExists(pattern string) bool {
	if strings.Contains(pattern, "*") {
		matches, _ := filepath.Glob(pattern)
		return len(matches) > 0
	}
	_, err := os.Stat(pattern)
	return err == nil
}

// Checks if a command is available in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func IsRoot() bool {
    return os.Geteuid() == 0
}
