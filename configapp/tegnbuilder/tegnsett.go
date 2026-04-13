package tegnbuilder

import (
	"maps"
	"fmt"
)

// Tegnsett represents a group of related Tegns.
type Tegnsett interface {
	TegnGeneral

	// Returns all child Tegns in this group.
	GetChildren() []Tegn

	// Post-install stage, after all Tegns (in ths Tegnsett) were installed.
	// 
	// The installedTegns contains only installed Tegns from the current Tegnsett.
	// The "already" includes features that were added by the children Tegns.
	ExecPostInstall(
		installedTegns []Tegn,
		osInfo OSInfoExt,
		already TegnInstalledFeaturesMap,
		tegnToParams map[string]TegnParameterMap,
	) error
}

type TegnsettBuildFunc func() Tegnsett
// type TegnsettOuterBuildFunc func (children []TegnBuildFunc) TegnsettBuildFunc

type TegnsettInitializeError struct {
	what string
}

var _ error = (*TegnsettInitializeError)(nil)

func (e *TegnsettInitializeError) Error() string {
	return e.what
}

type MapTegnsettByID map[string]Tegnsett
type MapTegnByID map[string]Tegn

type TegnsettInitializeResult struct {
	TegnsettByID MapTegnsettByID
	TegnByID     MapTegnByID
	AllIDsSet    map[string]struct{}
}

func (d *MapTegnsettByID) GetTegnsettIDs() []string {
	result := make([]string, 0, len(*d))
	for k, _ := range *d {
		result = append(result, k)
	}

	return result
}

func (d *MapTegnsettByID) GetTegnsetts() []Tegnsett {
	result := make([]Tegnsett, 0, len(*d))
	for _, v := range *d {
		result = append(result, v)
	}

	return result
}

// Initializes all Tegnsetts using their build functions.
// Also initializes all the children (Tegns).
func InitializeAllTegnsetts(tegnsetts []TegnsettBuildFunc) (*TegnsettInitializeResult, error) {
	tegnsettByID := make(map[string]Tegnsett, len(tegnsetts))
	tegnByID := make(map[string]Tegn /*probably capacity*/, len(tegnsetts))
	allIDsSet := make(map[string]struct{} /*probably capacity*/, len(tegnsetts)*2)

	for _, v := range tegnsetts {
		tegnsett := v()
		tegnsettID := tegnsett.GetID()
		if _, ok := allIDsSet[tegnsettID]; ok {
			return nil, &TegnsettInitializeError{
				what: fmt.Sprintf("Tegnsett ID '%s' is already taken", tegnsettID),
			}
		}
		allIDsSet[tegnsettID] = struct{}{}
		tegnsettByID[tegnsettID] = tegnsett

		for _, child := range tegnsett.GetChildren() {
			childID := child.GetID()
			if _, ok := allIDsSet[childID]; ok {
				return nil, &TegnsettInitializeError{
					what: fmt.Sprintf("Tegn ID '%s' (inside Tegnsett '%s') is already taken", childID, tegnsettID),
				}
			}

			allIDsSet[childID] = struct{}{}
			tegnByID[childID] = child
		}
	}

	return &TegnsettInitializeResult{
		TegnsettByID: tegnsettByID,
		TegnByID:     tegnByID,
		AllIDsSet:    allIDsSet,
	}, nil
}

type TegnsettsOrderResult struct {
	Tegnsett         []string
	TegnByTegnsettID map[string][]string
}

func GetTegnsettsOrder(tegnsettByID MapTegnsettByID) (*TegnsettsOrderResult, error) {
	generalArray := make([]TegnGeneral, 0, len(tegnsettByID))
	for _, v := range tegnsettByID {
		generalArray = append(generalArray, v.(TegnGeneral))
	}

	tegnsettOrder, err := GetGeneralOrder(generalArray)
	if err != nil {
		return nil, err
	}

	tegnOrderByTegnsettID := make(map[string][]string, len(tegnsettOrder))
	for _, id := range tegnsettOrder {
		tegnsett := tegnsettByID[id]
		children := tegnsett.GetChildren()

		generalArray := make([]TegnGeneral, 0, len(tegnsettByID))
		for _, v := range children {
			generalArray = append(generalArray, v.(TegnGeneral))
		}

		tegnOrder, err := GetGeneralOrder(generalArray)
		if err != nil {
			return nil, fmt.Errorf("In Tegnsett '%s': %w", id, err)
		}

		tegnOrderByTegnsettID[id] = tegnOrder
	}

	return &TegnsettsOrderResult{
		Tegnsett:         tegnsettOrder,
		TegnByTegnsettID: tegnOrderByTegnsettID,
	}, nil
}

type TegnGeneralAvailabilityByID  map[string]TegnAvailability

func GetTegnsettsAvailability(
	osInfo OSInfoExt,
	orders TegnsettsOrderResult,
	tegnsettByID MapTegnsettByID,
	tegnByID MapTegnByID,
	selectedIDs TegnGeneralEnabledIDsMap,
	alreadyInstalledFeatures TegnInstalledFeaturesMap,
) TegnGeneralAvailabilityByID {
	availableByID := make(TegnGeneralAvailabilityByID /* probably capacity */, len(tegnByID))

	currentFeatures := make(TegnInstalledFeaturesMap)
	maps.Copy(currentFeatures, alreadyInstalledFeatures)

	for _, tegnsettID := range orders.Tegnsett {
		tegnsett := tegnsettByID[tegnsettID]

		available := GetTegnGeneralAvailable(tegnsett, osInfo, currentFeatures, selectedIDs)
		availableByID[tegnsettID] = available
		if !available.Available {
			continue
		}
		if v, ok := selectedIDs[tegnsettID]; !(ok && v) {
			continue
		}

		for _, tegnID := range orders.TegnByTegnsettID[tegnsettID] {
			tegn := tegnByID[tegnID]

			available := GetTegnGeneralAvailable(tegn, osInfo, currentFeatures, selectedIDs)
			availableByID[tegnID] = available
			if !available.Available {
				continue
			}
			if v, ok := selectedIDs[tegnID]; !(ok && v) {
				continue
			}

			maps.Copy(currentFeatures, tegn.GetFeatures())
		}
	}

	return availableByID
}
