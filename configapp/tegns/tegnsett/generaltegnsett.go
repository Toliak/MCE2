package tegnsett

import (
	"encoding/json"
	"fmt"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
)

type GeneralTegnsett struct {
	info     tegnbuilder.TegnBuilderData

	id string
	name string
	description string
	beforeIDs []string

	children []tegnbuilder.Tegn
}

var _ tegnbuilder.Tegnsett = (*GeneralTegnsett)(nil)

func NewGeneralTegnsett(info tegnbuilder.TegnBuilderData) tegnbuilder.Tegnsett {
	return &GeneralTegnsett{
		info: info,
	}
}

func NewOuterGeneralTegnsett(
	id, name, description string,
	beforeIDs []string,
	children []tegnbuilder.TegnBuildFunc,
) tegnbuilder.TegnsettBuildFunc {
	return func(data tegnbuilder.TegnBuilderData) tegnbuilder.Tegnsett {
		v := NewGeneralTegnsett(data).(*GeneralTegnsett)
		v.id = id
		v.name = name
		v.description = description
		v.beforeIDs = beforeIDs

		v.children = make([]tegnbuilder.Tegn, len(children))
		for i, child := range children {
			v.children[i] = child(data)
		}

		return v
	}
}

var _ tegnbuilder.TegnsettBuildFunc = NewGeneralTegnsett
// var _ tegnbuilder.TegnsettOuterBuildFunc = NewOuterGeneralTegnsett

// GetID implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetID() string {
	return p.id
}

// GetName implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetName() string {
	return p.name
}

// GetDescription implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetDescription() string {
	return p.description
}

// GetAvailableCPUArch implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetAvailability() tegnbuilder.TegnAvailability {
	return tegnbuilder.TegnAvailability{
		Available: p.info.PkgManager.V != data.PkgMgrUnknown,
		Reason:    fmt.Sprintf("Package manager is unknown (%v)", p.info.PkgManager),
	}
}

// GetFeatures implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetFeatures() []string {
	// result := make([]string, 0)
	// for _, child := range p.children {
	// 	result = append(result, child.GetFeatures()...)
	// }

	return make([]string, 0)
}

// GetBeforeIDs implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetBeforeIDs() []string {
	return p.beforeIDs
}

// GetChildren implements [tegnbuilder.Tegnsett].
func (p *GeneralTegnsett) GetChildren() []tegnbuilder.Tegn {
	return p.children
}

func (p *GeneralTegnsett) GoString() string {
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
