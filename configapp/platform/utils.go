package platform

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
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

func OpenOrCreateFileAppend(path string) (*os.File, error) {
	if FileEntryExists(path) {
		v, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY, 0644)
		return v, err
	}

	v, err := os.Create(path)
	return v, err
}


func AppendFilepathString(path string, text string) error {
	outputFile, err := OpenOrCreateFileAppend(path)
    if err != nil {
        return err
    }
    defer outputFile.Close()

	_, err = outputFile.WriteString(text)
	if err != nil {
        return err
    }

	return nil
}

type ProgressWriter struct {
	Total      int64
	Downloaded int64
	Url        string
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)

	if pw.Total > 0 {
		percent := float64(pw.Downloaded) / float64(pw.Total) * 100
		fmt.Printf("\rDownloading... %.2f%%", percent)
	} else {
		fmt.Printf("\rDownloading... %d bytes", pw.Downloaded)
	}

	return n, nil
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
// See: https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9
func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Downloading '%s' -> '%s'\n", url, filepath)
	pw := &ProgressWriter{
		Total: resp.ContentLength,
		Url: url,
	}

	// Write the body to file
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	if err != nil {
		return err
	}
	fmt.Println("")

	return nil
}