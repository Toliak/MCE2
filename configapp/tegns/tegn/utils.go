package tegn

import (
	"fmt"
	"os"
	"regexp"

	tb "github.com/toliak/mce/tegnbuilder"
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

// GetGitCloneParameters generates the standard set of parameters for Git clone operations.
func GetGitCloneParameters(defaultRepoURL string, defaultInstallDir string) []tb.TegnParameter {
	return []tb.TegnParameter{
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue(defaultRepoURL),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch, leave empty to use the default branch"),
			tb.WithDefaultValue(""),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"use-shallow-clone",
			"Shallow Clone",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Perform a shallow clone with depth 1 (--depth 1) for faster download"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"install-dir",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation path (read-only)"),
			tb.WithDefaultValue(defaultInstallDir),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}
