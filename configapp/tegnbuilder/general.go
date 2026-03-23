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

type TegnFeature string

// Key -- Parameter ID
// Value -- Parameter Value
type TegnParameterMap map[string]string

// Key -- Tegn or Tegnsett ID
// Value -- is enabled (default -- false)
type TegnGeneralEnabledIDsMap map[string]bool

// The type MUST match [TegnGeneral.GetAvailability] function
type TegnAvailabilityFunc func (osInfo OSInfoExt, before []TegnFeature, enabledIds TegnGeneralEnabledIDsMap) TegnAvailability

// TegnGeneral defines the common interface for all Tegn types.
type TegnGeneral interface {
	// Returns the unique identifier.
	GetID() string

	// Returns the human-readable name.
	GetName() string

	// Returns a detailed description.
	GetDescription() string

	// Returns the CPU architectures this Tegn supports.
	// Returns nil if the Tegn is available on all architectures.
	// This condition is evaluated before GetAvailability.
	GetAvailableCPUArch() *[]data.CPUArchE

	// Returns the OS types this Tegn supports.
	// Returns nil if the Tegn is available on all OS types.
	// This condition is evaluated before GetAvailableCPUArch.
	GetAvailableOsType() *[]data.OSTypeE

	// Checks whether this Tegn is available for installation.
	// Received all features that are available just before that Tegn (from previous Tegns and Tegnsetts).
	//
	// In Tegnsett the method must not accumulate the children data!
	GetAvailability(osInfo OSInfoExt, before []TegnFeature, enabledIds TegnGeneralEnabledIDsMap) TegnAvailability

	// Returns IDs of Tegns that must be installed before this one.
	// Only IDs from the same category take effect; IDs from other categories
	// are ignored.
	GetBeforeIDs() []string
}

// Availability checker. Depends on the:
// - CPU arch
// - OS Type
// - own Tegn's Availability method
func GetTegnGeneralAvailable(obj TegnGeneral, osInfo OSInfoExt, before []TegnFeature, selectedIDs TegnGeneralEnabledIDsMap) TegnAvailability {
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

	return obj.GetAvailability(osInfo, before, selectedIDs)
}
