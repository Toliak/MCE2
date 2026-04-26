package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	tb "github.com/toliak/mce/tegnbuilder"
)

// NewHelpModal creates a help modal
func NewHelpModal(installedCache tb.AvailablePackagesMap, app *tview.Application, tegnGeneralPtr *tb.TegnGeneral, availability *tb.TegnAvailability, closer func()) *tview.Flex {
	var helpText string

	if tegnGeneralPtr == nil {
		helpText = `No help available`
	} else {
		tegnGeneral := *tegnGeneralPtr

		var sb strings.Builder
		sb.WriteString("[::b]Description[::-]\n")
		sb.WriteString(tegnGeneral.GetDescription())

		sb.WriteString("\n\n[::b]Availability[::-]\n")
		sb.WriteString("CPU Architectures: ")
		cpuArchList := tegnGeneral.GetAvailableCPUArch()
		if cpuArchList == nil {
			sb.WriteString("any\n")
		} else {
			if len(*cpuArchList) == 0 {
				sb.WriteString("nothing! (error?)\n")
			} else {
				sb.WriteString("\n")
				for _, arch := range *cpuArchList {
					sb.WriteString("- ")
					sb.WriteString(arch.String())
					sb.WriteString("\n")
				}
			}
		}

		sb.WriteString("OS Types: ")
		osTypeList := tegnGeneral.GetAvailableOsType()
		if osTypeList == nil {
			sb.WriteString("any\n")
		} else {
			if len(*osTypeList) == 0 {
				sb.WriteString("nothing! (error?)\n")
			} else {
				sb.WriteString("\n")
				for _, osType := range *osTypeList {
					sb.WriteString("- ")
					sb.WriteString(osType.String())
					sb.WriteString("\n")
				}
			}
		}

		sb.WriteString("\n")
		if availability == nil {
			sb.WriteString("No detailed availability info\n")
		} else {
			sb.WriteString("Available: ")
			if availability.Available {
				sb.WriteString("yes\n")
			} else {
				sb.WriteString("no\n")
				sb.WriteString("Reason: ")
				sb.WriteString(availability.Reason)
				sb.WriteString("\n")
			}
			
		}
		
		// tegnGeneral.
		if tegn, ok := tegnGeneral.(tb.Tegn); ok {
			sb.WriteString("\n[::b]Tegn Features[::-]\n")
			featureList := tegn.GetFeatures()
			if len(featureList) == 0 {
				sb.WriteString("no features (error?)\n")
			} else {
				for ft, _ := range featureList {
					sb.WriteString("- ")
					sb.WriteString(string(ft))
					sb.WriteString("\n")
				}
			}

			sb.WriteString("\n[::b]Tegn Already installed[::-]: ")
			if installedCache[tegnGeneral.GetID()] {
				sb.WriteString("yes\n")
			} else {
				sb.WriteString("no\n")
			}
		} else if _, ok := tegnGeneral.(tb.Tegnsett); ok {
			// No additional info
		}

		helpText = sb.String()
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
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
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
	
	confirmText := `[::b]Confirmation Summary[::]

Do you want to apply these configuration changes?`

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