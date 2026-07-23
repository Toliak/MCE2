package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshAutoSuggestions struct {}

var _ tb.Tegn = (*ZshAutoSuggestions)(nil)

func NewTegnZshAutoSuggestionsBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshAutoSuggestions{}
	}
}

func getInstallDirZshAutoSuggestions(osInfo tb.OSInfoExt) string {
	return filepath.Join(getInstallDirZshBaseConfig(osInfo), "custom", "plugins", "zsh-autosuggestions")
}

// GetID implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetID() string {
	return "cfg-zsh-autosuggestions"
}

// GetName implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetName() string {
	return "ZSH Autosuggestions"
}

// GetDescription implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetDescription() string {
	return `cfg-zsh-autosuggestions

URL: https://github.com/zsh-users/zsh-autosuggestions
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailability(
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
func (p *ZshAutoSuggestions) GetBeforeIDs() []string {
	return []string{"cfg-local-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/zsh-users/zsh-autosuggestions"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch"),
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
			tb.WithDefaultValue(getInstallDirZshAutoSuggestions(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-autosuggestions"): true, 
	}
}

func (p *ZshAutoSuggestions) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshAutoSuggestions(osInfo)
	return platform.FileExists(path)
}

func (p *ZshAutoSuggestions) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// Assume we have `cfg:zsh-local` as a requirement
	// localConfigPath := getZshLocalConfigPath(osInfo)
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)

	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirZshAutoSuggestions(osInfo)

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	err = platform.AppendFilepathString(
		localPreConfigPath,
		"\n# <BEGIN> MCE2 zsh-autosuggestions\nplugins+=(zsh-autosuggestions)\n# <END> MCE2 zsh-autosuggestions\n",
	)
	if err != nil {
		return fmt.Errorf("AppendFilepathString %s: %w", localPreConfigPath, err)
	}

	return nil
}

func (p *ZshAutoSuggestions) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the plugin registration from zsh local pre-config
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)
	if platform.FileExists(localPreConfigPath) {
		err := removeConfigBlockFromFile(localPreConfigPath, "MCE2 zsh-autosuggestions")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", localPreConfigPath, err)
		}
	}

		// Remove the cloned plugin
		installPath := getInstallDirZshAutoSuggestions(osInfo)
		if platform.FileExists(installPath) {
			err := os.RemoveAll(installPath)
			if err != nil {
				return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
			}
		}

	return nil
}