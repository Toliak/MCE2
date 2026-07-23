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

func getVimrcPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	vimrcPath := filepath.Join(userHomeDir, ".vimrc")
	return vimrcPath, err
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
	return []tb.TegnParameter {
		// TODO: encapsulate Git clone default parameter pack:
		// repo-url
		// repo-branch
		// use-shallow-clone
		// install-dir

		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/amix/vimrc"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch, leave empty to use the default branch"),
			tb.WithDefaultValue(""),
			tb.WithAvailabilityTrue(),
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
			"install-dir",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation path (read-only)"),
			tb.WithDefaultValue(getInstallDirUltimateVim(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
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

	vimrcOrigPath, err := getVimrcPath()
	if err != nil {
		return fmt.Errorf("failed to get vimrc path: %w", err)
	}

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
	vimrcPath, err := getVimrcPath()
	if err != nil {
		return fmt.Errorf("getVimrcPath error: %w", err)
	}

	if platform.FileExists(vimrcPath) {
		err = removeConfigBlockFromFile(vimrcPath, "Ultimate vim config (autogen mce2)")
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