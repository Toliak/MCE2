// Package tegn provides types and interfaces for defining installable
// configuration units.
// "Tegn" (Norwegian Bokmål for "character/symbol/sign")
// represents atomic functionality units to avoid collision with the "package" keyword.
package tegnbuilder

import (
	"fmt"
	"slices"

	"github.com/toliak/mce/osinfo/data"
)

// Utility struct to provide the availability information
type TegnAvailability struct {
	// Available indicates is the Tegn available.
	Available bool
	
	// Reason explains why the Tegn is unavailable (empty if Available is true).
	Reason string
}

// Returns the TegnAvailability with the Available=true
func NewTegnAvailable() TegnAvailability {
	return TegnAvailability {
		Available: true,
	}
}

// Returns the TegnAvailability with the Available=true
func NewTegnNotAvailable(reason string) TegnAvailability {
	return TegnAvailability {
		Available: false,
		Reason: reason,
	}
}

type TegnParameterOption func(*TegnParameter)

// Represents a parameter specification for a Tegn.
type TegnParameter struct {
	Name string
	Description string

	// Current parameter value.
	value string

	// Parameter's data type.
	ParamType TegnParameterType

	// Indicates the availability and contains the reason, if not available.
	Available TegnAvailability
}

func NewTegnParameter(name string, paramType TegnParameterType, opts ...TegnParameterOption) TegnParameter{
	p := TegnParameter{
		Name:        name,
		ParamType: paramType,
	}
	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// WithValue sets the private value field.
func WithValue(val string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.value = val
	}
}

// WithAvailability sets the availability status.
func WithAvailability(available TegnAvailability) TegnParameterOption {
	return func(p *TegnParameter) {
		p.Available = available
	}
}

// WithAvailability sets the availability status.
func WithAvailabilityTrue() TegnParameterOption {
	return func(p *TegnParameter) {
		p.Available = NewTegnAvailable()
	}
}

// WithAvailability sets the availability status.
func WithAvailabilityFalse(reason string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.Available = NewTegnNotAvailable(reason)
	}
}

// WithDescription overrides the description (useful if you want to set it optionally).
func WithDescription(desc string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.Description = desc
	}
}

// Data type of a TegnParameter.
type TegnParameterType int

const (
	TegnParameterTypeInt TegnParameterType = iota
	TegnParameterTypeBool
	TegnParameterTypeString
)

func (t TegnParameterType) String() string {
	switch t {
	case TegnParameterTypeInt:
		return "int"
	case TegnParameterTypeBool:
		return "bool"
	case TegnParameterTypeString:
		return "string"
	default:
		return "invalid"
	}
}

func TegnParameterFromBool(v bool) string {
	if v {
		return "y"
	}

	return "n"
}

func TegnParameterFromInt(v int) string {
	return fmt.Sprintf("%d", v)
}

func (p *TegnParameter) ToBool() bool {
	return TegnParameterToBool(p.value)
}

func TegnParameterToBool(s string) bool {
	return s == "y"
}

func (p *TegnParameter) ToInt(fallback int) int {
	return TegnParameterToInt(p.value, fallback)
}

func (p *TegnParameter) GetValue() string {
	return p.value
}

func TegnParameterToInt(s string, fallback int) int {
	var result *int
	_, err := fmt.Sscanf(s, "%d", result)
	if err != nil || result == nil {
		return fallback
	}

	return *result
}

// TegnGeneral defines the common interface for all Tegn types.
type TegnGeneral interface {
	// Returns the unique identifier.
	GetID() string

	// Returns the human-readable name.
	GetName() string

	// Returns a detailed description.
	GetDescription() string

	// Returns the OS types this Tegn supports.
	// Returns nil if the Tegn is available on all OS types.
	// This condition is evaluated before GetAvailableCPUArch.
	GetAvailableOsType() *[]data.OSTypeE

	// Returns the CPU architectures this Tegn supports.
	// Returns nil if the Tegn is available on all architectures.
	// This condition is evaluated before GetAvailability.
	GetAvailableCPUArch() *[]data.CPUArchE

	// Checks whether this Tegn is available for installation.
	//
	// In Tegnsett the method must not accumulate the children data!
	GetAvailability() TegnAvailability

	// Returns IDs of Tegns that must be installed before this one.
	// Only IDs from the same category take effect; IDs from other categories
	// are ignored.
	GetBeforeIDs() []string
}

func GetTegnGeneralAvailable(obj TegnGeneral, osInfo data.OSInfo) TegnAvailability {
	cpu := obj.GetAvailableCPUArch()
	if cpu != nil {
		if !slices.ContainsFunc(*cpu, func(v data.CPUArchE) bool {
			return v == osInfo.Arch.V
		}) {
			return NewTegnNotAvailable(
				fmt.Sprintf("Not available for the CPU Arch %v", osInfo.Arch),
			)
		}
	}

	osType := obj.GetAvailableOsType()
	if osType != nil {
		if !slices.ContainsFunc(*osType, func(v data.OSTypeE) bool {
			return v == osInfo.OsType.V
		}) {
			return NewTegnNotAvailable(
				fmt.Sprintf("Not available for the OS Type %v", osInfo.OsType),
			)
		}
	}

	return obj.GetAvailability()
}
