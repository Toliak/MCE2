package data

import (
	"fmt"
)

type Distrib struct {
	id string
	idLike string
	name string
	version Version
	unknown bool
}

// ---------- Functional Options ----------

func NewDistrib(id, idLike, name string, version Version) Distrib {
	return Distrib{
		id:      id,
		idLike:  idLike,
		name:    name,
		version: version,
		unknown: false,
	}
}

func NewUnknownDistrib() Distrib {
	return Distrib{
		unknown: true,
	}
}

// ID returns the distribution ID
func (d *Distrib) ID() string {
	return d.id
}

// IDLike returns the distribution ID_LIKE
func (d *Distrib) IDLike() string {
	return d.idLike
}

// Name returns the distribution name
func (d *Distrib) Name() string {
	return d.name
}

// Version returns the distribution version
func (d *Distrib) Version() *Version {
	return &d.version
}

func (d *Distrib) IsUnknown() bool {
	return d.unknown
}

// String returns a string representation of the distribution
func (d *Distrib) String() string {
	return fmt.Sprintf("%s %s", d.name, d.version.String())
}

// GoString returns a Go syntax representation of the distribution
func (d *Distrib) GoString() string {
	var versionGoStr string
	
	return fmt.Sprintf("&Distrib{id: %q, idLike: %q, name: %q, version: %s}", 
		d.id, d.idLike, d.name, versionGoStr)
}