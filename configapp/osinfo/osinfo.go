package osinfo

import "github.com/toliak/mce/osinfo/data"

func Harvest() data.OsInfo {
	distrib := harvestDistrib()

	return data.NewOsInfo(
		data.WithArch(harvestCPUArch()),
		data.WithOSType(harvestOSType()),
		data.WithPkgManager(
			harvestPkgManager(&distrib),
		),
		data.WithSysLib(harvestSysLib()),
		// data.WithKernelVersion(),
		data.WithDistrib(distrib),
	)
}