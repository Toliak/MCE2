package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	tb "github.com/toliak/mce/tegnbuilder"
)

// TegnList creates the tegn list view within a tegnsett
func NewTegnList(
	state *UIState,
	tegnsettID string,
	order []string,
	availability tb.TegnGeneralAvailabilityByID,
	app *App,
) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)
	
	for _, id := range order {
		tegn := state.InitResult.TegnByID[id]
		tegnName := tegn.GetName()
		enabled := state.EnabledIDsMap[id]

		isAvailable := availability[id].Available
		isInstalled := state.InstalledCache[id]
		canBeOpened := !isInstalled && isAvailable && enabled && len(tegn.GetParameters(state.OSInfExt)) > 0

		status := "[gray]---[white]"
		
		if isInstalled {
			status = "[yellow]-+-[white]"
		} else if isAvailable {
			if enabled {
				status = "[green](*)[white]"
			} else {
				status = "[red]( )[white]"
			}
		}

		statusEnd := ""
		if canBeOpened {
			statusEnd = ">>"
		}

		secondaryText := fmt.Sprintf("%s %s %s", status, tegnName, statusEnd)

		list.AddItem(secondaryText, "", 0, func() {
			// Enter key pressed
			if canBeOpened {
				state.SelectionTegnID = list.GetCurrentItem()
				app.showParameterList(id)
			}
		})
	}

	if sel := state.GetAndResetTegnSelection(); sel != -1 && sel < list.GetItemCount(){
		list.SetCurrentItem(sel)
	}

	// Set up key bindings for the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyLeft:
			state.SelectionTegnID = list.GetCurrentItem()
			app.navigateBack()
			return nil
		case tcell.KeyEnter, tcell.KeyRight:
			// Handled by item selection
			index := list.GetCurrentItem()
			state.SelectionTegnID = index
			if index < list.GetItemCount() {
				list.GetItemSelectedFunc(index)()
			}
			return nil
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				// Toggle enabled state
				index := list.GetCurrentItem()
				if index >= 0 && index < len(order) {
					id := order[index]
					
					if !state.InstalledCache[id] && availability[id].Available {
						state.ToggleID(id)

						// Refresh the list to show updated status
						state.SelectionTegnID = list.GetCurrentItem()
						app.showTegnList(tegnsettID)
					}
				}
				return nil
			}
			if event.Rune() == '?' {
				index := list.GetCurrentItem()
				state.SelectionTegnID = index
				if index >= 0 && index < len(order) {
					id := order[index]
					tegn := state.InitResult.TegnByID[id].(tb.TegnGeneral)
					availability := availability[id]
					app.showHelpModal(&tegn, &availability)
					return nil
				}
				app.showHelpModal(nil, nil)
				return nil
			}
		}
		return event
	})
	
	return list
}