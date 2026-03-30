package ui

import (
	tb "github.com/toliak/mce/tegnbuilder"
)

// UIView represents the current view state
type UIView int

const (
	ViewTegnsettList UIView = iota
	ViewTegnList
	ViewParameterList
)

// UIState manages the application state
type UIState struct {
	InitResult        tb.TegnsettInitializeResult
	OSInfExt 	  tb.OSInfoExt
	InstalledCache       tb.AvailablePackagesMap
	InstalledFeatures	 tb.TegnInstalledFeaturesMap

	EnabledIDsMap       tb.TegnGeneralEnabledIDsMap
	ParameterByIDMap     map[string]tb.TegnParameterMap

	CurrentView       UIView
	CurrentTegnsettID string
	CurrentTegnID     string
	ExitConfirmed     bool
	ExitError         error

	SelectionTegnsettID  int
	SelectionTegnID      int
	SelectionParameterID int
}

// NewUIState creates a new UI state
func NewUIState(
	initResult tb.TegnsettInitializeResult, 
	osInfoExt tb.OSInfoExt, 
	alreadyInstalled tb.AvailablePackagesMap,
) *UIState {
	state := &UIState{
		InitResult:        initResult,
		OSInfExt:        osInfoExt,
		InstalledCache: alreadyInstalled,
		EnabledIDsMap:     make(tb.TegnGeneralEnabledIDsMap),
		ParameterByIDMap:  make(map[string]tb.TegnParameterMap),
		ExitConfirmed:    false,
		SelectionTegnsettID: -1,
		SelectionTegnID: -1,
		SelectionParameterID: -1,
	}
	// Initialize maps for each tegnsett
	for k, _ := range initResult.AllIDsSet {
		state.EnabledIDsMap[k] = false
	}
	
	return state
}

// ToggleTegnsett toggles the enabled state of a tegnsett
func (s *UIState) ToggleID(id string) {
	s.EnabledIDsMap[id] = !s.EnabledIDsMap[id]
}

func (s *UIState) GetAndResetTegnsettSelection() int {
	result := s.SelectionTegnsettID
	s.SelectionTegnsettID = -1
	return result
}

func (s *UIState) GetAndResetTegnSelection() int {
	result := s.SelectionTegnID
	s.SelectionTegnID = -1
	return result
}

func (s *UIState) GetAndResetParameterSelection() int {
	result := s.SelectionParameterID
	s.SelectionParameterID = -1
	return result
}

// SetParameterValue sets a parameter value
// func (s *UIState) SetParameterValue(tegnID, paramName, value string) {
// 	// TODO: handle the error somehow
// 	s.InitResult.TegnByID[tegnID].SetParameter(paramName, value)
// }

// // GetParameterValue gets a parameter value
// func (s *UIState) GetParameterValue(tegnID, paramName string) string {
// 	params := s.InitResult.TegnByID[tegnID].GetParameters()
// 	index := slices.IndexFunc(params, func(p tegnbuilder.TegnParameter) bool {
// 		return p.Name == paramName
// 	})
// 	if index == -1 {
// 		return ""
// 	}

// 	return params[index].Value
// }

// // GetEnabledSummary returns a summary of enabled items
// func (s *UIState) GetEnabledSummary() (int, int, int) {
// 	tegnsettCount := 0
// 	tegnCount := 0
// 	paramCount := 0
	
// 	for tsID, enabled := range s.EnabledTegnsetts {
// 		if enabled {
// 			tegnsettCount++
// 			if tegns, ok := s.EnabledTegns[tsID]; ok {
// 				for tegnID, tegnEnabled := range tegns {
// 					if tegnEnabled {
// 						tegnCount++
// 						if params, ok := s.ParameterValues[tsID][tegnID]; ok {
// 							paramCount += len(params)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
	
// 	return tegnsettCount, tegnCount, paramCount
// }