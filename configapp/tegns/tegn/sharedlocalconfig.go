package tegn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type SharedLocalConfig struct {}

var _ tb.Tegn = (*SharedLocalConfig)(nil)

func NewTegnSharedLocalConfigBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &SharedLocalConfig{}
	}
}

func getSharedLocalConfigPath(osInfo tb.OSInfoExt) string {
	return filepath.Join(osInfo.GetFullDataDir(), "local-cfg.shared")
}

// GetID implements [tb.Tegn].
func (p *SharedLocalConfig) GetID() string {
	return "cfg-local-shared"
}

// GetName implements [tb.Tegn].
func (p *SharedLocalConfig) GetName() string {
	return "Local MCE2 config for shared environment variables and functions"
}

// GetDescription implements [tb.Tegn].
func (p *SharedLocalConfig) GetDescription() string {
	return `Local MCE2 config file for managing shared environment variables, functions, and settings for multiple shells
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *SharedLocalConfig) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *SharedLocalConfig) GetBeforeIDs() []string {
	// return []string{"base-cfg-zsh", "cfg-local-zsh", "cfg-local-bash"}
	return []string{}
}

var sharedConfigPATHValidator *regexp.Regexp = regexp.MustCompile(`^[^"']+$`)

// GetParameters implements [tb.Tegn].
func (p *SharedLocalConfig) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"path",
			"Local config path",
			tb.TegnParameterTypeString,
			tb.WithDescription("Local config path (read-only)"),
			tb.WithDefaultValue(getSharedLocalConfigPath(osInfo)),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),

		// Shared variables
		tb.NewTegnParameter(
			"editor",
			"EDITOR",
			tb.TegnParameterTypeString,
			tb.WithDescription("EDITOR variable.\nLeave the string empty to leave the variable unchanged"),
			tb.WithDefaultValue("vim"),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"add-local-bin-path",
			"PATH += HOME/.local/bin",
			tb.TegnParameterTypeBool,
			tb.WithDescription("Add \"$HOME/.local/bin\" to the PATH variable"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
		tb.NewTegnParameter(
			"path-additional",
			"Additional PATHs",
			tb.TegnParameterTypeString,
			tb.WithDescription("Add paths into the PATH variable (separated by the colon).\nLeave the string empty to leave the variable unchanged"),
			tb.WithDefaultValue(""),
			tb.WithAvailabilityTrue(),
			tb.WithValidator(func(self *tb.TegnParameter, newValue string) error {
				if !sharedConfigPATHValidator.MatchString(newValue) {
					return fmt.Errorf("The value '%s' did not match the regexp '%v'", newValue, sharedConfigPATHValidator)
				}

				return nil
			}),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *SharedLocalConfig) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:shared-local"): true, 
	}
}

func (p *SharedLocalConfig) IsInstalled(osInfo tb.OSInfoExt) bool {
	path := getSharedLocalConfigPath(osInfo)
	return platform.FileEntryExists(path)
}

func prepareSharedLocalConfigText(params tb.TegnParameterMap) string {
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString("###### Shared Local MCE2 config file.\n###### Managed by MCE2\n###### DO NOT EDIT. This file may be overwritten or removed at any time\n\n")
	sb.WriteString("# <BEGIN> Shared config options\n")

	if editor := params["editor"]; editor != "" {
		sb.WriteString("export EDITOR=\"")
		sb.WriteString(editor)
		sb.WriteString("\"\n")
	}
	addLocalBin := tb.TegnParameterToBool(params["add-local-bin-path"])
	additionsPath := params["path-additional"]
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

	sb.WriteString("# <END> Shared config options\n\n")

	return sb.String()
}

func (p *SharedLocalConfig) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	localConfigPath := getSharedLocalConfigPath(osInfo)
	err := MkdirAllParent(localConfigPath)
	if err != nil {
		return fmt.Errorf(" MkdirAll parent '%s' error: %w", localConfigPath, err)
	}

	{
		outputFile, err := os.Create(localConfigPath)
		if err != nil {
			return fmt.Errorf(" failed to create config file %s: %w", localConfigPath, err)
		}
		defer outputFile.Close()
		_, err = outputFile.WriteString(prepareSharedLocalConfigText(params))
		if err != nil {
			return fmt.Errorf(" failed to write to config file %s: %w", localConfigPath, err)
		}
	}

	if already["cfg:bash-local"] {
		bashLocalConfigPath := getBashLocalConfigPath(osInfo)
		platform.AppendFilepathString(
			bashLocalConfigPath,
			fmt.Sprintf("\n# <BEGIN> MCE2 shared config\nsource '%s'\n# <END> MCE2 shared config\n\n", localConfigPath),
		)
	}
	if already["cfg:zsh-local"] {
		bashLocalConfigPath := getZshLocalConfigPath(osInfo)
		platform.AppendFilepathString(
			bashLocalConfigPath,
			fmt.Sprintf("\n# <BEGIN> MCE2 shared config\nsource '%s'\n# <END> MCE2 shared config\n\n", localConfigPath),
		)
	}

	return nil
}

func (p *SharedLocalConfig) ExecUninstall(osInfo tb.OSInfoExt) error {
	// Remove references from bash local config if it exists
	bashLocalConfigPath := getBashLocalConfigPath(osInfo)
	if platform.FileEntryExists(bashLocalConfigPath) {
		err := removeConfigBlockFromFile(bashLocalConfigPath, "MCE2 shared config")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", bashLocalConfigPath, err)
		}
	}

	// Remove references from zsh local config if it exists
	zshLocalConfigPath := getZshLocalConfigPath(osInfo)
	if platform.FileEntryExists(zshLocalConfigPath) {
		err := removeConfigBlockFromFile(zshLocalConfigPath, "MCE2 shared config")
		if err != nil {
			return fmt.Errorf("removeConfigBlockFromFile error '%s': %w", zshLocalConfigPath, err)
		}
	}

	// Remove the shared config file
	localConfigPath := getSharedLocalConfigPath(osInfo)
	if platform.FileEntryExists(localConfigPath) {
		err := os.Remove(localConfigPath)
		if err != nil {
			return fmt.Errorf("failed to remove shared config file '%s': %w", localConfigPath, err)
		}
	}

	return nil
}
