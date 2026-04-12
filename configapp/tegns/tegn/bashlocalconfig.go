package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type BashLocalConfig struct {}

var _ tb.Tegn = (*BashLocalConfig)(nil)

func NewTegnBashLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &BashLocalConfig{}
	}
}

func getBashLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-cfg.bash")
}

func getBashrcPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	bashrcPath := filepath.Join(userHomeDir, ".bashrc")
	return bashrcPath, err
}

// GetID implements [tb.Tegn].
func (p *BashLocalConfig) GetID() string {
	return "cfg-local-bash"
}

// GetName implements [tb.Tegn].
func (p *BashLocalConfig) GetName() string {
	return "Local MCE2 config for Bash"
}

// GetDescription implements [tb.Tegn].
func (p *BashLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing plugins, aliases, etc add by the mce2 config app
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *BashLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *BashLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *BashLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// if err := tb.CheckFeatures(before, []tb.TegnFeature{""}); err != nil {
	// 	return tb.NewTegnNotAvailable(err.Error())
	// }

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *BashLocalConfig) GetBeforeIDs() []string {
	// return []string{"base-cfg-bash"}
	return []string{}
}

// GetParameters implements [tb.Tegn].
func (p *BashLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getBashLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *BashLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:bash-local"): true, 
	}
}

func (p *BashLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getBashLocalConfigPath(osInfo)
	return platform.FileEntryExists(path)
}

func (p *BashLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	localConfigPath := getBashLocalConfigPath(osInfo)

	{
		outputFile, err := os.Create(localConfigPath)
		if err != nil {
			return fmt.Errorf("ExecInstall failed to create config file %s: %w", localConfigPath, err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString("###### Config at the end of bashrc initialization.\n###### Managed by the MCE2\n\n")
		if err != nil {
			return fmt.Errorf("ExecInstall failed to write to config file %s: %w", localConfigPath, err)
		}

		if already["cfg:mce2"] {
			_, err = outputFile.WriteString(
				fmt.Sprintf("# <BEGIN> MCE2 config\nsource '%s'\n# <END> MCE2 config\n", filepath.Join(osInfo.MainInstallDir, "shell", "bash.sh")),
			)
			if err != nil {
				return fmt.Errorf("ExecInstall failed to write to config file %s: %w", localConfigPath, err)
			}
		}
	}

	bashrcPath, err := getBashrcPath()
	if err != nil {
		return fmt.Errorf("ExecInstall failed to get bashrc path: %w", err)
	}

	platform.AppendFilepathString(
		bashrcPath,
		fmt.Sprintf(`\n\n# <BEGIN> MCE2 local config\nsource '%s'\n# <END> MCE2 local config\n\n"`, localConfigPath),
	)

	return nil
}
