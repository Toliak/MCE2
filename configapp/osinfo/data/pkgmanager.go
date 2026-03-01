package data

import (
	"fmt"
)

// Package Manager enum with Rust-like unknown value
type PkgManager struct {
	v     PkgManagerE
	raw   string // only used for Unknown and AptGet variant. Describes raw command name
}

type PkgManagerE int

const (
	PkgMgrAptGet PkgManagerE = iota
	PkgMgrApk
	PkgMgrDnf
	PkgMgrYum
	PkgMgrPacman
	PkgMgrBrew
	PkgMgrWinget
	PkgMgrScoop
	PkgMgrUnknown
)

func NewPkgManager(mgr PkgManagerE, raw string) PkgManager {
	return PkgManager{
		v: mgr,
		raw: raw,
	}
}

func PkgManagerUnknown() PkgManager {
	return NewPkgManager(PkgMgrUnknown, "unknown")
}


func (os *PkgManager) V() PkgManagerE { 
	return os.v
}

func (os *PkgManager) Raw() string { 
	return os.raw
}

func (os *PkgManager) String() string {
	return os.v.String()
}

func (os *PkgManager) GoString() string {
	return fmt.Sprintf("%s(%s)", os.v.String(), os.raw)
}

func (os PkgManagerE) String() string {
	switch os {
	case PkgMgrAptGet:
		return "apt-get"
	case PkgMgrApk:
		return "apk"
	case PkgMgrDnf:
		return "dnf"
	case PkgMgrYum:
		return "yum"
	case PkgMgrPacman:
		return "pacman"
	case PkgMgrBrew:
		return "brew"
	case PkgMgrWinget:
		return "winget"
	case PkgMgrScoop:
		return "scoop"
	case PkgMgrUnknown:
		return "unknown"
	default:
		return "invalid"
	}
}
