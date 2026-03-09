package ui

import (
	// "fmt"
	// "strings"

	// "github.com/gdamore/tcell/v2"
	// "github.com/rivo/tview"
)

// ShowSearchModal creates and shows a search modal with proper functionality
// func (a *App) showSearchModal2() {
// 	form := tview.NewForm().
// 		AddInputField("Enter Tegn/Tegnsett ID", "", 40, nil, nil).
// 		AddButton("Search", func() {
// 			// Get the input field
// 			inputField := form.GetFormItem(0).(*tview.InputField)
// 			searchID := inputField.GetText()
			
// 			if searchID == "" {
// 				return
// 			}
			
// 			// Search for the ID
// 			found := a.searchByID(searchID)
			
// 			if !found {
// 				// Show error modal
// 				errorModal := tview.NewModal().
// 					SetText(fmt.Sprintf("ID '%s' not found", searchID)).
// 					AddButtons([]string{"OK"}).
// 					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 						a.pages.RemovePage("searchModal")
// 						a.pages.RemovePage("errorModal")
// 					})
// 				a.pages.AddPage("errorModal", errorModal, true, true)
// 				return
// 			}
			
// 			a.pages.RemovePage("searchModal")
// 		}).
// 		AddButton("Cancel", func() {
// 			a.pages.RemovePage("searchModal")
// 		})
	
// 	form.SetBorder(true).
// 		SetTitle("Search by ID").
// 		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 			if event.Key() == tcell.KeyEscape {
// 				a.pages.RemovePage("searchModal")
// 				return nil
// 			}
// 			return event
// 		})
	
// 	a.pages.AddPage("searchModal", form, true, true)
// }

// func (a *App) searchByID(id string) bool {
// 	id = strings.TrimSpace(id)
	
// 	// Search in Tegnsetts
// 	for _, ts := range a.state.Tegnsetts {
// 		if ts.GetID() == id {
// 			// Open this tegnsett
// 			a.state.CurrentTegnsettID = ts.GetID()
// 			a.showTegnList(ts.GetID())
// 			return true
// 		}
		
// 		// Search in children Tegns
// 		for _, tegn := range ts.GetChildren() {
// 			if tegn.GetID() == id {
// 				// Open this tegn's parameters
// 				a.state.CurrentTegnsettID = ts.GetID()
// 				a.state.CurrentTegnID = tegn.GetID()
// 				a.showTegnList(ts.GetID())
// 				// Would need to navigate to parameter list
// 				return true
// 			}
// 		}
// 	}
	
// 	return false
// }