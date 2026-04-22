package ui

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	tb "github.com/toliak/mce/tegnbuilder"
)

// UninstallState holds the state for the uninstall UI
type UninstallState struct {
	App            *App                 // Reference to the main app for data
	InstalledTegns []string             // List of installed Tegn IDs (sorted)
	SelectedTegns  []string             // IDs selected for uninstall
	ExitConfirmed  bool
}

// UninstallApp wraps the uninstall UI
type UninstallApp struct {
	app       *tview.Application
	pages     *tview.Pages
	rootView  *tview.Flex
	statusBar *tview.TextView
	State     *UninstallState
}

// NewUninstallApp creates a new uninstall UI application
func NewUninstallApp(
	initResult tb.TegnsettInitializeResult, 
	osInfoExt tb.OSInfoExt, 
	alreadyInstalled tb.AvailablePackagesMap,
	alreadyInstalledFeatures tb.TegnInstalledFeaturesMap,
) *UninstallApp {
	// Sort installed Tegns for consistent order
	// TODO: here
	sort.Strings(installedTegns)

	state := &UninstallState{
		App:            mainApp,
		InstalledTegns: installedTegns,
		SelectedTegns:  make([]string, 0),
		ExitConfirmed:  false,
	}

	ui := &UninstallApp{
		app:   tview.NewApplication(),
		State: state,
	}
	ui.setupUI()
	return ui
}

func (u *UninstallApp) setupUI() {
	u.rootView = tview.NewFlex().SetDirection(tview.FlexRow)

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b]Uninstall Tegns[::]").
		SetDynamicColors(true)

	u.pages = tview.NewPages()

	u.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	u.updateStatusBar()

	buttonBar := u.createButtonBar()

	u.rootView.AddItem(title, 1, 0, false)
	u.rootView.AddItem(u.pages, 0, 1, true)
	u.rootView.AddItem(u.statusBar, 1, 0, false)
	u.rootView.AddItem(buttonBar, 1, 0, false)

	u.app.SetRoot(u.rootView, true)
	u.setupGlobalKeyBindings()
}

func (u *UninstallApp) createButtonBar() *tview.Flex {
	buttonBar := tview.NewFlex().SetDirection(tview.FlexRow)
	buttons := tview.NewFlex().SetDirection(tview.FlexColumn)

	helpBtn := tview.NewButton(" (F1) Help ").
		SetSelectedFunc(func() { u.showHelpModal(nil) })

	confirmBtn := tview.NewButton(" (F5) Confirm ").
		SetSelectedFunc(func() { u.showConfirmModal() })

	exitBtn := tview.NewButton(" (F10) Exit ").
		SetSelectedFunc(func() { u.showExitModal() })

	buttons.AddItem(helpBtn, 0, 1, false)
	buttons.AddItem(confirmBtn, 0, 1, false)
	buttons.AddItem(exitBtn, 0, 1, false)

	buttonBar.AddItem(buttons, 0, 1, true)
	return buttonBar
}

func (u *UninstallApp) setupGlobalKeyBindings() {
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			// Show help for currently selected item if any
			index := u.getCurrentListIndex()
			if index >= 0 && index < len(u.State.InstalledTegns) {
				id := u.State.InstalledTegns[index]
				tegn := u.State.App.State.InitResult.TegnByID[id].(tb.TegnGeneral)
				u.showHelpModal(&tegn)
			} else {
				u.showHelpModal(nil)
			}
			return nil
		case tcell.KeyF5:
			u.showConfirmModal()
			return nil
		case tcell.KeyF10, tcell.KeyCtrlC:
			u.showExitModal()
			return nil
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				// Toggle selection
				index := u.getCurrentListIndex()
				if index >= 0 && index < len(u.State.InstalledTegns) {
					id := u.State.InstalledTegns[index]
					u.toggleSelection(id)
					u.refreshList()
				}
				return nil
			}
			if event.Rune() == '?' {
				index := u.getCurrentListIndex()
				if index >= 0 && index < len(u.State.InstalledTegns) {
					id := u.State.InstalledTegns[index]
					tegn := u.State.App.State.InitResult.TegnByID[id].(tb.TegnGeneral)
					u.showHelpModal(&tegn)
				} else {
					u.showHelpModal(nil)
				}
				return nil
			}
		}
		return event
	})
}

