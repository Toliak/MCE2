package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type SharedLocalConfig struct {}

var _ tb.Tegn = (*SharedLocalConfig)(nil)

func NewTegnSharedLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &SharedLocalConfig{}
	}
}

func getSharedLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-cfg.shared")
}

// GetID implements [tb.Tegn].
func (p *SharedLocalConfig) GetID() string {
	return "cfg-local-shared"
}

// GetName implements [tb.Tegn].
func (p *SharedLocalConfig) GetName() string {
	return "Local MCE2 config for shared environment variables and functions"
}

// GetDescription implements [tb.Tegn].
func (p *SharedLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing shared environment variables, functions, and settings for multiple shells
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *SharedLocalConfig) GetBeforeIDs() []string {
	// return []string{"base-cfg-zsh", "cfg-local-zsh", "cfg-local-bash"}
	return []string{}
}

// GetParameters implements [tb.Tegn].
func (p *SharedLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getSharedLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}

	// TODO: move PATH and EDITOR here
}

// GetFeatures implements [tb.Tegn].
func (p *SharedLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:shared-local"): true, 
	}
}

func (p *SharedLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getSharedLocalConfigPath(osInfo)
	return platform.FileEntryExists(path)
}

func (p *SharedLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	localConfigPath := getSharedLocalConfigPath(osInfo)
	err := MkdirAllParent(localConfigPath)
	if err != nil {
		return fmt.Errorf(" MkdirAll parent '%s' error: %w", localConfigPath, err)
	}

	{
		outputFile, err := os.Create(localConfigPath)
		if err != nil {
			return fmt.Errorf(" failed to create config file %s: %w", localConfigPath, err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString("###### Shared Local MCE2 config file.\n###### Managed by the MCE2\n\n")
		if err != nil {
			return fmt.Errorf(" failed to write to config file %s: %w", localConfigPath, err)
		}
	}

	if already["cfg:bash-local"] {
		bashLocalConfigPath := getBashLocalConfigPath(osInfo)
		platform.AppendFilepathString(
			bashLocalConfigPath,
			fmt.Sprintf("\n# <BEGIN> MCE2 shared config\nsource '%s'\n# <END> MCE2 shared config\n\n", localConfigPath),
		)
	}
	if already["cfg:zsh-local"] {
		bashLocalConfigPath := getZshLocalConfigPath(osInfo)
		platform.AppendFilepathString(
			bashLocalConfigPath,
			fmt.Sprintf("\n# <BEGIN> MCE2 shared config\nsource '%s'\n# <END> MCE2 shared config\n\n", localConfigPath),
		)
	}

	return nil
}