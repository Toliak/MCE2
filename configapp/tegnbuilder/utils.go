package tegnbuilder

import (
	"fmt"
	"regexp"
	"strings"
)


func NameToID(name string) (*string, error) {
	// Check if string contains only english letters, digits, spaces, underscores or dashes
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 _-]+$`, name)
	if err != nil || !matched {
		return nil, fmt.Errorf("Unable to convert name '%s' to ID", name)
	}

	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces and underscores with dashes
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Collapse multiple consecutive dashes into a single dash
	re := regexp.MustCompile(`-+`)
	name = re.ReplaceAllString(name, "-")

	return &name, nil
}

// For global usage only
func NameToIDUnsafe(name string) string {
	result, err := NameToID(name)
	if err != nil {
		panic(err)
	}

	return *result
}
