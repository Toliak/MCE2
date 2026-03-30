package tegnsett

import (
	"encoding/json"
	"fmt"

	"github.com/toliak/mce/osinfo/data"
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
		v.children = make([]tb.Tegn, len(children))
		for i, child := range children {
			child := child()
			
			if _, ok := child.(tegn.GenericPackageTegn); !ok {
				// TODO: log error
				fmt.Printf("The child '%s' cannot be add as a child of the OSPackages. Because it does not inherit GenericPackageTegn\n", child.GetID())
				continue
			}

			v.children[i] = child
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
func (p *OSPackages) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tb.Tegnsett].
func (p *OSPackages) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegnsett].
func (p *OSPackages) GetAvailability(_osInfo tb.OSInfoExt, _before tb.TegnInstalledFeaturesMap, _enabledIds tb.TegnGeneralEnabledIDsMap) tb.TegnAvailability {
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
	_installedTegns []tb.Tegn, 
	_osInfo tb.OSInfoExt, 
	_already tb.TegnInstalledFeaturesMap, 
	_tegnToParams map[string]tb.TegnParameterMap,
) {
	// TODO: implement
	panic("NOT IMPLEMENTED ExecPostInstall for OSPackages")
}

func (p *OSPackages) GoString(osInfo tb.OSInfoExt) string {
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
