package tegnsett

import (
	"encoding/json"
	"fmt"

	"github.com/toliak/mce/optional"
	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

type OSPackages struct {
	children []tb.Tegn
}

var _ tb.Tegnsett = (*OSPackages)(nil)

func NewOSPackages(children []tb.TegnBuildFunc) tb.TegnsettBuildFunc {
	return func() tb.Tegnsett {
		v := OSPackages{
			// children will be add later
		}
		v.children = make([]tb.Tegn, 0, len(children))
		for _, buildChild := range children {
			child := buildChild()
			if child == nil {
				fmt.Println("The OSPackages child build function returned nil, skipped")
				continue
			}

			if _, ok := child.(tegn.GenericPackageTegn); !ok {
				fmt.Printf("The child '%s' cannot be added as a child of the OSPackages. Because it does not inherit GenericPackageTegn\n", child.GetID())
				continue
			}

			v.children = append(v.children, child)
		}

		return &v
	}
}

// GetID implements [tb.Tegnsett].
func (p *OSPackages) GetID() string {
	return "os-packages"
}

// GetName implements [tb.Tegnsett].
func (p *OSPackages) GetName() string {
	return "os-packages"
}

// GetDescription implements [tb.Tegnsett].
func (p *OSPackages) GetDescription() string {
	return "os-packages"
}

// GetAvailableCPUArch implements [tb.Tegnsett].
func (p *OSPackages) GetAvailableCPUArch() optional.Optional[[]data.CPUArchE] {
	return optional.None[[]data.CPUArchE]()
}

// GetAvailableOsType implements [tb.Tegnsett].
func (p *OSPackages) GetAvailableOsType() optional.Optional[[]data.OSTypeE] {
	return optional.None[[]data.OSTypeE]()
}

// GetAvailability implements [tb.Tegnsett].
func (p *OSPackages) GetAvailability(_osInfo *tb.OSInfoExt, _before tb.TegnInstalledFeaturesMap, _enabledIds tb.TegnGeneralEnabledIDsMap) tb.TegnAvailability {
	// return tegnbuilder.TegnAvailability{
	// 	Available: p.info.PkgManager.V != data.PkgMgrUnknown,
	// 	Reason:    fmt.Sprintf("Package manager is unknown (%v)", p.info.PkgManager),
	// }
	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegnsett].
func (p *OSPackages) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetChildren implements [tb.Tegnsett].
func (p *OSPackages) GetChildren() []tb.Tegn {
	return p.children
}

// There is a "contract":
// we assume that the installedTegns are the children of the [OSPackages]
// therefore they all implement the [tegn.GenericPackageTegn] interface
func (p *OSPackages) ExecPostInstall(
	installedTegns []tb.Tegn,
	osInfo *tb.OSInfoExt,
	_already tb.TegnInstalledFeaturesMap,
	_tegnToParams map[string]tb.TegnParameterMap,
) error {
	if len(installedTegns) == 0 {
		// nothing to install
		return nil
	}

	packageTegns := make([]tegn.GenericPackageTegn, 0, len(installedTegns))
	for _, v := range installedTegns {
		packageTegn, ok := v.(tegn.GenericPackageTegn)
		if !ok {
			return fmt.Errorf("Tegn '%s' is not a GenericPackageTegn", v.GetID())
		}
		packageTegns = append(packageTegns, packageTegn)
	}

	pkgNames := make([]string, 0, len(packageTegns))
	for _, pkg := range packageTegns {
		pkgNames = append(pkgNames, pkg.GetPackageName())
	}

	return platform.InstallPackages(osInfo.PkgManager, pkgNames)
}

func (p *OSPackages) GoString(osInfo *tb.OSInfoExt) string {
	children := make([]map[string]any, len(p.children))
	for i, v := range p.children {
		children[i] = map[string]any{
			"id":     v.GetID(),
			"name":   v.GetName(),
			"params": v.GetParameters(osInfo),
		}
	}

	str, _ := json.Marshal(map[string]any{
		"id":       p.GetID(),
		"name":     p.GetName(),
		"children": children,
	})

	return string(str)
}
