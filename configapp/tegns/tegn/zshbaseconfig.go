package tegn

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type ZshBaseConfig struct {
	installDir string
}

var _ tb.Tegn = (*ZshBaseConfig)(nil)

func NewTegnZshBaseConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshBaseConfig{}
	}
}

// GetID implements [tb.Tegn].
func (p *ZshBaseConfig) getInstallDir(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "oh-my-zsh")
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
	return &[]data.OSTypeE{
		data.OSTypeLinux,
	}
}

// GetAvailability implements [tb.Tegn].
func (p *ZshBaseConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before []tb.TegnFeature, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.TegnAvailability{
		Available: slices.Contains(before, "pkg:zsh") /*|| platform.CommandExists("zsh")*/,
		Reason:    fmt.Sprintf("Feature pkg:zsh not found"),
	}
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshBaseConfig) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tb.Tegn].
func (p *ZshBaseConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	// TODO: cache that maybe?

	userHomeDir, err := os.UserHomeDir()
	zshrcPath := filepath.Join(userHomeDir, ".zshrc")
	var zshrcAvailabilityReason string = ""
	if err != nil {
		zshrcAvailabilityReason = fmt.Sprintf(
			"Error while retrieving user's HOME directory: %s",
			err,
		)
	} else if !platform.FileEntryExists(zshrcPath) {
		zshrcAvailabilityReason = fmt.Sprintf(
			"File .zshrc not found (%s)",
			zshrcPath,
		)
	}


	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh-my-zsh Repository URL"),
			tb.WithDefaultValue("https://github.com/ohmyzsh/ohmyzsh"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh-my-zsh Repository branch"),
			tb.WithDefaultValue("master"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"zshrc-backup",
			tb.TegnParameterTypeString,
			tb.WithDescription("Backup current .zshrc configuration?"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailability(tb.TegnAvailability{
				Available: zshrcAvailabilityReason != "",
				Reason: zshrcAvailabilityReason,
			}),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshBaseConfig) GetFeatures() []tb.TegnFeature {
	return []tb.TegnFeature{
		tb.TegnFeature("cfg:zsh-base"), 
		tb.TegnFeature("cfg:oh-my-zsh"),
}
}

func (p *ZshBaseConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := p.getInstallDir(osInfo)
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

func (p *ZshBaseConfig) ExecInstall(osInfo tb.OSInfoExt, already []tb.TegnFeature, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	zshrcBackup := tb.TegnParameterToBool(params["zshrc-backup"])

	path := p.getInstallDir(osInfo)
	repo, err := git.PlainClone(path, &git.CloneOptions{
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
	
	userHomeDir, err := os.UserHomeDir()
	zshrcOrigPath := filepath.Join(userHomeDir, ".zshrc")
	if zshrcBackup && err != nil && platform.FileEntryExists(zshrcOrigPath) {
		err := platform.CopyFile(zshrcOrigPath, zshrcOrigPath + ".backup-mce")
		if err != nil {
			return fmt.Errorf("ExecInstall zshrc backup error: %w", err)
		}
	}

	templateFile := filepath.Join(path, "templates/zshrc.zsh-template")
	if !platform.FileEntryExists(templateFile) {
		return fmt.Errorf("ExecInstall .zshrc template file does not exist (%s)", templateFile)
	}
	err = platform.CopyFile(templateFile, zshrcOrigPath)
	if err != nil {
		return err
	}

	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

// func (p *ZshBaseConfig)  ExecUninstall() error {

// }
