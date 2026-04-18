package tegn

import (
	"fmt"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/zsh-users/zsh-syntax-highlighting"),
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
			tb.WithDefaultValue(getInstallDirZshSyntaxHighlight(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshSyntaxHighlight) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-syntax-highlighting"): true,
	}
}

func (p *ZshSyntaxHighlight) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshSyntaxHighlight(osInfo)
	return platform.FileEntryExists(path)
}

func (p *ZshSyntaxHighlight) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]

	path := getInstallDirZshSyntaxHighlight(osInfo)
	repo, err := git.PlainClone(
		path, 
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
	
	// Assume that we have local config installed
	zshLocalPreConfig := getZshLocalPreOhMyZshConfigPath(osInfo)
	platform.AppendFilepathString(
		zshLocalPreConfig,
		"\n\n# <BEGIN> zsh-syntax-highlighting\nplugins+=(zsh-syntax-highlighting)\n# <END> zsh-syntax-highlighting\n\n",
	)

	return nil
}
