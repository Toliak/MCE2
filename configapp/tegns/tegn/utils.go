package tegn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	git "github.com/go-git/go-git/v6"
)

func defaultGitCloneOptions(additionalOpts ...func(*git.CloneOptions)) *git.CloneOptions {
	v := git.CloneOptions{
		Progress: os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		NoCheckout: true,
	}

	for _, opt := range additionalOpts {
		opt(&v)
	}

	return &v
}

func MkdirAllParent(path string) error {
	parentDir := filepath.Dir(path)
	return os.MkdirAll(parentDir, 0755)
}

func removeConfigBlockFromFile(filePath string, blockName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error '%s': %w", filePath, err)
	}

	// Pattern to match the entire block including the BEGIN and END lines
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)# <BEGIN> %s.*?# <END> %s\n*`, regexp.QuoteMeta(blockName), regexp.QuoteMeta(blockName)))
	newContent := pattern.ReplaceAll(content, []byte{})

	// Remove any extra blank lines that might have been left
	// doubleNewlinePattern := regexp.MustCompile(`\n{3,}`)
	// newContent = doubleNewlinePattern.ReplaceAll(newContent, []byte("\n\n"))

	// Trim trailing whitespace
	// newContent = regexp.MustCompile(`\s+\z`).ReplaceAll(newContent, []byte("\n"))
	err = os.WriteFile(filePath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile error '%s': %w", filePath, err)
	}

	return nil
}
