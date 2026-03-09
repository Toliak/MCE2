package ui

import (
	"github.com/toliak/mce/inspector"
	"github.com/toliak/mce/tegnbuilder"
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
	InitResult        tegnbuilder.TegnsettInitializeResult
	HarvestData 	  inspector.HarvestData

	EnabledIDsMap     tegnbuilder.EnabledIDsMap
	CurrentView       UIView
	CurrentTegnsettID string
	CurrentTegnID     string
	ExitConfirmed     bool
	ExitError         error
}

// NewUIState creates a new UI state
func NewUIState(initResult        tegnbuilder.TegnsettInitializeResult, harvestData inspector.HarvestData) *UIState {
	state := &UIState{
		InitResult:        initResult,
		HarvestData:        harvestData,
		EnabledIDsMap:     make(tegnbuilder.EnabledIDsMap),
		ExitConfirmed:    false,
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