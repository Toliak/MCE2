package osinfo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/toliak/mce/osinfo/data"
)

// Returns false, if nothing found. True -- if found.
func detectPkgManagerByBinary() (data.PkgManager, bool) {
	// Fallback: check for binaries
	switch {
	case commandExists("apt-get") || commandExists("apt"):
		return data.NewPkgManager(data.PkgMgrAptGet, "apt-get"), true
	case commandExists("apk"):
		return data.NewPkgManager(data.PkgMgrApk, "apk"), true
	case commandExists("dnf"):
		return data.NewPkgManager(data.PkgMgrDnf, "dnf"), true
	case commandExists("yum"):
		return data.NewPkgManager(data.PkgMgrYum, "yum"), true
	case commandExists("pacman"):
		return data.NewPkgManager(data.PkgMgrPacman, "pacman"), true
	default:
		return data.PkgManagerUnknown(), false
	}
}

// detectPkgManager infers package manager from distro ID and filesystem checks.
// Returns false, if nothing found. True -- if found.
func detectPkgManagerByOsVersion() (data.PkgManager, bool) {
	// First try to read /etc/os-release again for ID/ID_LIKE
	var id, idLike string
	if f, err := os.Open("/etc/os-release"); err == nil {
		props := parseKeyValueFile(f)
		f.Close()
		id = strings.ToLower(props["ID"])
		idLike = strings.ToLower(props["ID_LIKE"])
	} else {
		return data.PkgManagerUnknown(), false
	}

	// Distribution-based detection
	switch {
	case id == "alpine" || idLike == "alpine" ||fileExists("/etc/alpine-release"):
		if !commandExists("apk") {
			break;
		}
		return data.NewPkgManager(data.PkgMgrApk, "apk"), true
	case id == "arch" || id == "manjaro" || idLike == "arch" || fileExists("/etc/pacman.conf"):
		if !commandExists("pacman") {
			break;
		}
		return data.NewPkgManager(data.PkgMgrPacman, "pacman"), true
	case id == "ubuntu" || id == "debian" || idLike == "debian" || fileExists("/etc/apt/sources.list"):
		if commandExists("apt") {
			return data.NewPkgManager(data.PkgMgrAptGet, "apt"), true
		}
		if commandExists("apt-get") {
			return data.NewPkgManager(data.PkgMgrAptGet, "apt-get"), true
		}
	case id == "fedora" || idLike == "fedora" || fileExists("/etc/dnf/dnf.conf") || id == "rhel" || id == "centos" || id == "ol" || fileExists("/etc/yum.conf"):
		if commandExists("dnf") {
			return data.NewPkgManager(data.PkgMgrDnf, "dnf"), true
		}
		if commandExists("yum") {
			return data.NewPkgManager(data.PkgMgrYum, "yum"), true
		}
	}

	return data.PkgManagerUnknown(), false
}

func harvestPkgManager() data.PkgManager {
	if pkg, found := detectPkgManagerByOsVersion(); found == true {
		return pkg
	}

	pkg, _ := detectPkgManagerByBinary()
	return pkg;
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
		props := parseKeyValueFile(f)
		f.Close()

		id := strings.ToLower(props["ID"])
		idLike := strings.ToLower(props["ID_LIKE"])
		name := props["NAME"]
		versionStr := strings.ToLower(props["VERSION_ID"])
		version := data.ParseVersion(versionStr)

		return data.NewDistrib(id, idLike, name, version)
	} else {
		return data.NewUnknownDistrib()
	}
}
