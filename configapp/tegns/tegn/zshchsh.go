package tegn

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshChsh struct {}

var _ tb.Tegn = (*ZshChsh)(nil)

func NewTegnZshChshBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshChsh{}
	}
}

// GetID implements [tb.Tegn].
func (p *ZshChsh) GetID() string {
	return "zsh-chsh"
}

// GetName implements [tb.Tegn].
func (p *ZshChsh) GetName() string {
	return "Set Zsh as default shell"
}

// GetDescription implements [tb.Tegn].
func (p *ZshChsh) GetDescription() string {
	return `Set Zsh as the default user shell using chsh`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshChsh) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshChsh) GetAvailableOsType() *[]data.OSTypeE {
	return &[]data.OSTypeE{
		data.OSTypeLinux,
	}
}

// GetAvailability implements [tb.Tegn].
func (p *ZshChsh) GetAvailability(
	_osInfo tb.OSInfoExt, 
	_before tb.TegnInstalledFeaturesMap, 
	_enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshChsh) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetParameters implements [tb.Tegn].
func (p *ZshChsh) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter{
		tb.NewTegnParameter(
			"use-sudo",
			"Use sudo with chsh",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Use sudo with chsh?"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshChsh) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("os:zsh-chsh"): true,
	}
}

func (p *ZshChsh) IsInstalled(_osInfo tb.OSInfoExt) bool {
	shell := os.Getenv("SHELL")
	return strings.Contains(shell, "zsh")
}

func (p *ZshChsh) ExecInstall(_osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// get path to zshExecInstall
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		return fmt.Errorf("failed to find zsh: %w", err)
	}

	useSudo := tb.TegnParameterToBool(params["use-sudo"])
	userCurrent, err := user.Current()
	if err != nil && useSudo {
		fmt.Printf("Unable to get the user, sudo will not be used: %s\n", err)
	}

	// change shell
	if useSudo {
		_, err = platform.ExecCommand(
			platform.NewExecCommandWrapper(
				platform.WithThrowExitCodeError(true),
				platform.WithCaptureStdin(true),
				platform.WithNeedsRoot(true), // sudo
			),
			"chsh", "-s", zshPath, userCurrent.Username,
		)
	} else {
		_, err = platform.ExecCommand(
			platform.NewExecCommandWrapper(
				platform.WithThrowExitCodeError(true),
				platform.WithCaptureStdin(true),
			),
			"chsh", "-s", zshPath,
		)
	}
	if err != nil {
		return fmt.Errorf("failed to change shell: %w", err)
	}

	return nil
}

func (p *ZshChsh) ExecUninstall(osInfo tb.OSInfoExt) error {
	fmt.Printf("Tegn '%s' is not removable\n", p.GetID())
	return nil
}
