package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshFavoriteDir struct {}

var _ tb.Tegn = (*ZshFavoriteDir)(nil)

func NewTegnZshFavoriteDirBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshFavoriteDir{}
	}
}

func getInstallDirZshFavoriteDir(osInfo tb.OSInfoExt) string {
	return filepath.Join(getInstallDirZshBaseConfig(osInfo), "custom", "plugins", "zsh_favorite_dir")
}

// GetID implements [tb.Tegn].
func (p *ZshFavoriteDir) GetID() string {
	return "cfg-zsh-favorite-dir"
}

// GetName implements [tb.Tegn].
func (p *ZshFavoriteDir) GetName() string {
	return "ZSH Favorite Dir"
}

// GetDescription implements [tb.Tegn].
func (p *ZshFavoriteDir) GetDescription() string {
	return `ZSH Favorite Dir
Oh-my-zsh plugin / widget

URL: https://github.com/Toliak/zsh-favorite-dir
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshFavoriteDir) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshFavoriteDir) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshFavoriteDir) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:oh-my-zsh", "cfg:zsh-local"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshFavoriteDir) GetBeforeIDs() []string {
	return []string{"base-cfg-zsh", "cfg-local-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshFavoriteDir) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	gitParams := GetGitCloneParameters(
		"https://github.com/Toliak/zsh-favorite-dir", 
		getInstallDirZshFavoriteDir(osInfo),
	)
	return gitParams
}

// GetFeatures implements [tb.Tegn].
func (p *ZshFavoriteDir) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-favorite-dir"): true,
	}
}

func (p *ZshFavoriteDir) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshFavoriteDir(osInfo)
	return platform.FileExists(path)
}

func (p *ZshFavoriteDir) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirZshFavoriteDir(osInfo)

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	// Assume that we have local config installed
	zshLocalPreConfig := getZshLocalPreOhMyZshConfigPath(osInfo)
	platform.AppendFilepathString(
		zshLocalPreConfig,
		"\n\n# <BEGIN> zsh_favorite_dir\nplugins+=(zsh_favorite_dir)\n# <END> zsh_favorite_dir\n\n",
	)

	return nil
}

func (p *ZshFavoriteDir) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the plugin registration from zsh local pre-config
	zshLocalPreConfig := getZshLocalPreOhMyZshConfigPath(osInfo)
	if platform.FileExists(zshLocalPreConfig) {
		err := removeConfigBlockFromFile(zshLocalPreConfig, "zsh_favorite_dir")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", zshLocalPreConfig, err)
		}
	}

	// Remove the cloned plugin
	installPath := getInstallDirZshFavoriteDir(osInfo)
	if platform.FileExists(installPath) {
		err := os.RemoveAll(installPath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
		}
	}

	return nil
}
