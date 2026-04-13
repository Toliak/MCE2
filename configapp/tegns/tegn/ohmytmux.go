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

type OhMyTmux struct {
	// TODO:
}

var _ tb.Tegn = (*OhMyTmux)(nil)

func NewTegnOhMyTmuxBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &OhMyTmux{}
	}
}

func getInstallDirOhMyTmux(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "oh-my-tmux")
}

func getTmuxConfigDir(osInfo tb.OSInfoExt) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(userHomeDir, ".config", "tmux")
	return path, nil
}

// GetID implements [tb.Tegn].
func (p *OhMyTmux) GetID() string {
	return "base-cfg-tmux"
}

// GetName implements [tb.Tegn].
func (p *OhMyTmux) GetName() string {
	return "OhMyTmux"
}

// GetDescription implements [tb.Tegn].
func (p *OhMyTmux) GetDescription() string {
	return `Oh My Tmux! (~/.config/tmux)

URL: https://github.com/gpakosz/.tmux
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *OhMyTmux) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *OhMyTmux) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *OhMyTmux) GetAvailability(
	osInfo tb.OSInfoExt,
	before tb.TegnInstalledFeaturesMap,
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *OhMyTmux) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetParameters implements [tb.Tegn].
func (p *OhMyTmux) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter{
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/gpakosz/.tmux"),
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
			tb.WithDefaultValue(getInstallDirOhMyTmux(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"tmux-conf-backup",
			"Do tmux.conf backup",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Backup current tmux.conf configuration?"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *OhMyTmux) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:tmux-base"): true,
		tb.TegnFeature("cfg:oh-my-tmux"): true,
	}
}

func (p *OhMyTmux) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := filepath.Join(getInstallDirOhMyTmux(osInfo), ".git")
	return platform.FileEntryExists(path)
}

func (p *OhMyTmux) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	installPath := getInstallDirOhMyTmux(osInfo)
	err := MkdirAllParent(installPath)
	if err != nil {
		return fmt.Errorf("ExecInstall MkdirAll parent '%s' error: %w", installPath, err)
	}

	tmuxConfigDir, err := getTmuxConfigDir(osInfo)
	if err != nil {
		return fmt.Errorf("ExecInstall getTmuxConfigDir error: %w", err)
	}
	tmuxConfBackup := tb.TegnParameterToBool(params["tmux-conf-backup"])

	// Ensure target tmux config directory exists
	err = os.MkdirAll(tmuxConfigDir, 0755)
	if err != nil {
		return fmt.Errorf("ExecInstall failed to create tmux config dir: %w", err)
	}

	// Clone the repository
	repo, err := git.PlainClone(
		installPath,
		defaultGitCloneOptions(func(v *git.CloneOptions) {
			v.URL = url
		}),
	)
	if err != nil {
		return fmt.Errorf("ExecInstall PlainClone error: %w", err)
	}

	// Checkout the specified branch
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

	// Create symbolic link for .tmux.conf
	sourceTmuxConf := filepath.Join(installPath, ".tmux.conf")
	targetTmuxConf := filepath.Join(tmuxConfigDir, "tmux.conf")

	if tmuxConfBackup && platform.FileEntryExists(targetTmuxConf) {
		err := platform.CopyFile(targetTmuxConf, targetTmuxConf + ".backup-mce")
		if err != nil {
			return fmt.Errorf("ExecInstall targetTmuxConf backup error: %w", err)
		}
	}

	// Remove target if it exists (to avoid errors if it's a file or an existing symlink)
	if platform.FileEntryExists(targetTmuxConf) {
		err = os.Remove(targetTmuxConf)
		if err != nil {
			return fmt.Errorf("ExecInstall Remove targetTmuxConf error: %w", err)
		}
	}
	err = MkdirAllParent(targetTmuxConf)
	if err != nil {
		return fmt.Errorf("ExecInstall MkdirAll parent '%s' error: %w", targetTmuxConf, err)
	}

	tmuxConf, err := os.Create(targetTmuxConf)
	if err != nil {
		return fmt.Errorf("ExecInstall os.Create failed: %w", err)
	}
	defer tmuxConf.Close()

	fmt.Fprintf(
		tmuxConf, 
		"#### Auto-gen by MCE2\nset-environment -g TMUX_CONF \"%s\"\nset-environment -g TMUX_CONF_LOCAL \"%s.local\"\nsource-file \"%s\"\n\n", 
		sourceTmuxConf,
		targetTmuxConf,
		sourceTmuxConf,
	)

	// Copy .tmux.conf.local template to the config directory
	sourceLocalConf := filepath.Join(installPath, ".tmux.conf.local")
	targetLocalConf := filepath.Join(tmuxConfigDir, "tmux.conf.local")

	if tmuxConfBackup && platform.FileEntryExists(targetLocalConf) {
		err := platform.CopyFile(targetLocalConf, targetLocalConf + ".backup-mce")
		if err != nil {
			return fmt.Errorf("ExecInstall targetLocalConf backup error: %w", err)
		}
	}

	// Only copy if the target doesn't already exist, preserving user customizations
	err = platform.CopyFile(sourceLocalConf, targetLocalConf)
	if err != nil {
		return fmt.Errorf("ExecInstall platform.CopyFile for tmux.conf.local error: %w", err)
	}

	return nil
}