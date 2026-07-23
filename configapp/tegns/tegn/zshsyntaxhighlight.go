package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshSyntaxHighlight struct {
	// TODO:
}

var _ tb.Tegn = (*ZshSyntaxHighlight)(nil)

func NewTegnZshSyntaxHighlightBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshSyntaxHighlight{}
	}
}

func getInstallDirZshSyntaxHighlight(osInfo tb.OSInfoExt) string {
	return filepath.Join(getInstallDirZshBaseConfig(osInfo), "custom", "plugins", "zsh-syntax-highlighting")
}

// GetID implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetID() string {
	return "cfg-zsh-syntax-highlighting"
}

// GetName implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetName() string {
	return "ZSH Syntax Highlighting"
}

// GetDescription implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetDescription() string {
	return `cfg-zsh-syntax-highlighting
Oh-my-zsh plugin

URL: https://github.com/zsh-users/zsh-syntax-highlighting
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetAvailability(
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
func (p *ZshSyntaxHighlight) GetBeforeIDs() []string {
	return []string{"base-cfg-zsh", "cfg-local-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	gitParams := GetGitCloneParameters("https://github.com/zsh-users/zsh-syntax-highlighting", getInstallDirZshSyntaxHighlight(osInfo))
	return gitParams
}

// GetFeatures implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-syntax-highlighting"): true,
	}
}

func (p *ZshSyntaxHighlight) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshSyntaxHighlight(osInfo)
	return platform.FileExists(path)
}

func (p *ZshSyntaxHighlight) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirZshSyntaxHighlight(osInfo)

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	// Assume that we have local config installed
	zshLocalPreConfig := getZshLocalPreOhMyZshConfigPath(osInfo)
	platform.AppendFilepathString(
		zshLocalPreConfig,
		"\n\n# <BEGIN> zsh-syntax-highlighting\nplugins+=(zsh-syntax-highlighting)\n# <END> zsh-syntax-highlighting\n\n",
	)

	return nil
}

func (p *ZshSyntaxHighlight) ExecUninstall(osInfo tb.OSInfoExt) error {

	// Remove the plugin registration from zsh local pre-config
	zshLocalPreConfig := getZshLocalPreOhMyZshConfigPath(osInfo)
	if platform.FileExists(zshLocalPreConfig) {
		err := removeConfigBlockFromFile(zshLocalPreConfig, "zsh-syntax-highlighting")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", zshLocalPreConfig, err)
		}
	}

	// Remove the cloned plugin
	installPath := getInstallDirZshSyntaxHighlight(osInfo)
	if platform.FileExists(installPath) {
		err := os.RemoveAll(installPath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
		}
	}

	return nil
}