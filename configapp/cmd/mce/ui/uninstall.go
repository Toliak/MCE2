package ui

import (
	"maps"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	tb "github.com/toliak/mce/tegnbuilder"
)

// UninstallUIState holds the state for the uninstall UI
type UninstallUIState struct {
	InitResult        tb.TegnsettInitializeResult
	OSInfExt          tb.OSInfoExt
	InstalledCache    tb.AvailablePackagesMap

	OriginalOrder     tb.TegnsettsOrderResult

	// Tegn IDs, installed only. In the "reversed" order
	FilteredOrder     []string

	// This field is saved if the installation fails
	EnabledIDsMap tb.TegnGeneralEnabledIDsMap

	ExitConfirmed bool
	ExitError     error

	SelectionTegnID  int
}

func (s *UninstallUIState) GetAndResetTegnSelection() int {
	result := s.SelectionTegnID
	s.SelectionTegnID = -1
	return result
}

// UninstallApp wraps the uninstall UI
type UninstallApp struct {
	app       *tview.Application
	pages     *tview.Pages
	rootView  *tview.Flex
	statusBar *tview.TextView
	State     *UninstallUIState
}

type UninstallAppInitData struct {
	InitResult               tb.TegnsettInitializeResult
	OSInfExt                tb.OSInfoExt
	AlreadyInstalled         tb.AvailablePackagesMap
	OriginalOrder     tb.TegnsettsOrderResult
	FilteredOrder     []string
	EnabledIDsMap tb.TegnGeneralEnabledIDsMap
}

// Returns bool and the tegn-blocker
func canTegnBeSelected(state *UninstallUIState, tegnIdToSelect string) (bool, string, string) {
	m := make(tb.TegnGeneralEnabledIDsMap)
	m[tegnIdToSelect] = true
	return CanTegnBeSelected(state, m)
}

// Returns bool and the tegn-blocker id and reason
func CanTegnBeSelected(state *UninstallUIState, tegnIdToSelect tb.TegnGeneralEnabledIDsMap) (bool, string, string) {
	remainedFeatures := make(tb.TegnInstalledFeaturesMap)
	mustRemainAvailable := make([]string, 0, len(state.FilteredOrder))
	for _, tegnId := range state.FilteredOrder {
		if state.EnabledIDsMap[tegnId] == true {
			// Do not include the features of the removing Tegn candidate
			continue
		}
		if tegnIdToSelect[tegnId] {
			// Do not include the features of the new candidate
			continue
		}

		mustRemainAvailable = append(mustRemainAvailable, tegnId)

		tegn := state.InitResult.TegnByID[tegnId]
		maps.Copy(remainedFeatures, tegn.GetFeatures())
	}

	// Workaround: we must pretend that all Tegnsetts are selected
	tegnsettEnabledIDsMap := make(tb.TegnGeneralEnabledIDsMap, len(state.InitResult.TegnsettByID))
	for k, _ := range state.InitResult.TegnsettByID {
		tegnsettEnabledIDsMap[k] = true
	}

	remainedFeaturesWithCandidate := make(tb.TegnInstalledFeaturesMap, len(remainedFeatures))
	maps.Copy(remainedFeaturesWithCandidate, remainedFeatures)
	for tegnId, _ := range tegnIdToSelect {
		maps.Copy(remainedFeaturesWithCandidate, state.InitResult.TegnByID[tegnId].GetFeatures())
	}

	availabilityWithCandidate := tb.GetTegnsettsAvailability(
		state.OSInfExt,
		state.OriginalOrder,
		state.InitResult.TegnsettByID,
		state.InitResult.TegnByID,
		tegnsettEnabledIDsMap,
		remainedFeaturesWithCandidate,
	)

	availability := tb.GetTegnsettsAvailability(
		state.OSInfExt,
		state.OriginalOrder,
		state.InitResult.TegnsettByID,
		state.InitResult.TegnByID,
		tegnsettEnabledIDsMap,
		remainedFeatures,
	)

	for _, tegnId := range mustRemainAvailable {
		av := availability[tegnId]
		avWithCandidate := availabilityWithCandidate[tegnId]
		if !av.Available && avWithCandidate.Available {
			// Check, that it is not available. And it will be available without the candidate selection
			return false, tegnId, av.Reason
		}
	}
	return true, "", ""
}

