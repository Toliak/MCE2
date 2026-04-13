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

func prepareReplaceConfigTheme(zshrcPath string) error {
	inputFile, err := os.Open(zshrcPath)
    if err != nil {
        return fmt.Errorf("prepareReplaceConfigTheme failed to open file: %w", err)
    }

	var lines []string
	scanner := bufio.NewScanner(inputFile)

	lines, times, err := sedparody.
		NewReplacer(
			sedparody.ScannerToReplacerReader(scanner),
		).Replace(
			zshThemeRegexpReplace,
			fmt.Sprintf("%s'%s'", "$1", filepath.Join(zshPowerLevel10kThemeName, "powerlevel10k")),
			1,
		)

	inputFile.Close()
	if err != nil {
		return fmt.Errorf("prepareReplaceConfigTheme error: %w", err)
	}
	if times == 0 {
		return fmt.Errorf("prepareReplaceConfigTheme: Not found line to replace")
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
		return fmt.Errorf("prepareReplaceConfigTheme Flush error: %w", err)
	}

    return nil
}

func (p *ZshPowerLevel10k) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	url := params["repo-url"]
	branch := params["repo-branch"]

	path := getInstallDirZshPowerLevel10k(osInfo)
	repo, err := git.PlainClone(
		path, 
		defaultGitCloneOptions(func (v *git.CloneOptions) {
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

	zshrcOrigPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("ExecInstall failed to get zshrc path: %w", err)
	}

	err = prepareReplaceConfigTheme(zshrcOrigPath)
	if err != nil {
		return fmt.Errorf("ExecInstall prepare config error: %w", err)
	}
	return nil
}
