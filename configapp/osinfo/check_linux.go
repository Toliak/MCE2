package osinfo

import (
	"fmt"
	"syscall"
	"time"
)

func isPseudoFsMounted(path string, magic int64) bool {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		// TODO: verbose log
		fmt.Printf("Error! %s\n", err)
		return false
	}

	// fmt.Printf("stat.Type=%x\n", stat.Type)
	return stat.Type == magic
}

func isProcMounted() bool {
	// procfs magic: 0x9fa0
	return isPseudoFsMounted("/proc", 0x9fa0)
}

func isSysMounted() bool {
	// sysfs magic: 0x62656572
	return isPseudoFsMounted("/sys", 0x62656572)
}

func isPriviledgedModeAvailable() bool {
	if IsRoot() {
		fmt.Println("Running the app as root is not recommended")
		time.Sleep(500 * time.Millisecond)
		return true
	}

	// TODO: sync that check with the platform/execwrap/ExecCommand
	// cache maybe or something like that
	if CommandExists("sudo") || CommandExists("pkexec") {
		return true
	}

	return false
}

func CheckPlatform() (bool, []string) {
	errors := make([]string, 0)

	if !isProcMounted() {
		errors = append(errors, "procfs (/proc) is not mounted")
	}
	if !isSysMounted() {
		errors = append(errors, "sysfs (/sys) is not mounted")
	}
	if !isPriviledgedModeAvailable() {
		errors = append(errors, "unable to find program to evaluate the privileges (sudo or pkexec)")
	}

	return len(errors) == 0, errors
}
