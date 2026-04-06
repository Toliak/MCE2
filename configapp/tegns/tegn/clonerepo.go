package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type CloneRepo struct {
	installDir string
}

var _ tb.Tegn = (*CloneRepo)(nil)

func NewTegnCloneRepoBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &CloneRepo{}
	}
}

// GetID implements [tb.Tegn].
func (p *CloneRepo) GetID() string {
	return "mce2-repo"
}

// GetName implements [tb.Tegn].
func (p *CloneRepo) GetName() string {
	return "Clone MCE2"
}

// GetDescription implements [tb.Tegn].
func (p *CloneRepo) GetDescription() string {
	return `Make Configuration Easier 2

URL: https://github.com/toliak/mce2
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *CloneRepo) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *CloneRepo) GetAvailableOsType() *[]data.OSTypeE {
	return &[]data.OSTypeE{
		data.OSTypeLinux,
	}
}

// GetAvailability implements [tb.Tegn].
func (p *CloneRepo) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// TODO: do we really need the git package to be installed????
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"pkg:git"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *CloneRepo) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tb.Tegn].
func (p *CloneRepo) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL (read-only)"),
			tb.WithDefaultValue(osInfo.MceRepositoryURL),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository branch (read-only)"),
			tb.WithDefaultValue(osInfo.MceRepositoryBranch),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"install-path",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation path (read-only)"),
			tb.WithDefaultValue(osInfo.MainInstallDir),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"data-path",
			"Installation data path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Installation data path (read-only)"),
			tb.WithDefaultValue(osInfo.GetFullDataDir()),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *CloneRepo) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:mce2"): true,
}
}

func (p *CloneRepo) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := filepath.Join(osInfo.MainInstallDir, ".git")
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		// TODO: log the error somewhere
		return false
	} else {
		return true
	}
}

func (p *CloneRepo) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	installPath := osInfo.MainInstallDir

	repo, err := git.PlainClone(installPath, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		NoCheckout: true,
	})
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

	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

// func (p *ZshBaseConfig)  ExecUninstall() error {

// }
