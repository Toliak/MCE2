package inspector

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	osInfoHarvest "github.com/toliak/mce/osinfo/harvest"
	"github.com/toliak/mce/platform"
)

type InspectCheckError struct {
	ChecksFailed []string
}

var _ error = (*InspectCheckError)(nil)

func (e *InspectCheckError) Error() string {
	var builder strings.Builder

	builder.WriteString("Platform checks failed: ")
	for i, err := range e.ChecksFailed {
		builder.WriteString(err)

		if i != len(e.ChecksFailed) - 1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

type HarvestData struct {
	OSInfo *data.OSInfo
	AvailableManagerPackages *[]string
}

type InspectAndHarvestConfig struct {
	// Check enabled?
	Check bool

	// Harvest enabled?
	Harvest bool

	// Do the app need to update the packages using the OS package manager?
	PkgManagerUpdate bool
	
	// Do the app need to use the package manager to obtain available packages?
	PkgManagerGetAvailablePackages bool
}

func InspectAndHarvestConfigDefault() InspectAndHarvestConfig {
	return InspectAndHarvestConfig{
		Check: true,
		Harvest: true,
		PkgManagerUpdate: true,
		PkgManagerGetAvailablePackages: true,
	}
}

func InspectAndHarvest(options InspectAndHarvestConfig) (*HarvestData, error) {
	if options.Check {
		passed, errors := checkPlatform()
		if !passed {
			return nil, &InspectCheckError{ChecksFailed: errors}
		}
	}

	var osInfo *data.OSInfo
	if options.Harvest {
		v := osInfoHarvest.HarvestOSInfo()
		osInfo = &v
	}

	if options.PkgManagerUpdate {
		if osInfo == nil {
			fmt.Println("PkgManager repository metadata update will be skipped, due to the absence of the Harvested OSInfo")
		} else {
			if osInfo.PkgManager.V == data.PkgMgrUnknown {
				fmt.Printf("PkgManager repository metadata update will be skipped, due to the unknown PkgManager %s\n", osInfo.PkgManager)
			} else {
				err := platform.UpdateRepositories(&osInfo.PkgManager)
				if err != nil {
					return nil, fmt.Errorf("InspectAndHarvest UpdateRepositories error: %w", err)
				}
			}
		}
	}

	var availablePackages *[]string
	if options.PkgManagerGetAvailablePackages {
		if osInfo == nil {
			fmt.Println("PkgManager available packages obtaining will be skipped, due to the absence of the Harvested OSInfo")
		} else {
			if osInfo.PkgManager.V == data.PkgMgrUnknown {
				fmt.Printf("PkgManager available packages obtaining will be skipped, due to the unknown PkgManager %s\n", osInfo.PkgManager)
			} else {
				v, err := platform.GetAvailablePackages(&osInfo.PkgManager)
				if err != nil {
					return nil, fmt.Errorf("InspectAndHarvest GetAvailablePackages error: %w", err)
				}
				availablePackages = &v
			}
		}
	}
	
	result := HarvestData{
		OSInfo: osInfo,
		AvailableManagerPackages: availablePackages,
	}

	return &result, nil
}

