package tegn

import (
	"fmt"
	"os"
	"regexp"
)


func removeConfigBlockFromFile(filePath string, blockName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error '%s': %w", filePath, err)
	}

	// Pattern to match the entire block including the BEGIN and END lines
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)# <BEGIN> %s.*?# <END> %s\n*`, regexp.QuoteMeta(blockName), regexp.QuoteMeta(blockName)))
	newContent := pattern.ReplaceAll(content, []byte{})

	// Remove any extra blank lines that might have been left
	// WARN: do not do that, it can break the formatting
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
