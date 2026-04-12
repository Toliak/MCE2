package tegn

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	"github.com/toliak/mce/sedparody"
	tb "github.com/toliak/mce/tegnbuilder"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type ZshBaseConfig struct {}

var _ tb.Tegn = (*ZshBaseConfig)(nil)

func NewTegnZshBaseConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshBaseConfig{}
	}
}

func getInstallDirZshBaseConfig(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "oh-my-zsh")
}

func getZshrcPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	zshrcPath := filepath.Join(userHomeDir, ".zshrc")
	return zshrcPath, err
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
}

// GetParameters implements [tb.Tegn].
func (p *ZshBaseConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	// TODO: cache that maybe?

	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"repo-url",
			"Repository URL",
			tb.TegnParameterTypeString,
			tb.WithDescription("Repository URL"),
			tb.WithDefaultValue("https://github.com/ohmyzsh/ohmyzsh"),
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
			tb.WithAvailabilityTrue(),
		),

		// Oh-my-zsh related configuration
		tb.NewTegnParameter(
			"zshrc-editor",
			"EDITOR",
			tb.TegnParameterTypeString,
			tb.WithDescription("EDITOR variable.\nLeave the string empty to leave the variable unchanged"),
			tb.WithDefaultValue("vim"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"zshrc-add-local-bin-path",
			"PATH += HOME/.local/bin",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Add \"$HOME/.local/bin\" to the PATH variable"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"zshrc-path-additional",
			"Additional PATHs",
			tb.TegnParameterTypeString,
			tb.WithDescription("Add paths into the PATH variable (separated by the colon).\nLeave the string empty to leave the variable unchanged"),
			tb.WithDefaultValue(""),
			tb.WithAvailabilityTrue(),
			tb.WithValidator(func(self *tb.TegnParameter, newValue string) error {
				re := regexp.MustCompile(`^[^"']+$`)
				if !re.MatchString(newValue) {
					return fmt.Errorf("The value '%s' did not match the regexp '%v'", newValue, re)
				}

				return nil
			}),
		),
		tb.NewTegnParameter(
			"zshrc-update",
			"Enable auto-update",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Enable auto-update (with the confirmation)"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(false)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"zshrc-case-sensitive",
			"Enable case-sensitive completion",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Set to true to force case-sensitive completion"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"zshrc-hist-stamps",
			"History stamp format",
			tb.TegnParameterTypeString,
			tb.WithDescription("Oh My Zsh provides a wrapper for the history command.\nYou can use this setting to decide whether to show a timestamp for each command in the history output.\nLeave the string empty to leave the variable unchanged"),
			tb.WithDefaultValue("yyyy-mm-dd"),
			tb.WithAvailabilityTrue(),
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
	return platform.FileEntryExists(path)
}

var zshExportRegexpReplace = regexp.MustCompile(`^(?:\s*#\s*)?((?:export\s+)?ZSH=).+$`)
var zshSourceOhMyZshLine = regexp.MustCompile(`^\s*source\s+.+/oh-my-zsh.sh\s*$`)

func prepareZshrcConfigToLinesWithReplace(zshrcPath string, ohMyZshDir string, textBeforePluginSource string) ([]string, error) {
	inputFile, err := os.Open(zshrcPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	lines, times, err := sedparody.
		NewReplacer(
			sedparody.ScannerToReplacerReader(scanner),
		).Replace(
			zshExportRegexpReplace,
			fmt.Sprintf("%s'%s'", "$1", ohMyZshDir),
			1,
		)
	if err != nil {
		return nil, fmt.Errorf("Replacer error: %w", err)
	}

	// TODO: looks like we can union the replacements, but for now I don't know how
	// Maybe 
	// 1. Scan all lines
	// 2. Traverse all of them and replace
	if textBeforePluginSource != "" {
		replaced := false
		for i, line := range lines {
			if zshSourceOhMyZshLine.MatchString(line) {
				lines = slices.Insert(lines, i, textBeforePluginSource)
				replaced = true
				break
			}
		}
		if !replaced {
			return nil, fmt.Errorf("Unable to find line by regexp '%v' in .zshrc", zshSourceOhMyZshLine)
		}
	}

	if times == 0 {
		return nil, fmt.Errorf("Not found line to replace")
	}

	return lines, nil
}

// TODO: encapsulate into the separated function with the [prepareReplaceConfigTheme]
func prepareReplaceConfig(zshrcPath string, ohMyZshDir string, textBeforePluginSource string) error {
	linesToWrite, err := prepareZshrcConfigToLinesWithReplace(zshrcPath, ohMyZshDir, textBeforePluginSource)
	if err != nil {
		return fmt.Errorf("prepareZshrcConfigToLinesWithReplace error: %w", err)
	}

	// Write back to same file (truncates)
    outputFile, err := os.Create(zshrcPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
    for i, line := range linesToWrite {
		// TODO: handle errors (?)
        writer.WriteString(line)
        if i < len(linesToWrite)-1 {
            writer.WriteString("\n")
        }
    }

	return nil
}

func prepareZshrcBeforeSourceLine(params tb.TegnParameterMap) string {
	var sb strings.Builder
	sb.Grow(512)

	if editor := params["zshrc-editor"]; editor != "" {
		sb.WriteString("export EDITOR=\"")
		sb.WriteString(editor)
		sb.WriteString("\"\n")
	}
	addLocalBin := tb.TegnParameterToBool(params["zshrc-add-local-bin-path"])
	additionsPath := params["zshrc-path-additional"]
	if addLocalBin || additionsPath != "" {
		sb.WriteString("export PATH=\"$PATH")
		if addLocalBin {
			sb.WriteString(":$HOME/.local/bin")
		}
		if additionsPath != "" {
			sb.WriteString(":")
			sb.WriteString(additionsPath)
		}
		sb.WriteString("\"\n")
	}
	if !tb.TegnParameterToBool(params["zshrc-update"]) {
		sb.WriteString(`
zstyle ':omz:update' mode disabled
zstyle ':omz:update' frequency 999999
UPDATE_ZSH_DAYS=999999
DISABLE_AUTO_UPDATE=true`)
		sb.WriteString("\n\n")
	}
	sb.WriteString("CASE_SENSITIVE=")
	if tb.TegnParameterToBool(params["zshrc-case-sensitive"]) {
		sb.WriteString("true\n")
	} else {
		sb.WriteString("false\n")
	}
	
	if stamp := params["zshrc-hist-stamps"]; stamp != "" {
		sb.WriteString("HIST_STAMPS=\"")
		sb.WriteString(stamp)
		sb.WriteString("\"\n")
	}
	if sb.Len() != 0 {
		sb.WriteString("\n")
	}

	return sb.String()
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
	
	zshrcOrigPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("ExecInstall failed to get zshrc path: %w", err)
	}

	if zshrcBackup && platform.FileEntryExists(zshrcOrigPath) {
		err := platform.CopyFile(zshrcOrigPath, zshrcOrigPath + ".backup-mce")
		if err != nil {
			return fmt.Errorf("ExecInstall zshrc backup error: %w", err)
		}
	}

	templateFile := filepath.Join(path, "templates", "zshrc.zsh-template")
	if !platform.FileEntryExists(templateFile) {
		return fmt.Errorf("ExecInstall .zshrc template file does not exist (%s)", templateFile)
	}
	err = platform.CopyFile(templateFile, zshrcOrigPath)
	if err != nil {
		return fmt.Errorf("ExecInstall platform.CopyFile error: %w", err)
	}

	textBeforeSource := prepareZshrcBeforeSourceLine(params)
	err = prepareReplaceConfig(zshrcOrigPath, path, textBeforeSource)

	if err != nil {
		return fmt.Errorf("ExecInstall prepare config error: %w", err)
	}
	return nil
}

// func (p *ZshBaseConfig)  ExecUpdate() error {

// }

// func (p *ZshBaseConfig)  ExecUninstall() error {

// }
