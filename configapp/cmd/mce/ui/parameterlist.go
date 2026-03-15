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
		// paramDesc := param.Description
		paramValue := param.GetValue()
		paramType := param.ParamType.String()
		
		secondaryText := fmt.Sprintf("%s = [yellow]%s[white]", paramName, paramValue)
		if param.Available.Available {
			secondaryText = fmt.Sprintf("(%- 8s) %s", paramType, secondaryText)
		} else {
			secondaryText = fmt.Sprintf("--- %s", secondaryText)
		}
		
		list.AddItem(secondaryText, "", 0, func() {
			// Enter key pressed - edit parameter
			if param.Available.Available {
				state.SelectionParameterID = list.GetCurrentItem()
				app.showParameterEditModal(state, tegnID, param)
			}
		})
	}

	if sel := state.GetAndResetParameterSelection(); sel != -1 && sel < list.GetItemCount(){
		list.SetCurrentItem(sel)
	}

	// Set up key bindings for the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyLeft:
			state.SelectionParameterID = list.GetCurrentItem()
			app.navigateBack()
			return nil
		case tcell.KeyEnter, tcell.KeyRight:
			index := list.GetCurrentItem()
			state.SelectionParameterID = index
			if index < list.GetItemCount() {
				list.GetItemSelectedFunc(index)()
			}
			// Handled by item selection
			return nil
		case tcell.KeyRune:
			if event.Rune() == '?' {
				state.SelectionParameterID = list.GetCurrentItem()
				app.showHelpModal(nil)
				return nil
			}
			if event.Rune() == ' ' {
				// Do not act on Space
				return nil
			}
		}
		return event
	})

	return list
}

// WARN: do not change param value here!
func (a *App) showParameterEditModal(state *UIState, tegnID string, param tegnbuilder.TegnParameter) {
	newValue := param.GetValue()

	// For now, use a simple form modal
	form := tview.NewForm().
		AddTextView("Description", param.Description, 0, 2, false, true).
		AddInputField("Value", param.GetValue(), 0, nil, func (v string) {
			newValue = v
		}).
		AddButton("Save", func() {
			// Get value from form and save
			state.InitResult.TegnByID[tegnID].SetParameter(param.Name, newValue)
			a.pages.RemovePage("parameterEditModal")

			// Refresh
			a.showParameterList(state.CurrentTegnID)
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("parameterEditModal")

			// Refresh
			a.showParameterList(state.CurrentTegnID)
		})

	form.SetFocus(1); // Focus on the value
	
	form.SetBorder(true).
		SetTitle("Edit Parameter").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				a.pages.RemovePage("parameterEditModal")

				// Refresh
				a.showParameterList(state.CurrentTegnID)
				return nil
			}
			return event
		})
	
	a.pages.AddPage("parameterEditModal", form, true, true)
}