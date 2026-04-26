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

		status := "[gray]---[-]"
		if isAvailable {
			if enabled {
				status = "[green](*)[-]"
			} else {
				status = "[red]( )[-]"
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
				state.SelectionTegnsettID = list.GetCurrentItem()
				app.showTegnList(id)
			}
		})
	}

	if sel := state.GetAndResetTegnsettSelection(); sel != -1 && sel < list.GetItemCount(){
		list.SetCurrentItem(sel)
	}

	// Set up key bindings for the list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// fmt.Println("yeee22sss!")
		switch event.Key() {
		case tcell.KeyEscape:
			// Can't go back from root view
			state.SelectionTegnsettID = list.GetCurrentItem()
			app.showExitModal()
			return event
		case tcell.KeyEnter, tcell.KeyRight:
			// Handled by item selection
			index := list.GetCurrentItem()
			state.SelectionTegnsettID = index
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
					if availability[id].Available {
						state.ToggleID(id)

						// Refresh the list to show updated status
						state.SelectionTegnsettID = index
						app.showTegnsettList()
					}
				}
				return nil
			}
			if event.Rune() == '?' {
				index := list.GetCurrentItem()
				state.SelectionTegnsettID = index
				if index >= 0 && index < len(order) {
					id := order[index]
					availability := availability[id]
					tegn := state.InitResult.TegnsettByID[id].(tb.TegnGeneral)
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
