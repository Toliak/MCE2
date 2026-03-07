package data

import (
	"fmt"
)

type Distrib struct {
	Id      string  `json:"id"`
	IdLike  []string  `json:"id_like"`
	Name    string  `json:"name"`
	Version Version `json:"version"`
	Unknown bool    `json:"unknown"`
}

// ---------- Functional Options ----------

func NewDistrib(id string, idLike []string, name string, version Version) Distrib {
	return Distrib{
		Id:      id,
		IdLike:  idLike,
		Name:    name,
		Version: version,
		Unknown: false,
	}
}

func NewUnknownDistrib() Distrib {
	return Distrib{
		Unknown: true,
	}
}


func (d *Distrib) IsUnknown() bool {
	return d.Unknown
}

func (d *Distrib) GetIdLikeOrId() []string {
	if len(d.IdLike) != 0 {
		return d.IdLike
	}
	
	return []string{d.Id}
}

// String returns a string representation of the distribution
func (d *Distrib) String() string {
	return fmt.Sprintf("%s %s", d.Name, d.Version.String())
}

// GoString returns a Go syntax representation of the distribution
func (d *Distrib) GoString() string {
	var versionGoStr string
	
	return fmt.Sprintf("&Distrib{id: %q, idLike: %q, name: %q, version: %s}", 
		d.Id, d.IdLike, d.Name, versionGoStr)
}
