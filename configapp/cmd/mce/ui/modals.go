package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	tb "github.com/toliak/mce/tegnbuilder"
)

// NewHelpModal creates a help modal
func NewHelpModal(state *UIState, app *tview.Application, tegnGeneral *tb.TegnGeneral, closer func()) *tview.Flex {
	var helpText string

	if tegnGeneral == nil {
		helpText = `No help available`
	} else {
		helpText = (*tegnGeneral).GetDescription()
	}
	
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetText(helpText)

	textView.SetBackgroundColor(tcell.ColorDarkBlue)

	// button := tview.NewButton("Close").
	// 	SetSelectedFunc(func() {
	// 		closer()
	// 	})

	form := tview.NewFlex().
		SetDirection(tview.FlexRow)

	form.AddItem(textView, 0, 1, true)
	// form.AddItem(button, 1, 0, false)
	
	// AddItem(textView).
	// AddFormItem(textView).
	// AddButton("Close", func() {
	// 	closer()
	// })
	
	form.SetTitle("Help").
		SetBackgroundColor(tcell.ColorDarkBlue).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				closer()
				return nil
			}
			return event
		})

	// textView.Box
	// form.SetInputCapture()

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