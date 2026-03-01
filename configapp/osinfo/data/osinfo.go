package data

import (
	"fmt"
)

type OsInfo struct {
	Arch CPUArch 				`json:"arch"`
	OsType OSType				`json:"osType"`
	PkgManager PkgManager		`json:"pkgManager"`
	SysLib SysLib				`json:"sysLib"`
	// kernelVersion Version
	Distrib Distrib				`json:"distrib"`
}

// ---------- Functional Options ----------

// OsInfoOption defines a functional option for configuring OsInfo.
type OsInfoOption func(*OsInfo)

// WithArch sets the CPU architecture.
func WithArch(a CPUArch) OsInfoOption {
	return func(o *OsInfo) { o.Arch = a }
}

// WithOSType sets the operating system type.
func WithOSType(t OSType) OsInfoOption {
	return func(o *OsInfo) { o.OsType = t }
}

// WithPkgManager sets the package manager.
func WithPkgManager(p PkgManager) OsInfoOption {
	return func(o *OsInfo) { o.PkgManager = p }
}

// WithSysLib sets the system library information.
func WithSysLib(l SysLib) OsInfoOption {
	return func(o *OsInfo) { o.SysLib = l }
}

// WithKernelVersion sets the kernel version.
// func WithKernelVersion(v Version) OsInfoOption {
// 	return func(o *OsInfo) { o.kernelVersion = v }
// }

// WithDistrib sets the distribution name.
func WithDistrib(d Distrib) OsInfoOption {
	return func(o *OsInfo) { o.Distrib = d }
}

// ---------- Constructor ----------

// NewOsInfo creates a new OsInfo applying any number of functional options.
func NewOsInfo(opts ...OsInfoOption) OsInfo {
	info := OsInfo{}
	for _, opt := range opts {
		opt(&info)
	}
	return info
}

// ---------- Stringer Interfaces ----------

// String returns a human‑readable representation of OsInfo.
func (o *OsInfo) String() string {
	return fmt.Sprintf(
		"%s (%s) – %s, pkg manager %s, sys lib %s",
		o.Distrib.String(),
		o.OsType.String(),
		o.Arch.String(),
		// o.kernelVersion.String(),
		o.PkgManager.String(),
		o.SysLib.String(),
	)
}

// GoString returns a Go‑syntax representation of OsInfo.
func (o *OsInfo) GoString() string {
	return fmt.Sprintf(
		"&OsInfo{arch: %#v, pkgManager: %#v, sysLib: %#v, osType: %#v, distrib: %#v}",
		o.Arch,
		o.PkgManager,
		o.SysLib,
		o.OsType,
		// o.kernelVersion,
		o.Distrib,
	)
}
