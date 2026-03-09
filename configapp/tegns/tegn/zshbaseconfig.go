package tegn

import (
	"fmt"
	"slices"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
)

type ZshBaseConfig struct {
	info tegnbuilder.TegnBuilderData
	features []string
	
	// packageParams map[string]bool
}

var _ tegnbuilder.Tegn = (*ZshBaseConfig)(nil)

func NewTegnZshBaseConfig(info tegnbuilder.TegnBuilderData) tegnbuilder.Tegn {
	return &ZshBaseConfig{
		info: info,
	}
}

var _ tegnbuilder.TegnBuildFunc = NewTegnZshBaseConfig

// GetID implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetID() string {
	return "zsh-base-config"
}

// GetName implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetName() string {
	return "zsh-base-config"
}

// GetDescription implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetDescription() string {
	return "zsh-base-config"
}

// GetAvailableCPUArch implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetAvailableOsType() *[]data.OSTypeE {
	return &[]data.OSTypeE{
		data.OSTypeLinux,
	}
}

// GetAvailability implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetAvailability() tegnbuilder.TegnAvailability {
	return tegnbuilder.TegnAvailability{
		Available: slices.Contains(p.features, "app:zsh") /*|| platform.CommandExists("zsh")*/,
		Reason:    fmt.Sprintf("Feature app:zsh not enabled"),
	}
}

// GetFeatures implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetFeatures() []string {
	return []string{"config:zsh-base"}
}

// GetBeforeIDs implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) GetParameters() []tegnbuilder.TegnParameter {
	return make([]tegnbuilder.TegnParameter, 0)
}

func (p *ZshBaseConfig) SetParameter(name string, value string) error {
	return nil
}

// GetFeatures implements [tegnbuilder.Tegn].
func (p *ZshBaseConfig) SetContextFeatures(features []string) {
	p.features = features
}
