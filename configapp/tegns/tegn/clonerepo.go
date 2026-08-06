package tegn

import (
	"fmt"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type CloneRepo struct {
	installDir string
}

var _ tb.Tegn = (*CloneRepo)(nil)

func NewTegnCloneRepoBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &CloneRepo{}
	}
}

// GetID implements [tb.Tegn].
func (p *CloneRepo) GetID() string {
	return "mce2-repo"
}

// GetName implements [tb.Tegn].
func (p *CloneRepo) GetName() string {
	return "Clone MCE2"
}

// GetDescription implements [tb.Tegn].
func (p *CloneRepo) GetDescription() string {
	return `Make Configuration Easier 2

URL: https://github.com/toliak/mce2
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *CloneRepo) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *CloneRepo) GetAvailableOsType() *[]data.OSTypeE {
	return &[]data.OSTypeE{
		data.OSTypeAndroid,
		data.OSTypeLinux,
		data.OSTypeDarwin,
	}
}

// GetAvailability implements [tb.Tegn].
func (p *CloneRepo) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// TODO: do we really need the git package to be installed????
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"pkg:git"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *CloneRepo) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tb.Tegn].
func (p *CloneRepo) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		// WARN: do not use GetGitCloneParameters here, because Repository branch here is read-only
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL (read-only)"),
			tb.WithDefaultValue(osInfo.MceRepositoryURL),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch (read-only)"),
			tb.WithDefaultValue(osInfo.MceRepositoryBranch),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"use-shallow-clone",
			"Shallow Clone",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Perform a shallow clone with depth 1 (--depth 1) for faster download"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"install-path",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation path (read-only)"),
			tb.WithDefaultValue(osInfo.MainInstallDir),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"data-path",
			"Installation data path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation data path (read-only)"),
			tb.WithDefaultValue(osInfo.GetFullDataDir()),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *CloneRepo) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:mce2"): true,
}
}

func (p *CloneRepo) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := filepath.Join(osInfo.MainInstallDir, ".git")
	return platform.FileExists(path)
}

func (p *CloneRepo) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := osInfo.MainInstallDir

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

func (p *CloneRepo) ExecUninstall(osInfo tb.OSInfoExt) error {
	// We do not remove the entire data folder, if there are any files inside


	// installPath := osInfo.MainInstallDir
	// // Assuming that installPath is the part of the data dir
	// dataDir := osInfo.GetFullDataDir()

	// if platform.FileExists(dataDir) {
	// 	// Check if directory is empty or only contains files we created
	// 	entries, err := os.ReadDir(dataDir)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to read data directory '%s': %w", dataDir, err)
	// 	}

	// 	// If directory is empty, remove it
	// 	if len(entries) == 0 {
	// 		// Remove the entire MCE2 installation directory
	// 		if platform.FileExists(installPath) {
	// 			err := os.RemoveAll(installPath)
	// 			if err != nil {
	// 				return fmt.Errorf("failed to remove installation directory '%s': %w", installPath, err)
	// 			}
	// 		}
	// 	}
	// }

	fmt.Printf("Tegn '%s' is not removable\n", p.GetID())
	return nil
}