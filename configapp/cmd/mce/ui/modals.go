package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewHelpModal creates a help modal
func NewHelpModal(state *UIState, app *tview.Application, closer func()) *tview.Form {
	helpText := `
[::b]Keyboard Shortcuts[::]
  [yellow]F1[white] or [yellow]?[white]      - Show this help
  [yellow]F3[white] or [yellow]/[white]      - Search by ID
  [yellow]F5[white] or [yellow]Ctrl-X[white] - Confirm selection
  [yellow]F10[white] or [yellow]Ctrl-C[white]- Exit application
  [yellow]Space[white]          - Toggle enable/disable
  [yellow]Enter[white]          - Open item / Edit parameter
  [yellow]Esc[white]            - Go back

[::b]Navigation[::]
  [yellow]↑/↓[white] or [yellow]j/k[white]   - Navigate list
  [yellow]Enter[white]          - Select item

[::b]Current View[::]
`
	
	switch state.CurrentView {
	case ViewTegnsettList:
		helpText += "  Tegnsett Selector - Choose configuration categories"
	case ViewTegnList:
		helpText += fmt.Sprintf("  Tegn List - Configure items in '%s'", state.CurrentTegnsettID)
	case ViewParameterList:
		helpText += fmt.Sprintf("  Parameters - Edit settings for '%s'", state.CurrentTegnID)
	}
	
	helpText += `

[::b]Legend[::]
  [green]✓[white] - Enabled
  [red]○[white] - Disabled
  [yellow][Params][white] - Has configurable parameters
`
	
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetText(helpText)

	// TODO: flex
	form := tview.NewForm().
		AddFormItem(textView).
		AddButton("Close", func() {
			closer()
		})
	
	form.SetBorder(true).
		SetTitle("Help").
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				closer()
				return nil
			}
			return event
		})
	
	return form
}

// NewSearchModal creates a search modal
func NewSearchModal(state *UIState, app *tview.Application, closer func()) *tview.Form {
	form := tview.NewForm().
		AddInputField("Search ID", "", 40, nil, nil).
		AddButton("Search", func() {
			// Get the input value and search
			// This would need to be implemented properly

			// TODO: TODO:
		}).
		AddButton("Cancel", func() {
			closer()
		})
	
	form.SetBorder(true).
		SetTitle("Search by ID").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				closer()
				return nil
			}
			return event
		})
	
	return form
}

// NewConfirmModal creates a confirmation modal
func NewConfirmModal(state *UIState, app *tview.Application, closer func()) *tview.Modal {
	// tsCount, tegnCount, paramCount := state.GetEnabledSummary()
	
	confirmText := fmt.Sprintf(`
[::b]Confirmation Summary[::]

Enabled Tegnsetts: [green]%d[white]
Enabled Tegns:     [green]%d[white]
Parameters:        [green]%d[white]

Do you want to apply these configuration changes?
`, -1, -1, -1)

	buttonYesClick := func() {
		state.ExitConfirmed = true
		app.Stop()
	}
	buttonNoClick := func() {
		closer()
	}
	
	modal := tview.NewModal().
		SetText(confirmText).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				buttonYesClick()
			} else {
				buttonNoClick()
			}
		})
	
	modal.SetBorder(true).
		SetTitle("Confirm Changes").
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEscape:
				buttonNoClick()
				return nil
			case tcell.KeyRune:
				if event.Rune() == 'y' {
					buttonYesClick()
				} else if event.Rune() == 'n' {
					buttonNoClick()
				}
				return nil
			}

			return event
		})
		
	
	return modal
}

// NewExitModal creates an exit confirmation modal
func NewExitModal(state *UIState, app *tview.Application, closer func()) *tview.Modal {
	// tsCount, tegnCount, _ := state.GetEnabledSummary()
	
	exitText := `
[::b]Exit Application[::]

Do you want to exit?
`
	
	modal := tview.NewModal().
		SetText(exitText).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				app.Stop()
			} else {
				closer()
			}
		})
	
	modal.SetBorder(true).
		SetTitle("Exit Confirmation").
		SetBackgroundColor(tcell.ColorDarkRed)
	
	return modal
}