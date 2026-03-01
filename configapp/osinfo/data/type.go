package data

import (
	"fmt"
	"strings"
)

// OS Type enum with Rust-like unknown value
type OSType struct {
	v     OSTypeE
	raw   string // only used for Unknown variant
}

type OSTypeE int

const (
	// 
	OSTypeDarwin OSTypeE = iota

	// 
	OSTypeLinux

	// 
	OSTypeWindows
	
	// Experimental, untested
	OSTypeAndroid

	// Fallback value
	OSTypeUnknown
)

func ParseOsType(s string) OSType {
	return OSType{
		v: parseOsTypeE(s),
		raw: s,
	}
}

func parseOsTypeE(s string) OSTypeE {
	switch strings.ToLower(s) {
	case "darwin":
		return OSTypeDarwin
	case "linux":
		return OSTypeLinux
	case "windows":
		return OSTypeWindows
	case "android":
		return OSTypeAndroid
	default:
		return OSTypeUnknown
	}
}

func (os *OSType) V() OSTypeE { 
	return os.v
}

func (os *OSType) Raw() string { 
	return os.raw
}

func (os *OSType) String() string {
	return os.v.String()
}

func (os *OSType) GoString() string {
	return fmt.Sprintf("%s(%s)", os.v.String(), os.raw)
}

func (os OSTypeE) String() string {
	switch os {
	case OSTypeDarwin:
		return "darwin"
	case OSTypeLinux:
		return "linux"
	case OSTypeWindows:
		return "windows"
	case OSTypeAndroid:
		return "android"
	case OSTypeUnknown:
		return "unknown"
	default:
		return "invalid"
	}
}
