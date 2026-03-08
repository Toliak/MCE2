package harvest

import (
	"os"
	"path/filepath"
	"strings"
	"strconv"

	"github.com/toliak/mce/osinfo/data"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func commandExistsInPaths(cmd string, paths []string) bool {
	for _, path := range paths {
		fullPath := filepath.Join(path, cmd)
		if _, err := os.Stat(fullPath); err == nil {
			return true
		}
	}
	return false
}

// detectPkgManagerByBinary checks for installed package managers
func detectPkgManagerByBinary() (data.PkgManager, bool) {
	// Check common installation paths
	programFiles := os.Getenv("ProgramFiles")
	userProfile := os.Getenv("USERPROFILE")

	// Common paths for package managers
	paths := []string{
		filepath.Join(programFiles, "WindowsApps"),           // winget
		filepath.Join(userProfile, "scoop", "shims"),         // Scoop
	}

	if commandExistsInPaths("winget.exe", paths) || commandExists("winget.exe") {
		return data.NewPkgManager(data.PkgMgrWinget, "winget"), true
	}
	
	// if commandExistsInPaths("choco.exe", paths) || commandExists("choco.exe") {
	// 	return data.NewPkgManager(PkgMgrChoco, "choco"), true
	// }
	
	if commandExistsInPaths("scoop.exe", paths) || commandExists("scoop.exe") {
		return data.NewPkgManager(data.PkgMgrScoop, "scoop"), true
	}
	
	return data.PkgManagerUnknown(), false
}


// harvestPkgManager detects installed package managers on Windows
func harvestPkgManager(d *data.Distrib) data.PkgManager {
	pkg, _ := detectPkgManagerByBinary()
	
	return pkg
}

// -----------------------

// detectWindowsVersion returns Windows 10/11 specific version info
func detectWindowsVersion() (string, int) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE)
	if err != nil {
		return "unknown", 0
	}
	defer k.Close()
	
	productName, _, _ := k.GetStringValue("ProductName")
	currentBuild, _, _ := k.GetStringValue("CurrentBuild")
	// ubr, _, _ := k.GetIntegerValue("UBR")
	
	build := 0
	if currentBuild != "" {
		build, _ = strconv.Atoi(currentBuild)
	}
	
	return productName, build
}

// harvestDistrib returns Windows distribution information
func harvestDistrib() data.Distrib {
	productName, build := detectWindowsVersion()
	
	// Determine if Windows 10 or 11
	var id string
	
	if strings.Contains(strings.ToLower(productName), "windows 11") {
		id = "windows11"
	} else if strings.Contains(strings.ToLower(productName), "windows 10") {
		id = "windows10"
	} else {
		id = "windows"
	}

	version := data.NewVersion(build, 0, 0, "")
	
	return data.NewDistrib(id, make([]string, 0), productName, version)
}

// -----------------------

func harvestSysLib() data.SysLib {
	return data.NewSysLib(data.SysLibUnknown, "i don't care")
}
