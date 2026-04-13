package tegn

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type ZshAutoSuggestions struct {}

var _ tb.Tegn = (*ZshAutoSuggestions)(nil)

func NewTegnZshAutoSuggestionsBuilder() tb.TegnBuildFunc {
	return func() tb.Tegn {
		return &ZshAutoSuggestions{}
	}
}

// GetID implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetID() string {
	return "cfg-zsh-autosuggestions"
}

// GetName implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetName() string {
	return "cfg-zsh-autosuggestions"
}

// GetDescription implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetDescription() string {
	return `cfg-zsh-autosuggestions

With enhancements
`
}

// GetAvailableCPUArch implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	if err := tb.CheckFeatures(before, []tb.TegnFeature{"cfg:oh-my-zsh", "cfg:zsh-local"}); err != nil {
		return tb.NewTegnNotAvailable(err.Error())
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetBeforeIDs() []string {
	return []string{"cfg-local-zsh"}
}

// GetParameters implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	return []tb.TegnParameter {
		tb.NewTegnParameter(
			"additional-config",
			"Additional configuration",
			tb.TegnParameterTypeBool,
			tb.WithDescription("TODO:"),
			tb.WithDefaultValue(tb.TegnParameterFromBool(true)),
			tb.WithAvailabilityTrue(),
		),
	}
}

// GetFeatures implements [tb.Tegn].
func (p *ZshAutoSuggestions) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("cfg:zsh-autosuggestions"): true, 
	}
}

func (p *ZshAutoSuggestions) IsInstalled(osInfo tb.OSInfoExt) bool {
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)
	if !platform.FileEntryExists(localPreConfigPath) {
		return false
	}

	inputFile, err := os.Open(localPreConfigPath)
    if err != nil {
        return false
    }
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "<BEGIN> MCE2 zsh-autosuggestions") {
			return true
		}
	}

	return false
}

func (p *ZshAutoSuggestions) ExecInstall(osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, params tb.TegnParameterMap) error {
	// Assume we have `cfg:zsh-local` as a requirement
	// localConfigPath := getZshLocalConfigPath(osInfo)
	localPreConfigPath := getZshLocalPreOhMyZshConfigPath(osInfo)

	// TODO: git clone, because it is not available by default in oh-my-zsh

	err := platform.AppendFilepathString(
		localPreConfigPath,
		"\n# <BEGIN> MCE2 zsh-autosuggestions\nplugins+=(zsh-autosuggestions)\n# <END> MCE2 zsh-autosuggestions\n",
	)
	if err != nil {
		return fmt.Errorf("ExecInstall AppendFilepathString %s: %w", localPreConfigPath, err)
	}

	additionalConfig := tb.TegnParameterToBool(params["additional-config"])
	if additionalConfig {
		err := platform.AppendFilepathString(
			localPreConfigPath,
			"\n# <BEGIN> MCE2 zsh-autosuggestions cfg\nif [ \"$TERM\" = \"linux\" ]; then\n  ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE=\"fg=yellow\"\nfi\n# <END> MCE2 zsh-autosuggestions cfg\n",
		)
		if err != nil {
			return fmt.Errorf("ExecInstall AppendFilepathString %s: %w", localPreConfigPath, err)
		}
	}

	return nil
}
