// Extension for the [data.OSInfo]
package tegnbuilder

import (
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
)

type AvailablePackagesMap map[string] bool

// Extended OSInfo
type OSInfoExt struct {
	data.OSInfo

	// Packages, that are available from the package managers
	AvailableManagerPackages AvailablePackagesMap

	// Directory to clone the MCE2 project and to clone.
	// Should contain absolute path
	MainInstallDir string

	// Path inside the MainInstallDir.
	// Where the configs and other data will be put
	DataDir string

	// MCE2 repository URL
	MceRepositoryURL string

	// MCE2 branch
	MceRepositoryBranch string
}

func (data *OSInfoExt) GetFullDataDir() string {
	return filepath.Join(data.MainInstallDir, data.DataDir)
}