package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/toliak/mce/tegnbuilder"
	// "strings"
)

// ParameterList creates the parameter list view for a tegn
func NewParameterList(state *UIState, tegnID string, app *App) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)
	
	tegn := state.InitResult.TegnByID[tegnID]

	if tegn == nil {
		list.AddItem("Error", "Tegn not found", 0, nil)
		return list
	}

	params := tegn.GetParameters()
	
	for _, param := range params {
		paramName := param.Name
		paramDesc := param.Description
		paramValue := param.GetValue()
		paramType := param.ParamType.String()
		
		secondaryText := fmt.Sprintf("[%s] %s = [yellow]%s[white]", paramType, paramDesc, paramValue)
		if !param.Available.Available {
			secondaryText = fmt.Sprintf("--- %s", secondaryText)
		}
		
		list.AddItem(paramName, secondaryText, 0, func() {
			// Enter key pressed - edit parameter
			if param.Available.Available {
				app.showParameterEditModal(state, tegnID, param)
			}
		})
	}
	
	// Set up key bindings for the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.navigateBack()
			return nil
		case tcell.KeyEnter:
			// Handled by item selection
			return nil
		}
		return event
	})
	
	return list
}

// WARN: do not change param value here!
func (a *App) showParameterEditModal(state *UIState, tegnID string, param tegnbuilder.TegnParameter) {
	// modal := tview.NewModal().
	// 	SetText(fmt.Sprintf("Edit parameter: %s", param.Name)).
	// 	AddButtons([]string{"Save", "Cancel"}).
	// 	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 		if buttonLabel == "Save" {
	// 			// Get the input value from the form
	// 			// This would nee`d a form modal instead of simple modal
	// 		}
	// 		a.pages.RemovePage("parameterEditModal")
	// 	})
	
	// TODO: provide accept predicate

	newValue := param.GetValue()

	// For now, use a simple form modal
	form := tview.NewForm().
		AddInputField("Value", param.GetValue(), 0, nil, func (v string) {
			newValue = v
		}).
		AddButton("Save", func() {
			// Get value from form and save
			state.InitResult.TegnByID[tegnID].SetParameter(param.Name, newValue)

			a.pages.RemovePage("parameterEditModal")
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("parameterEditModal")
		})
	
	form.SetBorder(true).
		SetTitle("Edit Parameter").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				a.pages.RemovePage("parameterEditModal")
				return nil
			}
			return event
		})
	
	a.pages.AddPage("parameterEditModal", form, true, true)
}