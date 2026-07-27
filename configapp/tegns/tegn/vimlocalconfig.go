package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type VimLocalConfig struct {}

var _ tb.Tegn = (*VimLocalConfig)(nil)

func NewTegnVimLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &VimLocalConfig{}
	}
}

func getVimLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-vim.conf")
}

// GetID implements [tb.Tegn].
func (p *VimLocalConfig) GetID() string {
	return "cfg-local-vim"
}

// GetName implements [tb.Tegn].
func (p *VimLocalConfig) GetName() string {
	return "Local MCE2 config for Vim"
}

// GetDescription implements [tb.Tegn].
func (p *VimLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing plugins, aliases, etc add by the mce2 config app
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *VimLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *VimLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *VimLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// fmt.Printf("before: %#v\n", before)
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:vim-ultimate"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *VimLocalConfig) GetBeforeIDs() []string {
	return []string{"base-cfg-vim"}
}

// GetParameters implements [tb.Tegn].
func (p *VimLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getVimLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *VimLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:vim-local"): true, 
	}
}

func (p *VimLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getVimLocalConfigPath(osInfo)
	return platform.FileExists(path)
}

func (p *VimLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	localConfigPath := getVimLocalConfigPath(osInfo)

	err := platform.MkdirAllParent(localConfigPath)
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

		_, err = outputFile.WriteString("\"\"\"\"\"\" Config inserted after ultimate vim initialization.\n\"\"\"\"\"\" Managed by MCE2\n\"\"\"\"\"\" DO NOT EDIT. This file may be overwritten or removed at any time\n\n")
		if err != nil {
			return fmt.Errorf("failed to write to config file '%s': %w", localConfigPath, err)
		}

		mce2vimConfig := filepath.Join(osInfo.MainInstallDir, "vim", "local.conf")
		if already["cfg:mce2"] && platform.FileExists(mce2vimConfig) {
			_, err = fmt.Fprintf(
				outputFile,
				"\" <BEGIN> MCE2 config\nsource %s\n\" <END> MCE2 config\n", 
				mce2vimConfig,
			)
			if err != nil {
				return fmt.Errorf("failed to write to config file '%s': %w", localConfigPath, err)
			}
		}
	}

	// Insert the local config entry
	vimrcConfPath := getVimrcPath(osInfo)

	err = platform.AppendFilepathString(
		vimrcConfPath,
		fmt.Sprintf("\n\n\" <BEGIN> MCE2 local config\nsource %s\n\" <END> MCE2 local config\n\n", localConfigPath),
	)

	return nil
}

func (p *VimLocalConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the entry from .vimrc
	vimrcPath := getVimrcPath(osInfo)

	if platform.FileExists(vimrcPath) {
		err := removeConfigBlockFromFile(vimrcPath, "MCE2 local config")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", vimrcPath, err)
		}
	}

	// Remove the local config file
	localConfigPath := getVimLocalConfigPath(osInfo)
	if platform.FileExists(localConfigPath) {
		err := os.Remove(localConfigPath)
		if err != nil {
			return fmt.Errorf("os.Remove error '%s': %w", localConfigPath, err)
		}
	}

	return nil
}
