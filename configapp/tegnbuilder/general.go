// Package tegn provides types and interfaces for defining installable
// configuration units.
// "Tegn" (Norwegian Bokmål for "character/symbol/sign")
// represents atomic functionality units to avoid collision with the "package" keyword.
package tegnbuilder

import (
	"fmt"
	"slices"
	"strings"

	"github.com/toliak/mce/optional"
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
type TegnInstalledFeaturesMap map[TegnFeature]bool

// Key -- Parameter ID
// Value -- Parameter Value
type TegnParameterMap map[string]string

// Key -- Tegn or Tegnsett ID
// Value -- is enabled (default -- false)
type TegnGeneralEnabledIDsMap map[string]bool

// The type MUST match [TegnGeneral.GetAvailability] function.
//
// osInfo must not be nil.
type TegnAvailabilityFunc func (osInfo *OSInfoExt, before TegnInstalledFeaturesMap, enabledIds TegnGeneralEnabledIDsMap) TegnAvailability

// TegnGeneral defines the common interface for all Tegn types.
type TegnGeneral interface {
	// Returns the unique identifier.
	GetID() string

	// Returns the human-readable name.
	GetName() string

	// Returns a detailed description.
	GetDescription() string

	// Returns the CPU architectures this Tegn supports.
	// An empty Optional means the Tegn is available on all architectures.
	// This condition is evaluated before GetAvailability.
	GetAvailableCPUArch() optional.Optional[[]data.CPUArchE]

	// Returns the OS types this Tegn supports.
	// An empty Optional means the Tegn is available on all OS types.
	// This condition is evaluated before GetAvailableCPUArch.
	GetAvailableOsType() optional.Optional[[]data.OSTypeE]

	// Checks whether this Tegn is available for installation.
	// Received all features that are available just before that Tegn (from previous Tegns and Tegnsetts).
	//
	// In Tegnsett the method must not accumulate the children data!
	// osInfo must not be nil.
	GetAvailability(osInfo *OSInfoExt, before TegnInstalledFeaturesMap, enabledIds TegnGeneralEnabledIDsMap) TegnAvailability

	// Returns IDs of Tegns that must be installed before this one.
	// Only IDs from the same category take effect; IDs from other categories
	// are ignored.
	GetBeforeIDs() []string
}

// Availability checker. Depends on the:
// - CPU arch
// - OS Type
// - own Tegn's Availability method
func GetTegnGeneralAvailable(obj TegnGeneral, osInfo *OSInfoExt, before TegnInstalledFeaturesMap, selectedIDs TegnGeneralEnabledIDsMap) TegnAvailability {
	cpu, hasCPU := obj.GetAvailableCPUArch().Get()
	if hasCPU {
		if !slices.ContainsFunc(cpu, func(v data.CPUArchE) bool {
			return v == osInfo.Arch.V
		}) {
			return NewTegnNotAvailable(
				fmt.Sprintf("Not available for the CPU Arch %v", osInfo.Arch),
			)
		}
	}

	osType, hasOsType := obj.GetAvailableOsType().Get()
	if hasOsType {
		if !slices.ContainsFunc(osType, func(v data.OSTypeE) bool {
			return v == osInfo.OsType.V
		}) {
			return NewTegnNotAvailable(
				fmt.Sprintf("Not available for the OS Type %v", osInfo.OsType),
			)
		}
	}

	return obj.GetAvailability(osInfo, before, selectedIDs)
}

func CheckFeatures(before TegnInstalledFeaturesMap, expected []TegnFeature) error {
	notBeforeNames := make([]string, 0)
	for _, v := range expected {
		if before[v] == false {
			notBeforeNames = append(notBeforeNames, string(v))
		}
	}

	if len(notBeforeNames) == 0 {
		return nil
	}

	notBeforeStr := strings.Join(notBeforeNames, ", ")
	return fmt.Errorf("Features not found: %s", notBeforeStr)
}