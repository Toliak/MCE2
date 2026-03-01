package data

import (
	"fmt"
)

// Package Manager enum with Rust-like unknown value
type SysLib struct {
	v     SysLibE
	raw   string // only used for Unknown variant
}

type SysLibE int

const (
	SysLibGlibc SysLibE = iota
	SysLibMusl
	SysLibUnknown
)

func NewSysLib(mgr SysLibE, raw string) SysLib {
	return SysLib{
		v: mgr,
		raw: raw,
	}
}


func (os *SysLib) V() SysLibE { 
	return os.v
}

func (os *SysLib) Raw() string { 
	return os.raw
}

func (os *SysLib) String() string {
	return os.v.String()
}

func (os *SysLib) GoString() string {
	return fmt.Sprintf("%s(%s)", os.v.String(), os.raw)
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
