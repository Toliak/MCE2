package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	tb "github.com/toliak/mce/tegnbuilder"
)

// TegnsettList creates the tegnsett selector view
func NewTegnsettList(
	state *UIState, 
	order []string, 
	availability tb.TegnGeneralAvailabilityByID, 
	app *App,
) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)
	
	for _, id := range order {
		tegnsett := state.InitResult.TegnsettByID[id]

		// tsName := ts.GetName()
		tsDesc := tegnsett.GetDescription()

		isAvailable := availability[id].Available
		enabled := state.EnabledIDsMap[id]
		canBeOpened := isAvailable && enabled && len(tegnsett.GetChildren()) != 0

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

		secondaryText := fmt.Sprintf("%s %s %s", status, tsDesc, statusEnd)
		
		list.AddItem(secondaryText, "", 0, func() {
			// Enter key pressed - open tegnsett
			if canBeOpened {
				app.showTegnList(id)
			}
		})
	}
	
	// Set up key bindings for the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// fmt.Println("yeee22sss!")
		switch event.Key() {
		case tcell.KeyEscape:
			// Can't go back from root view
			return event
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
						app.showTegnsettList()
					}
				}
				return nil
			}
		}
		return event
	})
	
	return list
}