package data

import (
	"encoding/json"
	"fmt"
)

// Package Manager enum with Rust-like unknown value
type PkgManager struct {
	V     PkgManagerE			`json:"value"`

	// only used for Unknown and AptGet variant. Describes Raw command name
	Raw   string				`json:"raw"`
}

type PkgManagerE int

const (
	PkgMgrAptGet PkgManagerE = iota
	PkgMgrApk
	PkgMgrDnf
	PkgMgrMicroDnf
	PkgMgrYum
	PkgMgrPacman
	PkgMgrBrew
	PkgMgrWinget
	PkgMgrScoop
	PkgMgrUnknown
)

func (s PkgManagerE) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func NewPkgManager(mgr PkgManagerE, raw string) PkgManager {
	return PkgManager{
		V: mgr,
		Raw: raw,
	}
}

func PkgManagerUnknown() PkgManager {
	return NewPkgManager(PkgMgrUnknown, "unknown")
}

func (os *PkgManager) String() string {
	return os.V.String()
}

func (os *PkgManager) GoString() string {
	return fmt.Sprintf("%s(%s)", os.V.String(), os.Raw)
}

func (os PkgManagerE) String() string {
	switch os {
	case PkgMgrAptGet:
		return "apt-get"
	case PkgMgrApk:
		return "apk"
	case PkgMgrDnf:
		return "dnf"
	case PkgMgrMicroDnf:
		return "microdnf"
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
