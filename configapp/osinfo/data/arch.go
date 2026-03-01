package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CPU Architecture enum with Rust-like unknown value
type CPUArch struct {
	V   CPUArchE	`json:"value"`
	Raw string		`json:"raw"`
}

// Enum variant of the CPUArchE
type CPUArchE int

const (
	// 
	CPUArchAMD64 CPUArchE = iota

	// Experimental, untested
	CPUArchI386

	// 
	CPUArchAARCH64

	// Experimental, untested
	CPUArchARMv7

	// Experimental, untested
	CPUArchMIPS64Le

	// Experimental, untested
	CPUArchPPC64

	// Experimental, untested
	CPUArchRISCV64

	// Fallback value
	CPUArchUnknown
)

// MarshalJSON implements json.Marshaler
func (s CPUArchE) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func ParseCPUArch(s string) CPUArch {
	return CPUArch{
		V: parseCPUArchE(s),
		Raw: s,
	}
}

func parseCPUArchE(s string) CPUArchE {
	switch strings.ToLower(s) {
	case "amd64", "x86_64", "x64":
		return CPUArchAMD64
	case "i386", "x86", "386":
		return CPUArchI386
	case "aarch64", "arm64":
		return CPUArchAARCH64
	case "arm", "armv7", "armhf", "armv7l":
		return CPUArchARMv7
	case "mips64le":
		return CPUArchMIPS64Le
	case "ppc64":
		return CPUArchPPC64
	case "riscv64":
		return CPUArchRISCV64
	default:
		return CPUArchUnknown
	}
}

func (c *CPUArch) String() string {
	return c.V.String()
}

func (c *CPUArch) GoString() string {
	return fmt.Sprintf("%s(%s)", c.V.String(), c.Raw)
}

func (c CPUArchE) String() string {
	switch c {
	case CPUArchAMD64:
		return "amd64"
	case CPUArchI386:
		return "i386"
	case CPUArchAARCH64:
		return "aarch64"
	case CPUArchARMv7:
		return "armv7"
	case CPUArchMIPS64Le:
		return "mips64le"
	case CPUArchPPC64:
		return "ppc64"
	case CPUArchRISCV64:
		return "riscv64"
	case CPUArchUnknown:
		return "unknown"
	default:
		return "invalid"
	}
}
