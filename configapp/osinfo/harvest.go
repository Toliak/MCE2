package osinfo

import (
	"runtime"

	"github.com/toliak/mce/osinfo/data"
)

// See https://github.com/golang/go/blob/9777ceceec8fee294d038182739cab7c845ad2d1/src/internal/syslist/syslist.go#L58
func harvestCPUArch() data.CPUArch {
	raw := runtime.GOARCH
	return data.ParseCPUArch(raw)
}

// See https://github.com/golang/go/blob/9777ceceec8fee294d038182739cab7c845ad2d1/src/internal/syslist/syslist.go#L58
func harvestOSType() data.OSType {
	raw := runtime.GOOS
	return data.ParseOsType(raw)
}

