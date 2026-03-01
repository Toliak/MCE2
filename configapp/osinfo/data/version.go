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
	v VersionT
	raw string
}

type VersionT struct {
	major int
	minor int
	patch int
}

func NewVersion(major int, minor int, patch int, raw string) Version {
	return Version{
		v: VersionT{
			major: major,
			minor: minor,
			patch: patch,
		},
		raw: raw,
	}
}

// ParseVersion parses a version string like "6.5.4" into a Version.
// The parsing is lenient: missing parts default to 0, invalid parts become 0,
// and extra parts beyond three are ignored. 
// The original raw string is preserved.
func ParseVersion(raw string) Version {
	parts := strings.Split(raw, ".")
	v := VersionT{
		major: 0,
		minor: 0,
		patch: 0,
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
			v.major = val
		case 1:
			v.minor = val
		case 2:
			v.patch = val
		}
	}

	return Version{
		v: v,
		raw: raw,
	}
}

func (v *Version) V() VersionT {
	return v.v
}

func (v *Version) Raw() string {
	return v.raw
}

func (v *Version) String() string {
	return v.v.String()
}

func (v *Version) GoString() string {
	return fmt.Sprintf("%s (%s)", v.v.String(), v.raw)
}

func (v *VersionT) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

// Major returns the major version
func (v *VersionT) Major() int {
    return v.major
}

// Minor returns the minor version
func (v *VersionT) Minor() int {
    return v.minor
}

// Patch returns the patch version
func (v *VersionT) Patch() int {
    return v.patch
}