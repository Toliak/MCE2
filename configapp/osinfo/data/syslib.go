package data

import (
	"encoding/json"
	"fmt"
)

// Package Manager enum with Rust-like unknown value
type SysLib struct {
	V     SysLibE				`json:"value"`

	// only used for Unknown variant
	Raw   string				`json:"raw"`
}

type SysLibE int

const (
	SysLibGlibc SysLibE = iota
	SysLibMusl
	SysLibUnknown
)

func (s SysLibE) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func NewSysLib(mgr SysLibE, raw string) SysLib {
	return SysLib{
		V: mgr,
		Raw: raw,
	}
}

func (os *SysLib) String() string {
	return os.V.String()
}

func (os *SysLib) GoString() string {
	return fmt.Sprintf("%s(%s)", os.V.String(), os.Raw)
}

func (os SysLibE) String() string {
	switch os {
	case SysLibGlibc:
		return "glibc"
	case SysLibMusl:
		return "musl"
	case SysLibUnknown:
		return "unknown"
	default:
		return "invalid"
	}
}