// NewUninstallApp creates a new uninstall UI application
func NewUninstallApp(
	data *UninstallAppInitData,
) *UninstallApp {
	// Sort installed Tegns for consistent order

	state := &UninstallUIState{
		InitResult: data.InitResult,
		OSInfExt: data.OSInfExt,
		InstalledCache: data.AlreadyInstalled,
		OriginalOrder: data.OriginalOrder,
		FilteredOrder: data.FilteredOrder,
		EnabledIDsMap:     data.EnabledIDsMap,
		ExitConfirmed:  false,
		SelectionTegnID: -1,
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

	helpBtn := tview.NewButton(" (F1) (?) Help ").
		SetSelectedFunc(func() { u.showHelpModal(nil) })

	confirmBtn := tview.NewButton(" (F5) (C-x) Confirm ").
		SetSelectedFunc(func() { u.showConfirmModal() })

	exitBtn := tview.NewButton(" (F10) (C-c) Exit ").
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
		case tcell.KeyF5, tcell.KeyCtrlX:
			u.showConfirmModal()
			return nil
		case tcell.KeyF10, tcell.KeyCtrlC:
			u.showExitModal()
			return nil
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

func (u *UninstallApp) updateStatusBar() {
	u.statusBar.SetText(
		"[yellow]Space[-]: Toggle | [yellow]F1[-]: Help | [yellow]F5[-]: Confirm | [yellow]F10[-]: Exit",
	)
}

func (u *UninstallApp) showTegnList() {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	for _, id := range u.State.FilteredOrder {
		tegn := u.State.InitResult.TegnByID[id]
		name := tegn.GetName()
		selected := u.State.EnabledIDsMap[id]
		canBeSelected, blocker, reason := canTegnBeSelected(u.State, id)

		status := "[red]( )[-]"
		if !canBeSelected {
			status = "[grey]---[-]"
		} else if selected {
			status = "[green](*)[-]"
		}

		// TODO: the availability

		secondaryText := fmt.Sprintf("%s %s", status, name)
		if !canBeSelected {
			secondaryText = fmt.Sprintf("%s %s [%s -- %s]", status, name, blocker, reason)
		}
		list.AddItem(secondaryText, "", 0, nil)
	}

	if sel := u.State.GetAndResetTegnSelection(); sel != -1 && sel < list.GetItemCount() {
		list.SetCurrentItem(sel)
	}

	// TODO: focus on the Tegn that was focused before

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentItem := list.GetCurrentItem()
		u.State.SelectionTegnID = currentItem

		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == ' ' {
				index := currentItem
				if index >= 0 && index < len(u.State.FilteredOrder) {
					id := u.State.FilteredOrder[index]

					canBeSelected, _, _ := canTegnBeSelected(u.State, id)
					if canBeSelected {
						u.State.EnabledIDsMap[id] = !u.State.EnabledIDsMap[id]
						u.refreshList()
					}
				}
				return nil
			}
			if event.Rune() == '?' {
				index := currentItem
				if index >= 0 && index < len(u.State.FilteredOrder) {
					id := u.State.FilteredOrder[index]
					tegn := u.State.InitResult.TegnByID[id].(tb.TegnGeneral)
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
	// removeAllPagesExceptOne(u.pages, "tegnList")
	u.updateStatusBar()
}

func (u *UninstallApp) refreshList() {
	u.showTegnList()
}

func (u *UninstallApp) showHelpModal(tegnGeneral *tb.TegnGeneral) {
	// Reuse the existing HelpModal from modals.go
	modal := NewHelpModal(u.State.InstalledCache, u.app, tegnGeneral, nil, func() {
		u.pages.RemovePage("helpModal")
	})
	u.pages.AddPage("helpModal", modal, true, true)
}

func (u *UninstallApp) showConfirmModal() {
	selected := 0
	for _, v := range u.State.EnabledIDsMap {
		if v {
			selected += 1
		}
	}

	text := fmt.Sprintf("Uninstall %d selected Tegns?\n\nThis action cannot be undone.", selected)
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
	modal.
		SetBorder(true).
		SetTitle("Confirm Uninstall").
		SetBackgroundColor(tcell.ColorDarkRed)
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
	modal.
		SetBorder(true).
		SetTitle("Exit")
	u.pages.AddPage("exitModal", modal, true, true)
}

// Run starts the uninstall UI
func (u *UninstallApp) Run() error {
	u.showTegnList()
	return u.app.Run()
}
