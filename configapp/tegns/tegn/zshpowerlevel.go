package tegn

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	"github.com/toliak/mce/sedparody"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

// oh-my-zsh -- tab-completion
// powerlevel 10k -- about prompt

type ZshPowerLevel10k struct {
	// TODO:
}

var _ tb.Tegn = (*ZshPowerLevel10k)(nil)

func NewTegnZshPowerLevel10kBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshPowerLevel10k{}
	}
}

var zshPowerLevel10kThemeName string = "powerlevel10k"

func getInstallDirZshPowerLevel10k(osInfo tb.OSInfoExt) string {
	return filepath.Join(getInstallDirZshBaseConfig(osInfo), "custom", "themes", zshPowerLevel10kThemeName)
}

func getZshLocalP10kPromptPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-p10k-prompt.zsh")
}

// GetID implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetID() string {
	return "cfg-zsh-p10k"
}

// GetName implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetName() string {
	return "PowerLevel10k"
}

// GetDescription implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetDescription() string {
	return `PowerLevel10k

URL: https://github.com/romkatv/powerlevel10k`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:oh-my-zsh", "cfg:mce2"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetBeforeIDs() []string {
	return []string{"base-cfg-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/romkatv/powerlevel10k"),
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
			tb.WithDefaultValue(getInstallDirZshPowerLevel10k(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"additional-config",
			"Auto config p10k",
			tb.TegnParameterTypeString,
			tb.WithDescription("Configure p10k prompt automatically (skip `p10k configure`)"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshPowerLevel10k) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-p10k"): true, 
	}
}

func (p *ZshPowerLevel10k) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getInstallDirZshPowerLevel10k(osInfo)
	return platform.FileEntryExists(path)
}

var zshThemeRegexpReplace = regexp.MustCompile(`^(?:\s*#\s*)?((?:export\s+)?ZSH_THEME=).+$`)

func prepareReplaceConfigThemeText(zshrcPath string, newThemeName string) ([]string, error) {
	inputFile, err := os.Open(zshrcPath)
    if err != nil {
        return nil, fmt.Errorf("prepareReplaceConfigTheme failed to open file: %w", err)
    }
	defer inputFile.Close()

	var lines []string
	scanner := bufio.NewScanner(inputFile)

	lines, times, err := sedparody.
		NewReplacer(
			sedparody.ScannerToReplacerReader(scanner),
		).Replace(
			zshThemeRegexpReplace,
			fmt.Sprintf("%s'%s'", "$1", newThemeName),
			1,
		)

	if err != nil {
		return nil, fmt.Errorf("prepareReplaceConfigTheme error: %w", err)
	}
	if times == 0 {
		return nil, fmt.Errorf("prepareReplaceConfigTheme: Not found line to replace")
	}

	return lines, nil
}

func prepareReplaceConfigTheme(zshrcPath string, newThemeName string) error {
	lines, err := prepareReplaceConfigThemeText(zshrcPath, newThemeName)
	if err != nil {
		return err
	}

	// Write back to same file (truncates)
    outputFile, err := os.Create(zshrcPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
    for i, line := range lines {
        writer.WriteString(line)
        if i < len(lines)-1 {
            writer.WriteString("\n")
        }
    }

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Flush error: %w", err)
	}

    return nil
}

var zshrcP10kInstantPromptBlock string = `
# <BEGIN> p10k instant prompt
# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
	source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi
# <END> p10k instant prompt
`
var zshrcP10kInstantPromptBlockBytes = []byte(zshrcP10kInstantPromptBlock)

func (p *ZshPowerLevel10k) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]
	autoConf := tb.TegnParameterToBool(params["additional-config"])

	path := getInstallDirZshPowerLevel10k(osInfo)
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

	zshrcOrigPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("failed to get zshrc path: %w", err)
	}

	err = prepareReplaceConfigTheme(
		zshrcOrigPath,
		filepath.Join(zshPowerLevel10kThemeName, "powerlevel10k"),
	)
	if err != nil {
		return fmt.Errorf("prepare config error: %w", err)
	}

	if autoConf {
		zshrcBytes, err := os.ReadFile(zshrcOrigPath)
		zshrcText := string(zshrcBytes)
		if err != nil {
			return fmt.Errorf("os.ReadFile error '%s': %w", zshrcOrigPath, err)
		}
		
		purePromptPath := filepath.Join(getInstallDirZshPowerLevel10k(osInfo), "config", "p10k-lean.zsh")
		purePromptCopyPath := getZshLocalP10kPromptPath(osInfo)
		promptConfig := ""
		if !platform.FileEntryExists(purePromptPath) {
			fmt.Printf("Unable to find pure p10k prompt config: %s\n", purePromptPath)
		} else {
			err = platform.CopyFile(purePromptPath, purePromptCopyPath)
			if err != nil {
				fmt.Printf("platform.CopyFile '%s' -> '%s' error: %s\n", purePromptPath, purePromptCopyPath, err)
			} else {
				promptConfig = fmt.Sprintf("# <BEGIN> p10k prompt config\nsource '%s'\n# <END> p10k prompt config", purePromptCopyPath)
			}
		}

		err = os.WriteFile(
			zshrcOrigPath,
			fmt.Append(zshrcP10kInstantPromptBlockBytes, zshrcText, promptConfig),
			0644,
		)
		if err != nil {
			return fmt.Errorf("os.WriteFile error '%s': %w", zshrcOrigPath, err)
		}
	}

	return nil
}

func (p *ZshPowerLevel10k) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Restore the original ZSH_THEME setting in .zshrc
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("getZshrcPath error: %w", err)
	}

	if platform.FileEntryExists(zshrcPath) {
		// Remove p10k instant prompt block
		err = removeConfigBlockFromFile(zshrcPath, "p10k instant prompt")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error: %w", err)
		}

		// Remove p10k prompt config block
		err = removeConfigBlockFromFile(zshrcPath, "p10k prompt config")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error: %w", err)
		}

		// Reset the zsh theme
		err = prepareReplaceConfigTheme(
			zshrcPath,
			"robbyrussell",
		)
		if err != nil {
			return fmt.Errorf("prepare config error: %w", err)
		}
	}

	// Remove the p10k prompt config file if it was created
	purePromptCopyPath := getZshLocalP10kPromptPath(osInfo)
	if platform.FileEntryExists(purePromptCopyPath) {
		err := os.Remove(purePromptCopyPath)
		if err != nil {
			return fmt.Errorf("os.Remove error '%s': %w", purePromptCopyPath, err)
		}
	}

	// Remove the cloned theme
	installPath := getInstallDirZshPowerLevel10k(osInfo)
	if platform.FileEntryExists(installPath) {
		err := os.RemoveAll(installPath)
		if err != nil {
			return fmt.Errorf("os.RemoveAll error '%s': %w", installPath, err)
		}
	}

	return nil
}