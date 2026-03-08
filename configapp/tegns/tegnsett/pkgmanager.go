package tegnsett

import (
	"encoding/json"
	"fmt"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
)

type OSPackages struct {
	info     tegnbuilder.TegnBuilderData
	children []tegnbuilder.Tegn
}

var _ tegnbuilder.Tegnsett = (*OSPackages)(nil)

func NewOSPackages(info tegnbuilder.TegnBuilderData) tegnbuilder.Tegnsett {
	return &OSPackages{
		info: info,
	}
}

func NewOuterOSPackages(children []tegnbuilder.TegnBuildFunc) tegnbuilder.TegnsettBuildFunc {
	return func(data tegnbuilder.TegnBuilderData) tegnbuilder.Tegnsett {
		v := NewOSPackages(data).(*OSPackages)
		v.children = make([]tegnbuilder.Tegn, len(children))
		for i, child := range children {
			v.children[i] = child(data)
		}

		return v
	}
}

var _ tegnbuilder.TegnsettBuildFunc = NewOSPackages
var _ tegnbuilder.TegnsettOuterBuildFunc = NewOuterOSPackages

// GetID implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetID() string {
	return "os-packages"
}

// GetName implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetName() string {
	return "os-packages"
}

// GetDescription implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetDescription() string {
	return "os-packages"
}

// GetAvailableCPUArch implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetAvailability(features []string) tegnbuilder.TegnAvailability {
	return tegnbuilder.TegnAvailability{
		Available: p.info.PkgManager.V != data.PkgMgrUnknown,
		Reason:    fmt.Sprintf("Package manager is unknown (%v)", p.info.PkgManager),
	}
}

// GetFeatures implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetFeatures() []string {
	result := make([]string, 0)
	for _, child := range p.children {
		result = append(result, child.GetFeatures()...)
	}

	return result
}

// GetBeforeIDs implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetChildren implements [tegnbuilder.Tegnsett].
func (p *OSPackages) GetChildren() []tegnbuilder.Tegn {
	return p.children
}

func (p *OSPackages) GoString() string {
	children := make([]map[string]any, len(p.children))
	for i, v := range p.children {
		children[i] = map[string]any{
			"id":     v.GetID(),
			"name":   v.GetName(),
			"params": v.GetParameters(),
		}
	}

	str, _ := json.Marshal(map[string]any{
		"id":       p.GetID(),
		"name":     p.GetName(),
		"features": p.GetFeatures(),
		"children": children,
	})

	return string(str)
}
