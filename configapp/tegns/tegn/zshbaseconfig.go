package tegn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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
	// TODO: cache that maybe?

	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/ohmyzsh/ohmyzsh"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch"),
			tb.WithDefaultValue("master"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"install-dir",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation path (read-only)"),
			tb.WithDefaultValue(getInstallDirZshBaseConfig(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
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
	return platform.FileEntryExists(path)
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

	ohmyzshDir := getInstallDirZshBaseConfig(osInfo)
	err := MkdirAllParent(ohmyzshDir)
	if err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", ohmyzshDir, err)
	}

	repo, err := git.PlainClone(
		ohmyzshDir, 
		defaultGitCloneOptions(func (v *git.CloneOptions) {
			v.URL = url
		}),
	)
	if err != nil {
		return fmt.Errorf("PlainClone error: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Worktree error: %w", err)
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  false,
	}
	err = w.Checkout(&branchCoOpts)
	if err != nil {
		return fmt.Errorf("Checkout error: %w", err)
	}
	
	zshrcOrigPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("failed to get zshrc path: %w", err)
	}
	err = platform.AppendFilepathString(zshrcOrigPath, getOhMyZshTemplateConfig(ohmyzshDir))
	if err != nil {
		return fmt.Errorf("AppendFilepathString error '%s': %w", zshrcOrigPath, err)
	}

	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

func (p *ZshBaseConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the configuration block from .zshrc
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("getZshrcPath error: %w", err)
	}

	if platform.FileEntryExists(zshrcPath) {
		err = removeConfigBlockFromFile(zshrcPath, "oh-my-zsh config (autogen mce2)")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", zshrcPath, err)
		}
	}

	// Remove the Oh My ZSH installation
	ohmyzshDir := getInstallDirZshBaseConfig(osInfo)
	if platform.FileEntryExists(ohmyzshDir) {
		err := os.RemoveAll(ohmyzshDir)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", ohmyzshDir, err)
		}
	}

	return nil
}