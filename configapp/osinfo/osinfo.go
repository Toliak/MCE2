package osinfo

import "github.com/toliak/mce/osinfo/data"

func Harvest() data.OsInfo {
	return data.NewOsInfo(
		data.WithArch(harvestCPUArch()),
		data.WithOSType(harvestOSType()),
		data.WithPkgManager(harvestPkgManager()),
		data.WithSysLib(harvestSysLib()),
		// data.WithKernelVersion(),
		data.WithDistrib(harvestDistrib()),
	)
}