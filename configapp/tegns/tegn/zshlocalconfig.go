package tegn

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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
		tb.NewTegnParameter(
			"enable-plugins",
			"Enable featured plugins",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getZshLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
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

		_, err = outputFile.WriteString("###### Config after the oh-my-zsh initialization.\n###### Managed by the MCE2\n\n")
		if err != nil {
			return fmt.Errorf("failed to write to config file %s: %w", localConfigPath, err)
		}

		if already["cfg:mce2"] {
			_, err = outputFile.WriteString(
				fmt.Sprintf("# <BEGIN> MCE2 config\nsource '%s'\n# <END> MCE2 config\n", filepath.Join(osInfo.MainInstallDir, "shell", "zsh.sh")),
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

		_, err = outputFile.WriteString("###### Config before the oh-my-zsh initialization.\n###### Managed by the MCE2\n\n")
		if err != nil {
			return fmt.Errorf("failed to write to config file %s: %w", localPreConfigPath, err)
		}

		if already["cfg:mce2"] {
			_, err = outputFile.WriteString(
				fmt.Sprintf("# <BEGIN> MCE2 pre-config\nsource '%s'\n# <END> MCE2 pre-config\n", filepath.Join(osInfo.MainInstallDir, "shell", "pre-zsh.sh")),
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
