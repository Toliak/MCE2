package harvest

import (
	"os"
	"path/filepath"
	"slices" // Import the standard slices package
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
)

// Returns false, if nothing found. True -- if found.
func detectPkgManagerByBinary() (data.PkgManager, bool) {
	// Fallback: check for binaries
	switch {
	case platform.CommandExists("apt-get") && platform.CommandExists("apt-cache"):
		// Do not trust without `apt-cache` because we have to search packages in some tasks
		return data.NewPkgManager(data.PkgMgrAptGet, "apt-get"), true
	case platform.CommandExists("apk"):
		return data.NewPkgManager(data.PkgMgrApk, "apk"), true
	case platform.CommandExists("dnf"):
		return data.NewPkgManager(data.PkgMgrDnf, "dnf"), true
	case platform.CommandExists("microdnf"):
		return data.NewPkgManager(data.PkgMgrMicroDnf, "microdnf"), true
	case platform.CommandExists("yum"):
		return data.NewPkgManager(data.PkgMgrYum, "yum"), true
	case platform.CommandExists("pacman"):
		return data.NewPkgManager(data.PkgMgrPacman, "pacman"), true
	default:
		return data.PkgManagerUnknown(), false
	}
}

// detectPkgManager infers package manager from distro ID and filesystem checks.
// Returns false, if nothing found. True -- if found.
func detectPkgManagerByOsVersion(d *data.Distrib) (data.PkgManager, bool) {
	// Distribution-based detection
	switch {
	case d.Id == "alpine" || slices.Contains(d.IdLike, "alpine") || platform.FileExists("/etc/alpine-release"):
		if !platform.CommandExists("apk") {
			break
		}
		return data.NewPkgManager(data.PkgMgrApk, "apk"), true
	case d.Id == "arch" || d.Id == "manjaro" || slices.Contains(d.IdLike, "arch") || platform.FileExists("/etc/pacman.conf"):
		if !platform.CommandExists("pacman") {
			break
		}
		return data.NewPkgManager(data.PkgMgrPacman, "pacman"), true
	case d.Id == "ubuntu" || d.Id == "debian" || slices.Contains(d.IdLike, "debian") || platform.FileExists("/etc/apt/sources.list"):
		if !platform.CommandExists("apt-cache") {
			break;
		}

		if platform.CommandExists("apt-get") {
			return data.NewPkgManager(data.PkgMgrAptGet, "apt-get"), true
		}
		if platform.CommandExists("apt") {
			return data.NewPkgManager(data.PkgMgrAptGet, "apt"), true
		}
	case d.Id == "fedora" ||
		slices.Contains(d.IdLike, "fedora") ||
		platform.FileExists("/etc/dnf/dnf.conf") ||
		slices.Contains(d.IdLike, "rhel") ||
		slices.Contains(d.IdLike, "centos") ||
		platform.FileExists("/etc/yum.conf"):

		if platform.CommandExists("dnf") {
			return data.NewPkgManager(data.PkgMgrDnf, "dnf"), true
		}
		if platform.CommandExists("microdnf") {
			return data.NewPkgManager(data.PkgMgrMicroDnf, "microdnf"), true
		}
		if platform.CommandExists("yum") {
			return data.NewPkgManager(data.PkgMgrYum, "yum"), true
		}
	}

	return data.PkgManagerUnknown(), false
}

func harvestPkgManager(d *data.Distrib) data.PkgManager {
	if pkg, found := detectPkgManagerByOsVersion(d); found == true {
		return pkg
	}

	pkg, _ := detectPkgManagerByBinary()
	return pkg
}

// ------------

// Returns true if found something. False otherwise
func checkSysLibPattern(dirs []string, patterns []string) bool {
	// Check musl first (combine dirs and patterns)
	for _, dir := range dirs {
		for _, pattern := range patterns {
			matches, err := filepath.Glob(filepath.Join(dir, pattern))
			if err == nil && len(matches) > 0 {
				return true
			}
		}
	}

	return false
}

func detectSysLibBySo() (data.SysLib, bool) {
	// Common library directories
	dirs := []string{
		"/lib",
		"/usr/lib",
		"/usr/local/lib",
		"/lib64",
		"/usr/lib64",
		"/lib32",
		"/usr/lib32",
	}

	// Musl patterns
	muslPatterns := []string{
		"ld-musl*.so*",
		"libc.musl*.so*",
	}

	// Glibc patterns
	glibcPatterns := []string{
		"ld-linux*.so*",
		"ld-*.so",
		"libc.so*",
		"libc-*.so",
	}

	if checkSysLibPattern(dirs, glibcPatterns) {
		// TODO: do we need the library path here?
		return data.NewSysLib(data.SysLibGlibc, "glibc"), true
	}
	if checkSysLibPattern(dirs, muslPatterns) {
		return data.NewSysLib(data.SysLibMusl, "musl"), true
	}

	return data.NewSysLib(data.SysLibUnknown, ""), false
}

func harvestSysLib() data.SysLib {
	res, _ := detectSysLibBySo()
	return res
}

// ------------

func harvestDistrib() data.Distrib {
	if f, err := os.Open("/etc/os-release"); err == nil {
		props := platform.ParseKeyValueFile(f)
		f.Close()

		id := strings.ToLower(props["ID"])
		idLike := strings.Split(strings.ToLower(props["ID_LIKE"]), " ")
		name := props["NAME"]
		versionStr := strings.ToLower(props["VERSION_ID"])
		version := data.ParseVersion(versionStr)

		return data.NewDistrib(id, idLike, name, version)
	} else {
		return data.NewUnknownDistrib()
	}
}
