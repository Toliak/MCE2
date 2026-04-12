package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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
	return filepath.Join(osInfo.GetFullDataDir(), "vimrc-zmix")
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
	return "cfg-vim-zmix"
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
			tb.WithDescription("Repository branch"),
			tb.WithDefaultValue("master"),
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
		tb.NewTegnParameter(
			"vimrc-backup",
			"Do vimrc backup",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Backup current .vimrc configuration?"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *UltimateVim) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:vim-ultimate"): true,
	}
}

func (p *UltimateVim) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := filepath.Join(getInstallDirUltimateVim(osInfo), ".git")
	return platform.FileEntryExists(path)
}

func (p *UltimateVim) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	path := getInstallDirUltimateVim(osInfo)
	vimrcBackup := tb.TegnParameterToBool(params["vimrc-backup"])

	repo, err := git.PlainClone(
		path, 
		defaultGitCloneOptions(func(v *git.CloneOptions) {
			v.URL = url
		}),
	)
	if err != nil {
		return fmt.Errorf("ExecInstall PlainClone error: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("ExecInstall Worktree error: %w", err)
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  false,
	}
	err = w.Checkout(&branchCoOpts)
	if err != nil {
		return fmt.Errorf("ExecInstall Checkout error: %w", err)
	}

	vimrcOrigPath, err := getVimrcPath()
	if err != nil {
		return fmt.Errorf("ExecInstall failed to get vimrc path: %w", err)
	}

	if vimrcBackup && platform.FileEntryExists(vimrcOrigPath) {
		err := platform.CopyFile(vimrcOrigPath, vimrcOrigPath + ".backup-mce")
		if err != nil {
			return fmt.Errorf("ExecInstall vimrc backup error: %w", err)
		}
	}

	templateFile := filepath.Join(path, "vimrcs", "basic.vim")
	err = platform.CopyFile(templateFile, vimrcOrigPath)
	if err != nil {
		return fmt.Errorf("ExecInstall platform.CopyFile error: %w", err)
	}

	return nil
}
