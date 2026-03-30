package tegnsett

import (
	"encoding/json"

	"github.com/toliak/mce/osinfo/data"
	tb "github.com/toliak/mce/tegnbuilder"
)

type GenericTegnsett struct {
	id string
	name string
	description string
	beforeIDs []string

	children []tb.Tegn

	availabilityFunc tb.TegnAvailabilityFunc
}

var _ tb.Tegnsett = (*GenericTegnsett)(nil)

func NewGeneralTegnsett(
	id, name, description string,
	beforeIDs []string,
	children []tb.TegnBuildFunc,
	availabilityFunc tb.TegnAvailabilityFunc,
) tb.TegnsettBuildFunc {
	return func() tb.Tegnsett {
		v := GenericTegnsett{
			id: id,
			name: name,
			description: description,
			beforeIDs: beforeIDs,
			// children will be add later
			availabilityFunc: availabilityFunc,
		}

		v.children = make([]tb.Tegn, len(children))
		for i, child := range children {
			v.children[i] = child()
		}

		return &v
	}
}

// ----------------------

// GetID implements [tb.Tegnsett].
func (p *GenericTegnsett) GetID() string {
	return p.id
}

// GetName implements [tb.Tegnsett].
func (p *GenericTegnsett) GetName() string {
	return p.name
}

// GetDescription implements [tb.Tegnsett].
func (p *GenericTegnsett) GetDescription() string {
	return p.description
}

// GetAvailableCPUArch implements [tb.Tegnsett].
func (p *GenericTegnsett) GetAvailableCPUArch() *[]data.CPUArchE {
	// TODO: fill that in the constructor
	return nil
}

// GetAvailableOsType implements [tb.Tegnsett].
func (p *GenericTegnsett) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tb.Tegnsett].
func (p *GenericTegnsett) GetAvailability(
	osInfo tb.OSInfoExt, 
	before tb.TegnInstalledFeaturesMap, 
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	if p.availabilityFunc != nil {
		return p.availabilityFunc(osInfo, before, enabledIds)
	}

	return tb.NewTegnAvailable()
}

// GetBeforeIDs implements [tb.Tegnsett].
func (p *GenericTegnsett) GetBeforeIDs() []string {
	return p.beforeIDs
}

// GetChildren implements [tb.Tegnsett].
func (p *GenericTegnsett) GetChildren() []tb.Tegn {
	return p.children
}

func (p *GenericTegnsett) ExecPostInstall(
	installedTegns []tb.Tegn, 
	osInfo tb.OSInfoExt, 
	already tb.TegnInstalledFeaturesMap, 
	tegnToParams map[string]tb.TegnParameterMap,
) {
	// Do nothing here
}

func (p *GenericTegnsett) GoString(osInfo tb.OSInfoExt) string {
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
