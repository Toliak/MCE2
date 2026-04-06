package platform

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// parseKeyValueFile reads shell-style KEY=VALUE files (supports quoted values)
func ParseKeyValueFile(r io.Reader) map[string]string {
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

func FileEntryExists(path string) bool {
	_, err := os.Stat(path)
	return !(err != nil && os.IsNotExist(err))
}

// See https://opensource.com/article/18/6/copying-files-go
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
			return err
	}

	if !sourceFileStat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
			return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
			return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}