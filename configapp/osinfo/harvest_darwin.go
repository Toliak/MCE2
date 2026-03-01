
package osinfo

import (
	"os/exec"
	"strings"

	"github.com/toliak/mce/osinfo/data"
)

// Returns false, if nothing found. True -- if found.
func detectPkgManagerByBinary() (data.PkgManager, bool) {
	// Darwin package managers
	switch {
	case commandExists("brew"):
		return data.NewPkgManager(data.PkgMgrBrew, "brew"), true
	// For now we do support only the brew pkg manager
	// case commandExists("port"):
	// 	return data.NewPkgManager(data.PkgMgrMacPorts, "port"), true
	// case commandExists("pkgin"):
	// 	return data.NewPkgManager(data.PkgMgrPkgin, "pkgin"), true // For NetBSD/macOS
	default:
		return data.PkgManagerUnknown(), false
	}
}

// detectPkgManagerByOsVersion infers package manager from Darwin version and filesystem checks.
// Returns false, if nothing found. True -- if found.
func detectPkgManagerByOsVersion() (data.PkgManager, bool) {
	// Check for common Darwin package manager paths
	switch {
	case fileExists("/opt/homebrew/bin/brew") || fileExists("/usr/local/bin/brew") || fileExists("/home/linuxbrew/.linuxbrew/bin/brew"):
		if commandExists("brew") {
			return data.NewPkgManager(data.PkgMgrBrew, "brew"), true
		}
	// case fileExists("/opt/local/bin/port"):
	// 	if commandExists("port") {
	// 		return data.NewPkgManager(data.PkgMgrMacPorts, "port"), true
	// 	}
	// case fileExists("/usr/pkg/bin/pkgin"):
	// 	if commandExists("pkgin") {
	// 		return data.NewPkgManager(data.PkgMgrPkgin, "pkgin"), true
	// 	}
	}

	return data.PkgManagerUnknown(), false
}

func harvestPkgManager() data.PkgManager {
	if pkg, found := detectPkgManagerByOsVersion(); found == true {
		return pkg
	}

	pkg, _ := detectPkgManagerByBinary()
	return pkg
}

// -----------------------

func harvestSysLib() data.SysLib {
	return data.NewSysLib(data.SysLibUnknown, "i don't care")
}

// -----------------------

func getMacOSProductName() string {
	cmd := exec.Command("sw_vers", "-productName")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getMacOSVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func harvestDistrib() data.Distrib {
	name := getMacOSProductName()
	versionStr := getMacOSVersion()

	return data.NewDistrib(
		"",
		"",
		name,
		data.ParseVersion(versionStr),
	)
}
