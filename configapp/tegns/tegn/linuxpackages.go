package tegn

import (
	// "fmt"

	// "github.com/toliak/mce/inspector"
	"fmt"
	"slices"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
)

type LinuxPackages struct {
	info tegnbuilder.TegnBuilderData
	
	packageParams map[string]bool
}

var _ tegnbuilder.Tegn = (*LinuxPackages)(nil)

func NewPkgLinuxPackages(info tegnbuilder.TegnBuilderData) tegnbuilder.Tegn {
	defaultPackages := []string {
		"git",
		"zsh",
		"vim",
		// "mc",
	}

	packages := make([]string, 0, len(defaultPackages))
	for _, v := range defaultPackages {
		_, found := slices.BinarySearch(info.AvailableManagerPackages, v)
		if !found {
			fmt.Printf("Package '%s' not found in the OS package manager list %d\n", v, len(info.AvailableManagerPackages))
			continue
		}

		packages = append(packages, v)
	}

	packageParams := make(map[string]bool, len(packages))
	for _, v := range packages {
		packageParams[v] = false
	}

	return &LinuxPackages{
		info: info,
		packageParams: packageParams,
	}
}

var _ tegnbuilder.TegnBuildFunc = NewPkgLinuxPackages

// GetID implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetID() string {
	return "packages-linux"
}

// GetName implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetName() string {
	return "packages-linux"
}

// GetDescription implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetDescription() string {
	return "packages-linux"
}

// GetAvailableCPUArch implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetAvailableOsType() *[]data.OSTypeE {
	return &[]data.OSTypeE{
		data.OSTypeLinux,
	}
}

// GetAvailability implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetAvailability(features []string) tegnbuilder.TegnAvailability {
	return tegnbuilder.TegnAvailability{
		Available: p.info.PkgManager.V != data.PkgMgrUnknown,
		Reason:    fmt.Sprintf("Package manager is unknown (%v)", p.info.PkgManager),
	}
}

// GetFeatures implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetFeatures() []string {
	result := make([]string, 0, len(p.packageParams))
	for k, v := range p.packageParams {
		if !v {
			continue
		}
		result = append(result, fmt.Sprintf("app:%s", k))
	}

	return result
}

// GetBeforeIDs implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tegnbuilder.Tegn].
func (p *LinuxPackages) GetParameters() []tegnbuilder.TegnParameter {
	result := make([]tegnbuilder.TegnParameter, 0, len(p.packageParams))
	for k, v := range p.packageParams {
		result = append(result, tegnbuilder.TegnParameter {
			Name: k,
			Description: k,
			Value: tegnbuilder.TegnParameterFromBool(v),
			ParamType: tegnbuilder.TegnParameterTypeBool,
			Available: tegnbuilder.NewTegnAvailable(),
		})
	}

	slices.SortFunc(result, func(tp1, tp2 tegnbuilder.TegnParameter) int {
		return strings.Compare(tp1.Name, tp2.Name)
	})
	return result
}

func (p *LinuxPackages) SetParameter(name string, value string) error {
	_, ok := p.packageParams[name]
	if !ok {
		return nil
	}

	p.packageParams[name] = tegnbuilder.TegnParameterToBool(value)
	return nil
}
