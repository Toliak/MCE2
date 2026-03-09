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
		canBeOpened := isAvailable && enabled && len(tegn.GetParameters()) > 0

		status := "[gray]---[white]"
		
		if isAvailable {
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
				app.showParameterList(id)
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
			return event
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				// Toggle enabled state
				index := list.GetCurrentItem()
				if index >= 0 && index < len(order) {
					id := order[index]
					
					if availability[id].Available {
						state.ToggleID(id)

						// Refresh the list to show updated status
						app.showTegnList(tegnsettID)
					}
				}
				return nil
			}
		}
		return event
	})
	
	return list
}