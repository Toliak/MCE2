package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type UltimateVim struct {
	// TODO:
}

var _ tb.Tegn = (*UltimateVim)(nil)

func NewTegnUltimateVimBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &UltimateVim{}
	}
}

func getInstallDirUltimateVim(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "vimrc-amix")
}

func getVimrcPath(osInfo tb.OSInfoExt) string {
	vimrcPath := filepath.Join(osInfo.HomeDir, ".vimrc")
	return vimrcPath
}

// GetID implements [tb.Tegn].
func (p *UltimateVim) GetID() string {
	return "base-cfg-vim"
}

// GetName implements [tb.Tegn].
func (p *UltimateVim) GetName() string {
	return "UltimateVim"
}

// GetDescription implements [tb.Tegn].
func (p *UltimateVim) GetDescription() string {
	return `UltimateVim (basic version)

URL: https://github.com/amix/vimrc
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *UltimateVim) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *UltimateVim) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *UltimateVim) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *UltimateVim) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetParameters implements [tb.Tegn].
func (p *UltimateVim) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	gitParams := GetGitCloneParameters("https://github.com/amix/vimrc", getInstallDirUltimateVim(osInfo))
	return gitParams
}

// GetFeatures implements [tb.Tegn].
func (p *UltimateVim) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:vim-base"): true,
		tb.TegnFeature("cfg:vim-ultimate"): true,
	}
}

func (p *UltimateVim) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := filepath.Join(getInstallDirUltimateVim(osInfo), ".git")
	return platform.FileExists(path)
}

func (p *UltimateVim) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirUltimateVim(osInfo)

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	vimrcOrigPath := getVimrcPath(osInfo)

	// Insert the config entry
	templateFile := filepath.Join(path, "vimrcs", "basic.vim")
	if !platform.FileExists(templateFile) {
		return fmt.Errorf("templateFile '%s' does not exist", templateFile)
	}

	err = platform.AppendFilepathString(
		vimrcOrigPath, 
		fmt.Sprintf("\" <BEGIN> Ultimate vim config (autogen mce2)\nsource %s\n\" <END> Ultimate vim config (autogen mce2)\n", templateFile),
	)
	if err != nil {
		return fmt.Errorf("AppendFilepathString error '%s': %w", vimrcOrigPath, err)
	}

	return nil
}

func (p *UltimateVim) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the configuration block from .vimrc
	vimrcPath := getVimrcPath(osInfo)

	if platform.FileExists(vimrcPath) {
		err := removeConfigBlockFromFile(vimrcPath, "Ultimate vim config (autogen mce2)")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", vimrcPath, err)
		}
	}

	// Remove the cloned repository
	installPath := getInstallDirUltimateVim(osInfo)
	if platform.FileExists(installPath) {
		err := os.RemoveAll(installPath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
		}
	}

	return nil
}