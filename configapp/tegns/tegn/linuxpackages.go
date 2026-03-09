package tegn

import (
	// "fmt"

	// "github.com/toliak/mce/inspector"
	"fmt"
	"slices"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	tb "github.com/toliak/mce/tegnbuilder"
)

type LinuxPackages struct {
	info tb.TegnBuilderData
	
	packageParams map[string]bool
}

var _ tb.Tegn = (*LinuxPackages)(nil)

func NewTegnLinuxPackages(info tb.TegnBuilderData) tb.Tegn {
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

	// features []string

	return &LinuxPackages{
		info: info,
		packageParams: packageParams,
	}
}

var _ tb.TegnBuildFunc = NewTegnLinuxPackages

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
func (p *LinuxPackages) GetAvailability() tb.TegnAvailability {
	return tb.TegnAvailability{
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
func (p *LinuxPackages) GetParameters() []tb.TegnParameter {
	result := make([]tb.TegnParameter, 0, len(p.packageParams))
	for k, v := range p.packageParams {
		result = append(result, tb.NewTegnParameter(
			k,
			tb.TegnParameterTypeBool,
			tb.WithDescription(k),
			tb.WithValue(tb.TegnParameterFromBool(v)),
			tb.WithAvailabilityTrue(),
		))
	}

	slices.SortFunc(result, func(tp1, tp2 tb.TegnParameter) int {
		return strings.Compare(tp1.Name, tp2.Name)
	})
	return result
}

func (p *LinuxPackages) SetParameter(name string, value string) error {
	_, ok := p.packageParams[name]
	if !ok {
		return nil
	}

	p.packageParams[name] = tb.TegnParameterToBool(value)
	return nil
}


// GetFeatures implements [tegnbuilder.Tegn].
func (p *LinuxPackages) SetContextFeatures(features []string) {}
