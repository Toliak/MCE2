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

	params := tegn.GetParameters(state.OSInfExt)
	
	for _, param := range params {
		paramID := param.GetID()
		paramName := param.GetName()
		// paramDesc := param.Description
		paramValue, ok := state.ParameterByIDMap[tegnID][paramID]
		if !ok {
			paramValue = param.GetDefaultValue()
		}
		paramType := param.GetParamType().String()
		
		secondaryText := fmt.Sprintf("%s = [yellow]%s[white]", paramName, paramValue)
		if param.GetAvailable().Available {
			secondaryText = fmt.Sprintf("(%- 8s) %s", paramType, secondaryText)
		} else {
			secondaryText = fmt.Sprintf("--- %s", secondaryText)
		}

		list.AddItem(secondaryText, "", 0, func() {
			// Enter key pressed - edit parameter
			if param.GetAvailable().Available {
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
				app.showHelpModal(nil, nil)
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
	paramValue, ok := state.ParameterByIDMap[tegnID][param.GetID()]
	if !ok {
		paramValue = param.GetDefaultValue()
	}

	getDefaultDescriptionValue := func () string {
		return param.GetDescription()
	}

	var setDescriptionError func (value error)

	// For now, use a simple form modal
	form := tview.NewForm().
		AddTextView("Description", getDefaultDescriptionValue(), 0, 2, false, true).
		AddInputField("Value", paramValue, 0, nil, func (v string) {
			paramValue = v
		}).
		AddButton("Save", func() {
			// TODO: validate and add the error message
			// Get value from form and save
			// TODO: resolved?
			err := param.Validate(paramValue)
			if err != nil {
				setDescriptionError(err)
				return
			}

			state.ParameterByIDMap[tegnID][param.GetID()] = paramValue
			a.pages.RemovePage("parameterEditModal")

			// Refresh
			a.showParameterList(state.CurrentTegnID)
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("parameterEditModal")

			// Refresh
			a.showParameterList(state.CurrentTegnID)
		})

	setDescriptionError = func (value error) {
		newText := fmt.Sprintf("ERROR: %s\n\n%s", value, getDefaultDescriptionValue())
		// TODO: Create TextView before the `form`. And use AddFormItem to add it
		// This makes able to get rid of the "non-verifiable" GetFormItem(0)
		textView := form.GetFormItem(0).(*tview.TextView)
		textView.SetText(newText)
	}

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