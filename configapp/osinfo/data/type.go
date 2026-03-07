package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OS Type enum with Rust-like unknown value
type OSType struct {
	V     OSTypeE			`json:"value"`

	// only used for Unknown variant
	Raw   string			`json:"raw"`
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

func (s OSTypeE) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func ParseOsType(s string) OSType {
	return OSType{
		V: parseOsTypeE(s),
		Raw: s,
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

func (os *OSType) String() string {
	return os.V.String()
}

func (os *OSType) GoString() string {
	return fmt.Sprintf("%s(%s)", os.V.String(), os.Raw)
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
