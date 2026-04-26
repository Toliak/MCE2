package tegn

import (
	// "fmt"

	// "github.com/toliak/mce/inspector"
	"fmt"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

type GenericPackageTegn interface {
	tb.Tegn

	// Gets the package name
	GetPackageName() string
}

// The type describes the package that can be installed via the OS package manager
type GenericPackage struct {
	idSuffix string
	name string
	description string
	pkgName string
	osTypes *[]data.OSTypeE
}

var _ tb.Tegn = (*GenericPackage)(nil)
var _ GenericPackageTegn = (*GenericPackage)(nil)

func NewGenericPackageBuilder(
	idSuffix string, 
	name string, 
	description string, 
	pkgName string, 
	osTypes *[]data.OSTypeE,
) tb.TegnBuildFunc {
	// defaultPackages := []string {
	// 	"git",
	// 	"zsh",
	// 	"vim",
	// 	"mc",
	// }

	// TODO: idSuffix from the name

	return func () tb.Tegn {
		return &GenericPackage {
			idSuffix: idSuffix,
			name: name,
			description: description,
			pkgName: pkgName,
			osTypes: osTypes,
		}
	}

	// packages := make([]string, 0, len(defaultPackages))
	// for _, v := range defaultPackages {
	// 	_, found := slices.BinarySearch(info.AvailableManagerPackages, v)
	// 	if !found {
	// 		fmt.Printf("Package '%s' not found in the OS package manager list %d\n", v, len(info.AvailableManagerPackages))
	// 		continue
	// 	}

	// 	packages = append(packages, v)
	// }

	// packageParams := make(map[string]bool, len(packages))
	// for _, v := range packages {
	// 	packageParams[v] = false
	// }

	// features []string

	// return &LinuxPackages{
	// 	info: info,
	// 	packageParams: packageParams,
	// }
}

func (p *GenericPackage) GetPackageName() string {
	return p.pkgName
}

// var _ tb.TegnBuildFunc = NewTegnLinuxPackages
// GetID implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetID() string {
	return "package-" + p.idSuffix
}

// GetName implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetName() string {
	return p.name
}

// GetDescription implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetDescription() string {
	return p.description
}

// GetAvailableCPUArch implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetAvailableOsType() *[]data.OSTypeE {
	return p.osTypes
}

// GetAvailability implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetAvailability(
	osInfo tb.OSInfoExt,
	_before tb.TegnInstalledFeaturesMap,
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	// TODO: make that more complex (request the checker from the user)
	return tb.TegnAvailability {
		Available: osInfo.AvailableManagerPackages[p.pkgName],
		Reason: fmt.Sprintf("Package '%s' not found in the package manager available packages list", p.pkgName),
	}
	// return tb.TegnAvailability{
	// 	Available: p.info.PkgManager.V != data.PkgMgrUnknown,
	// 	Reason:    fmt.Sprintf("Package manager is unknown (%v)", p.info.PkgManager),
	// }
}

// GetFeatures implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetFeatures() tb.TegnInstalledFeaturesMap {
	// result := make([]string, 0, len(p.packageParams))
	// for k, v := range p.packageParams {
	// 	if !v {
	// 		continue
	// 	}
	// 	result = append(result, fmt.Sprintf("app:%s", k))
	// }

	// return result
	return tb.TegnInstalledFeaturesMap{
		tb.TegnFeature("pkg:" + p.idSuffix): true,
	}
}

// GetBeforeIDs implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetBeforeIDs() []string {
	return make([]string, 0)
	// return []
}

// GetParameters implements [tegnbuilder.Tegn].
func (p *GenericPackage) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	// result := make([]tb.TegnParameter, 0, len(p.packageParams))
	// for k, v := range p.packageParams {
	// 	result = append(result, tb.NewTegnParameter(
	// 		k,
	// 		tb.TegnParameterTypeBool,
	// 		tb.WithDescription(k),
	// 		tb.WithDefaultValue(tb.TegnParameterFromBool(v)),
	// 		tb.WithAvailabilityTrue(),
	// 	))
	// }

	// slices.SortFunc(result, func(tp1, tp2 tb.TegnParameter) int {
	// 	return strings.Compare(tp1.Name, tp2.Name)
	// })
	// return result
	return make([]tb.TegnParameter, 0)
}

// func (p *LinuxPackages) SetParameter(name string, value string) error {
// 	_, ok := p.packageParams[name]
// 	if !ok {
// 		return nil
// 	}

// 	p.packageParams[name] = tb.TegnParameterToBool(value)
// 	return nil
// }


// GetFeatures implements [tegnbuilder.Tegn].
// func (p *LinuxPackages) SetContextFeatures(features []string) {}


func (p *GenericPackage) IsInstalled(osInfo tb.OSInfoExt) bool {
	// TODO: it is not true
	return platform.CommandExists(p.pkgName)
}

func (p *GenericPackage) ExecInstall(_osInfo tb.OSInfoExt, _already tb.TegnInstalledFeaturesMap, _params tb.TegnParameterMap) error {
	// Do not install packages here. They must be installed in the Tegnsett stage.
	return nil
}

func (p *GenericPackage) ExecUninstall(osInfo tb.OSInfoExt) error {
	fmt.Printf("Tegn '%s' (pkgName: %s) is not removable\n", p.GetID(), p.pkgName)
	return nil
}
