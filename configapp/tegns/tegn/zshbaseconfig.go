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

type ZshBaseConfig struct {
	installDir string
}

var _ tb.Tegn = (*ZshBaseConfig)(nil)

func NewTegnZshBaseConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshBaseConfig{}
	}
}

func getInstallDirZshBaseConfig(osInfo tb.OSInfoExt) string {
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
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh-my-zsh Repository URL"),
			tb.WithDefaultValue("https://github.com/ohmyzsh/ohmyzsh"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"repo-branch",
			"Repository branch",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh-my-zsh Repository branch"),
			tb.WithDefaultValue("master"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"install-dir",
			"Installation path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh-my-zsh Installation path (read-only)"),
			tb.WithDefaultValue(getInstallDirZshBaseConfig(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"zshrc-backup",
			"Do zshrc backup",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Backup current .zshrc configuration?"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailability(
				zshrcAvailabilityReason != "",
				zshrcAvailabilityReason,
			),
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

var zshExportRegexpReplace = regexp.MustCompile(`^(?:#\s*)?((?:export\s+)?ZSH=).+$`)

func prepareReplaceConfig(zshrcPath string, ohMyZshDir string) error {
	inputFile, err := os.Open(zshrcPath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }

	var lines []string
	scanner := bufio.NewScanner(inputFile)

	lines, times, err := sedparody.
		NewReplacer(
			sedparody.ScannerToReplacerReader(scanner),
		).Replace(
			zshExportRegexpReplace,
			fmt.Sprintf("%s'%s'", "$1", ohMyZshDir),
			1,
		)

	inputFile.Close()
	if err != nil {
		return fmt.Errorf("prepareReplaceConfig error: %w", err)
	}
	if times == 0 {
		return fmt.Errorf("Not found line to replace")
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
    return writer.Flush()
}

// TODO: shouldn't this thing be idempotent
func (p *ZshBaseConfig) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if debug build -> check the params
	url := params["repo-url"]
	branch := params["repo-branch"]
	zshrcBackup := tb.TegnParameterToBool(params["zshrc-backup"])

	path := getInstallDirZshBaseConfig(osInfo)
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

	err = prepareReplaceConfig(zshrcOrigPath, path)
	if err != nil {
		return fmt.Errorf("ExecInstall prepare config error: %w", err)
	}
	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

// func (p *ZshBaseConfig)  ExecUninstall() error {

// }
