package data

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: make the fields private and create the methods to access them (getters only)

// Version both the parsed version and the original raw string.
// The parsed version represents a tuple of three integers.
type Version struct {
	V     VersionT			`json:"value"`
	Raw   string			`json:"raw"`
}

type VersionT struct {
	Major int			`json:"major"`
	Minor int			`json:"minor"`
	Patch int			`json:"patch"`
}

func NewVersion(major int, minor int, patch int, raw string) Version {
	return Version{
		V: VersionT{
			Major: major,
			Minor: minor,
			Patch: patch,
		},
		Raw: raw,
	}
}

// ParseVersion parses a version string like "6.5.4" into a Version.
// The parsing is lenient: missing parts default to 0, invalid parts become 0,
// and extra parts beyond three are ignored. 
// The original raw string is preserved.
func ParseVersion(raw string) Version {
	parts := strings.Split(raw, ".")
	v := VersionT{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}

	// Parse up to 3 parts, ignore extra parts if more than 3
	for i := 0; i < 3 && i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		val, _ := strconv.Atoi(part)

		switch i {
		case 0:
			v.Major = val
		case 1:
			v.Minor = val
		case 2:
			v.Patch = val
		}
	}

	return Version{
		V: v,
		Raw: raw,
	}
}

func (v *Version) String() string {
	return v.V.String()
}

func (v *Version) GoString() string {
	return fmt.Sprintf("%s (%s)", v.V.String(), v.Raw)
}

func (v *VersionT) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
