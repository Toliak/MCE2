package tegn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
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

// Config path starting with tmux version 3.1
func getTmuxConfigPath(osInfo tb.OSInfoExt) string {
	dir := filepath.Join(osInfo.HomeDir, ".config", "tmux")
	return filepath.Join(dir, "tmux.conf")
}

// Config path with tmux version below 3.1
func getTmuxConfigPathOld(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.HomeDir, ".tmux.conf")
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
	gitParams := GetGitCloneParameters("https://github.com/gpakosz/.tmux", getInstallDirOhMyTmux(osInfo))
	params := []tb.TegnParameter {
		tb.NewTegnParameter(
			"cfg-compat",
			"Config compatibility",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Config compatibility with tmux version <3.1"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
	params = append(gitParams, params...)
	return params
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
	return platform.FileExists(path)
}

func (p *OhMyTmux) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	shallowClone := tb.TegnParameterToBool(params["use-shallow-clone"])
	path := getInstallDirOhMyTmux(osInfo)

	targetTmuxConf := getTmuxConfigPath(osInfo)
	err := platform.MkdirAllParent(targetTmuxConf)
	if err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", targetTmuxConf, err)
	}

	cloner := platform.NewGitClone(url, path, platform.WithGitBranch(branch), platform.WithGitShallowMode(shallowClone))
	err = cloner.Clone()
	if err != nil {
		return fmt.Errorf("GitClone error: %w", err)
	}

	// Create "reference" to .tmux.conf from the config
	sourceTmuxConf := filepath.Join(path, ".tmux.conf")

	err = platform.MkdirAllParent(targetTmuxConf)
	if err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", targetTmuxConf, err)
	}

	// TMUX_CONF and TMUX_CONF_LOCAL are required by oh-my-tmux
	err = platform.AppendFilepathString(
		targetTmuxConf,
		fmt.Sprintf(
			`# <BEGIN> Oh-my-tmux config (autogen mce2)
set-environment -g TMUX_CONF "%s"
set-environment -gu TMUX_CONF_LOCAL
source-file "%s"
# <END> Oh-my-tmux config (autogen mce2)` + "\n\n", 
			sourceTmuxConf,
			sourceTmuxConf,
		),
	)

	// Insert the local config entry into the ~/.tmux.conf (backwards compatibilty)
	cfgCompat := tb.TegnParameterToBool(params["cfg-compat"])
	if cfgCompat {
		tmuxOldConfPath := getTmuxConfigPathOld(osInfo)
		if !platform.IsFileSymlinkTo(tmuxOldConfPath, targetTmuxConf) {
			err = platform.AppendFilepathString(
				tmuxOldConfPath,
				fmt.Sprintf("\n\n# <BEGIN> MCE2 local config ref\nsource '%s'\n# <END> MCE2 local config ref\n\n", targetTmuxConf),
			)
		} else {
			fmt.Printf(
				"Config compatibility skipped, because '%s' is a symlink to a new config file '%s'\n",
				tmuxOldConfPath, targetTmuxConf,
			)
		}
	}

	// WARN: We are not backing up the files
	// We are not using .local file here! Just put your configs into the tmux.conf

	// Copy .tmux.conf.local template to the config directory
	// sourceLocalConf := filepath.Join(installPath, ".tmux.conf.local")
	// targetLocalConf := filepath.Join(tmuxConfigDir, "tmux.conf.local")

	// if tmuxConfBackup && platform.FileExists(targetLocalConf) {
	// 	err := platform.CopyFile(targetLocalConf, targetLocalConf + ".backup-mce")
	// 	if err != nil {
	// 		return fmt.Errorf("targetLocalConf backup error: %w", err)
	// 	}
	// }

	// Only copy if the target doesn't already exist, preserving user customizations
	// err = platform.CopyFile(sourceLocalConf, targetLocalConf)
	// if err != nil {
	// 	return fmt.Errorf("platform.CopyFile for tmux.conf.local error: %w", err)
	// }

	return nil
}

func (p *OhMyTmux) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the configuration block from tmux.conf
	tmuxConfPath := getTmuxConfigPath(osInfo)

	if platform.FileExists(tmuxConfPath) {
		err := removeConfigBlockFromFile(tmuxConfPath, "Oh-my-tmux config (autogen mce2)")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", tmuxConfPath, err)
		}
	}

	// Remove the cloned repository
	installPath := getInstallDirOhMyTmux(osInfo)
	if platform.FileExists(installPath) {
		err := os.RemoveAll(installPath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
		}
	}

	return nil
}