package tegn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshBaseConfig struct {}

var _ tb.Tegn = (*ZshBaseConfig)(nil)

func NewTegnZshBaseConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshBaseConfig{}
	}
}

func getInstallDirZshBaseConfig(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "oh-my-zsh")
}

func getZshrcPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	zshrcPath := filepath.Join(userHomeDir, ".zshrc")
	return zshrcPath, err
}

// GetID implements [tb.Tegn].
func (p *ZshBaseConfig) GetID() string {
	return "base-cfg-zsh"
}

// GetName implements [tb.Tegn].
func (p *ZshBaseConfig) GetName() string {
	return "Oh My ZSH"
}

// GetDescription implements [tb.Tegn].
func (p *ZshBaseConfig) GetDescription() string {
	return `Oh My ZSH

URL: https://github.com/ohmyzsh/ohmyzsh
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshBaseConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshBaseConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshBaseConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"pkg:zsh", "cfg:mce2"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshBaseConfig) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetParameters implements [tb.Tegn].
func (p *ZshBaseConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	gitParams := GetGitCloneParameters("https://github.com/ohmyzsh/ohmyzsh", getInstallDirZshBaseConfig(osInfo))
	return gitParams
}

// GetFeatures implements [tb.Tegn].
func (p *ZshBaseConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-base"): true, 
		tb.TegnFeature("cfg:oh-my-zsh"): true,
	}
}

func (p *ZshBaseConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshBaseConfig(osInfo)
	return platform.FileExists(path)
}

var zshSourceOhMyZshLine = regexp.MustCompile(`^\s*source\s+.+/oh-my-zsh.sh\s*$`)

func getOhMyZshTemplateConfig(ohMyZshDir string) string {
	return fmt.Sprintf(`
# <BEGIN> oh-my-zsh config (autogen mce2)
export ZSH='%s'
ZSH_THEME="robbyrussell"

plugins=(git)

source $ZSH/oh-my-zsh.sh
# <END> oh-my-zsh config (autogen mce2)
`,
		ohMyZshDir,
	)
}

// TODO: shouldn't this thing be idempotent
func (p *ZshBaseConfig) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirZshBaseConfig(osInfo)

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err := cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}
	
	zshrcOrigPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("failed to get zshrc path: %w", err)
	}
	err = platform.AppendFilepathString(zshrcOrigPath, getOhMyZshTemplateConfig(path))
	if err != nil {
		return fmt.Errorf("AppendFilepathString error '%s': %w", zshrcOrigPath, err)
	}

	return nil
}

func (p *ZshBaseConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the configuration block from .zshrc
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("getZshrcPath error: %w", err)
	}

	if platform.FileExists(zshrcPath) {
		err = removeConfigBlockFromFile(zshrcPath, "oh-my-zsh config (autogen mce2)")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", zshrcPath, err)
		}
	}

	// Remove the Oh My ZSH installation
	ohmyzshDir := getInstallDirZshBaseConfig(osInfo)
	if platform.FileExists(ohmyzshDir) {
		err := os.RemoveAll(ohmyzshDir)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", ohmyzshDir, err)
		}
	}

	return nil
}