func (u *UninstallApp) getCurrentListIndex() int {
    primitive := u.pages.GetPage("tegnList")
    if primitive == nil {
        return -1
    }
    if list, ok := primitive.(*tview.List); ok {
        return list.GetCurrentItem()
    }
    return -1
}

func (u *UninstallApp) toggleSelection(id string) {
	found := false
	for i, selID := range u.State.SelectedTegns {
		if selID == id {
			// Remove
			u.State.SelectedTegns = append(u.State.SelectedTegns[:i], u.State.SelectedTegns[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		u.State.SelectedTegns = append(u.State.SelectedTegns, id)
	}
	u.updateStatusBar()
}

func (u *UninstallApp) updateStatusBar() {
	selectedCount := len(u.State.SelectedTegns)
	statusText := fmt.Sprintf("Selected: %d Tegns | [yellow]Space[white]: Toggle | [yellow]F1[white]: Help | [yellow]F5[white]: Confirm | [yellow]F10[white]: Exit", selectedCount)
	u.statusBar.SetText(statusText)
}

func (u *UninstallApp) showTegnList() {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	for _, id := range u.State.InstalledTegns {
		tegn := u.State.App.State.InitResult.TegnByID[id]
		name := tegn.GetName()
		selected := false
		for _, selID := range u.State.SelectedTegns {
			if selID == id {
				selected = true
				break
			}
		}

		status := "[red]( )[white]"
		if selected {
			status = "[green](*)[white]"
		}

		secondaryText := fmt.Sprintf("%s %s", status, name)
		list.AddItem(secondaryText, "", 0, nil)
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				index := list.GetCurrentItem()
				if index >= 0 && index < len(u.State.InstalledTegns) {
					id := u.State.InstalledTegns[index]
					u.toggleSelection(id)
					u.refreshList()
				}
				return nil
			}
			if event.Rune() == '?' {
				index := list.GetCurrentItem()
				if index >= 0 && index < len(u.State.InstalledTegns) {
					id := u.State.InstalledTegns[index]
					tegn := u.State.App.State.InitResult.TegnByID[id].(tb.TegnGeneral)
					u.showHelpModal(&tegn)
				} else {
					u.showHelpModal(nil)
				}
				return nil
			}
		}
		return event
	})

	u.pages.AddPage("tegnList", list, true, true)
	removeAllPagesExceptOne(*u.pages, "tegnList")
	u.updateStatusBar()
}

func (u *UninstallApp) refreshList() {
	u.showTegnList()
}

func (u *UninstallApp) showHelpModal(tegnGeneral *tb.TegnGeneral) {
	// Reuse the existing HelpModal from modals.go
	modal := NewHelpModal(u.State.App.State, u.app, tegnGeneral, nil, func() {
		u.pages.RemovePage("helpModal")
	})
	u.pages.AddPage("helpModal", modal, true, true)
}

func (u *UninstallApp) showConfirmModal() {
	selected := u.State.SelectedTegns
	text := fmt.Sprintf("Uninstall %d selected Tegns?\n\nThis action cannot be undone.", len(selected))
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				u.State.ExitConfirmed = true
				u.app.Stop()
			} else {
				u.pages.RemovePage("confirmModal")
			}
		})
	modal.SetBorder(true).SetTitle("Confirm Uninstall")
	u.pages.AddPage("confirmModal", modal, true, true)
}

func (u *UninstallApp) showExitModal() {
	modal := tview.NewModal().
		SetText("Exit without uninstalling?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				u.app.Stop()
			} else {
				u.pages.RemovePage("exitModal")
			}
		})
	modal.SetBorder(true).SetTitle("Exit")
	u.pages.AddPage("exitModal", modal, true, true)
}

// Run starts the uninstall UI
func (u *UninstallApp) Run() error {
	u.showTegnList()
	return u.app.Run()
}