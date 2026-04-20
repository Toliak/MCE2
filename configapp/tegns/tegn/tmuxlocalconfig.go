package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type TmuxLocalConfig struct {}

var _ tb.Tegn = (*TmuxLocalConfig)(nil)

func NewTegnTmuxLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &TmuxLocalConfig{}
	}
}

func getTmuxLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-tmux.conf")
}

// GetID implements [tb.Tegn].
func (p *TmuxLocalConfig) GetID() string {
	return "cfg-local-tmux"
}

// GetName implements [tb.Tegn].
func (p *TmuxLocalConfig) GetName() string {
	return "Local MCE2 config for Tmux"
}

// GetDescription implements [tb.Tegn].
func (p *TmuxLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing plugins, aliases, etc add by the mce2 config app
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *TmuxLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *TmuxLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *TmuxLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// fmt.Printf("before: %#v\n", before)
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:oh-my-tmux"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *TmuxLocalConfig) GetBeforeIDs() []string {
	return []string{"base-cfg-tmux"}
}

// GetParameters implements [tb.Tegn].
func (p *TmuxLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getTmuxLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"pre-path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getTmuxLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *TmuxLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:tmux-local"): true, 
	}
}

func (p *TmuxLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getTmuxLocalConfigPath(osInfo)
	return platform.FileEntryExists(path)
}

func (p *TmuxLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	localConfigPath := getTmuxLocalConfigPath(osInfo)

	err := MkdirAllParent(localConfigPath)
	if err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", localConfigPath, err)
	}

	// Create the local config
	{
		outputFile, err := os.Create(localConfigPath)
		if err != nil {
			return fmt.Errorf("failed to create config file '%s': %w", localConfigPath, err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString("###### Config inserted after oh-my-tmux initialization.\n###### Managed by MCE2\n###### DO NOT EDIT. This file may be overwritten or removed at any time\n\n")
		if err != nil {
			return fmt.Errorf("failed to write to config file '%s': %w", localConfigPath, err)
		}

		mce2tmuxConfig := filepath.Join(osInfo.MainInstallDir, "tmux", "local.conf")
		if already["cfg:mce2"] && platform.FileEntryExists(mce2tmuxConfig) {
			_, err = fmt.Fprintf(
				outputFile,
				"# <BEGIN> MCE2 config\nsource '%s'\n# <END> MCE2 config\n", 
				mce2tmuxConfig,
			)
			if err != nil {
				return fmt.Errorf("failed to write to config file '%s': %w", localConfigPath, err)
			}
		}
	}

	// Insert the local config entry into the tmux.conf
	tmuxConfPath, err := getTmuxConfigPath(osInfo)
	if err != nil {
		return fmt.Errorf("getTmuxConfigPath error: %w", err)
	}
	err = platform.AppendFilepathString(
		tmuxConfPath,
		fmt.Sprintf("\n\n# <BEGIN> MCE2 local config\nsource '%s'\n# <END> MCE2 local config\n\n", localConfigPath),
	)

	return nil
}

func (p *TmuxLocalConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the entry from tmux.conf
	tmuxConfPath, err := getTmuxConfigPath(osInfo)
	if err != nil {
		return fmt.Errorf("getTmuxConfigPath error: %w", err)
	}

	if platform.FileEntryExists(tmuxConfPath) {
		err = removeConfigBlockFromFile(tmuxConfPath, "MCE2 local config")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", tmuxConfPath, err)
		}
	}

	// Remove the local config file
	localConfigPath := getTmuxLocalConfigPath(osInfo)
	if platform.FileEntryExists(localConfigPath) {
		err := os.Remove(localConfigPath)
		if err != nil {
			return fmt.Errorf("os.Remove error '%s': %w", localConfigPath, err)
		}
	}

	return nil
}
