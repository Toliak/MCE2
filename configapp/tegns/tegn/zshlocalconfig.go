package tegn

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshLocalConfig struct {}

var _ tb.Tegn = (*ZshLocalConfig)(nil)

func NewTegnZshLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshLocalConfig{}
	}
}

func getZshLocalPreOhMyZshConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-pre-cfg.zsh")
}

func getZshLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-cfg.zsh")
}

// TODO: use
func getZshLocalPreConfigPlugins() []string {
	return []string{
		"git",
		"git-lfs",
		"docker",
		"docker-compose",
		"npm",
		"python",
		"tmux",
	}
}

// GetID implements [tb.Tegn].
func (p *ZshLocalConfig) GetID() string {
	return "cfg-local-zsh"
}

// GetName implements [tb.Tegn].
func (p *ZshLocalConfig) GetName() string {
	return "Local MCE2 config for Zsh"
}

// GetDescription implements [tb.Tegn].
func (p *ZshLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing plugins, aliases, etc add by the mce2 config app
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// fmt.Printf("before: %#v\n", before)
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:oh-my-zsh"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshLocalConfig) GetBeforeIDs() []string {
	return []string{"base-cfg-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getZshLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"pre-path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getZshLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		// TODO: implement
		tb.NewTegnParameter(
			"enable-plugins",
			"Enable featured plugins",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Enable plugins: :TODO:"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),

		// Oh-my-zsh related configuration
		// Assuming there is a feature-dependency `cfg:oh-my-zsh`
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
func (p *ZshLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-local"): true, 
	}
}

func (p *ZshLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getZshLocalConfigPath(osInfo)
	return platform.FileEntryExists(path)
}

func prepareZshrcLocalConfigLinesWithReplace(osInfo tb.OSInfoExt, zshrcPath string) ([]string, error) {
	localConfigPath := getZshLocalConfigPath(osInfo)
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)

	inputFile, err := os.Open(zshrcPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
	defer inputFile.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(inputFile)
	foundPreConfigLine := false
	for scanner.Scan() {
		line := scanner.Text()
		if !foundPreConfigLine && zshSourceOhMyZshLine.MatchString(line) {
			lines = append(
				lines,
				fmt.Sprintf("\n\n# <BEGIN> MCE2 local pre-config\nsource '%s'\n# <END> MCE2 local pre-config\n\n", localPreConfigPath),
			)
			foundPreConfigLine = true
		} 
		lines = append(lines, line)
	}
	if !foundPreConfigLine {
		return nil, fmt.Errorf("Unable to find line by regexp '%v' in .zshrc", zshSourceOhMyZshLine)
	}

	lines = append(
		lines, 
		fmt.Sprintf("\n\n# <BEGIN> MCE2 local config\nsource '%s'\n# <END> MCE2 local config\n\n", localConfigPath),
	)
	return lines, nil
}

func prepareZshPreLocalConfigText(params tb.TegnParameterMap) string {
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString("###### Config inserted before the oh-my-zsh initialization.\n###### Managed by MCE2\n###### DO NOT EDIT. This file may be overwritten or removed at any time\n\n")
	sb.WriteString("# <BEGIN> oh-my-zsh config options\n")

	if !tb.TegnParameterToBool(params["zshrc-update"]) {
		sb.WriteString("zstyle ':omz:update' mode disabled\nzstyle ':omz:update' frequency 999999\nUPDATE_ZSH_DAYS=999999\nDISABLE_AUTO_UPDATE=true\n\n")
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
	sb.WriteString("# <END> oh-my-zsh config options\n\n")

	return sb.String()
}

func (p *ZshLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// TODO: if _already contains mce2 configs, then use it

	localConfigPath := getZshLocalConfigPath(osInfo)
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)

	err := MkdirAllParent(localConfigPath)
	if err != nil {
		return fmt.Errorf("MkdirAll parent '%s' error: %w", localConfigPath, err)
	}

	{
		outputFile, err := os.Create(localConfigPath)
		if err != nil {
			return fmt.Errorf("failed to create config file %s: %w", localConfigPath, err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString("###### Config inserted after oh-my-zsh initialization.\n###### Managed by MCE2\n###### DO NOT EDIT. This file may be overwritten or removed at any time\n\n")
		if err != nil {
			return fmt.Errorf("failed to write to config file %s: %w", localConfigPath, err)
		}

		mce2zshConfig := filepath.Join(osInfo.MainInstallDir, "shell", "zsh.sh")
		if already["cfg:mce2"] && platform.FileEntryExists(mce2zshConfig) {
			_, err = fmt.Fprintf(
				outputFile,
				"# <BEGIN> MCE2 config\nsource '%s'\n# <END> MCE2 config\n", 
				mce2zshConfig,
			)
			if err != nil {
				return fmt.Errorf("failed to write to config file %s: %w", localConfigPath, err)
			}
		}
	}
	{
		outputFile, err := os.Create(localPreConfigPath)
		if err != nil {
			return fmt.Errorf("failed to create config file %s: %w", localPreConfigPath, err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString(prepareZshPreLocalConfigText(params))
		if err != nil {
			return fmt.Errorf("failed to write to config file %s: %w", localPreConfigPath, err)
		}

		mce2zshConfig := filepath.Join(osInfo.MainInstallDir, "shell", "pre-zsh.sh")
		if already["cfg:mce2"] && platform.FileEntryExists(mce2zshConfig) {
			_, err = fmt.Fprintf(
				outputFile,
				"# <BEGIN> MCE2 pre-config\nsource '%s'\n# <END> MCE2 pre-config\n", 
				mce2zshConfig,
			)
			if err != nil {
				return fmt.Errorf("failed to write to config file %s: %w", localPreConfigPath, err)
			}
		}
	}

	zshrcPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("failed to get zshrc path: %w", err)
	}

	lines, err := prepareZshrcLocalConfigLinesWithReplace(osInfo, zshrcPath)
	if err != nil {
		return fmt.Errorf("prepareZshrcLocalConfigLinesWithReplace error: %w", err)
	}

	outputFile, err := os.Create(zshrcPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
    for i, line := range lines {
		// TODO: handle errors (?)
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

func (p *ZshLocalConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove the entries from .zshrc
	zshrcPath, err := getZshrcPath()
	if err != nil {
		return fmt.Errorf("failed to get zshrc path: %w", err)
	}

	if platform.FileEntryExists(zshrcPath) {
		// Remove pre-config block
		err = removeConfigBlockFromFile(zshrcPath, "MCE2 local pre-config")
		if err != nil {
			return fmt.Errorf("removeBlockFromFile error '%s': %w", zshrcPath, err)
		}

		// Remove config block
		err = removeConfigBlockFromFile(zshrcPath, "MCE2 local config")
		if err != nil {
			return fmt.Errorf("removeBlockFromFile error '%s': %w", zshrcPath, err)
		}
	}

	// Remove the local config files
	localConfigPath := getZshLocalConfigPath(osInfo)
	if platform.FileEntryExists(localConfigPath) {
		err := os.Remove(localConfigPath)
		if err != nil {
			return fmt.Errorf("os.Remove error '%s': %w", localConfigPath, err)
		}
	}
	
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)
	if platform.FileEntryExists(localPreConfigPath) {
		err := os.Remove(localPreConfigPath)
		if err != nil {
			return fmt.Errorf("os.Remove error '%s': %w", localPreConfigPath, err)
		}
	}

	return nil
}