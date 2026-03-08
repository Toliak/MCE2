// Package tegn provides types and interfaces for defining installable
// configuration units.
// "Tegn" (Norwegian Bokmål for "character/symbol/sign")
// represents atomic functionality units to avoid collision with the "package" keyword.
package tegnbuilder

import (
	"fmt"

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

// Represents a parameter specification for a Tegn.
type TegnParameter struct {
	Name string
	Description string

	// Current parameter value.
	Value string

	// Parameter's data type.
	ParamType TegnParameterType

	// Indicates the availability and contains the reason, if not available.
	Available TegnAvailability
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
	return TegnParameterToBool(p.Value)
}

func TegnParameterToBool(s string) bool {
	return s == "y"
}

func (p *TegnParameter) ToInt(fallback int) int {
	return TegnParameterToInt(p.Value, fallback)
}

func TegnParameterToInt(s string, fallback int) int {
	var result *int
	_, err := fmt.Sscanf(s, "%d", result)
	if err != nil {
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
	// The features parameter contains features provided 
	// by previously selected Tegns (according to installation order).
	GetAvailability(features []string) TegnAvailability

	// Returns the features provided by this Tegn.
	// For Tegnsett, returns features provided by all child Tegns.
	GetFeatures() []string

	// Returns IDs of Tegns that must be installed before this one.
	// Only IDs from the same category take effect; IDs from other categories
	// are ignored.
	GetBeforeIDs() []string
}