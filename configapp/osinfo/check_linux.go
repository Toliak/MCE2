package osinfo

import (
	"fmt"
	"syscall"
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

func CheckPlatform() (bool, []string) {
	errors := make([]string, 0)

	if !isProcMounted() {
		errors = append(errors, "procfs (/proc) is not mounted")
	}
	if !isSysMounted() {
		errors = append(errors, "sysfs (/sys) is not mounted")
	}

	return len(errors) == 0, errors
}
