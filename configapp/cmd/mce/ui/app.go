package ui

import (
	// "slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	// "github.com/toliak/mce/inspector"
	tb "github.com/toliak/mce/tegnbuilder"
)

// App represents the main TUI application
type App struct {
	app        *tview.Application
	pages      *tview.Pages
	State      *UIState
	rootView   *tview.Flex
	statusBar  *tview.TextView
}

// NewApp creates a new TUI application
func NewApp(
	initResult tb.TegnsettInitializeResult, 
	osInfoExt tb.OSInfoExt, 
	alreadyInstalled tb.AvailablePackagesMap,
	alreadyInstalledFeatures tb.TegnInstalledFeaturesMap,
) App {
	app := tview.NewApplication()
	state := NewUIState(initResult, osInfoExt, alreadyInstalled)
	state.InstalledFeatures = alreadyInstalledFeatures

	ui := App{
		app:       app,
		State:     state,
	}
	
	ui.setupUI()
	return ui
}

func (a *App) setupUI() {
	// Create main layout
	a.rootView = tview.NewFlex().SetDirection(tview.FlexRow)
	
	// Create title
	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[::b]Configuration[::]").
		SetDynamicColors(true)
	
	// Create main content area
	a.pages = tview.NewPages()
	
	// Create status bar
	a.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	a.updateStatusBar()
	
	// Create button bar
	buttonBar := a.createButtonBar()
	
	// Assemble the layout
	a.rootView.AddItem(title, 1, 0, false)
	a.rootView.AddItem(a.pages, 0, 1, true)
	a.rootView.AddItem(a.statusBar, 1, 0, false)
	a.rootView.AddItem(buttonBar, 1, 0, false)
	
	a.app.SetRoot(a.rootView, true)
	a.setupGlobalKeyBindings()
}

func (a *App) createButtonBar() *tview.Flex {
	buttonBar := tview.NewFlex().SetDirection(tview.FlexRow)
	
	buttons := tview.NewFlex().SetDirection(tview.FlexColumn)
	
	helpBtn := tview.NewButton(" (F1) (?) Help ").
		SetSelectedFunc(func() { a.showHelpModal(nil, nil) })
	
	// searchBtn := tview.NewButton(" (F3) Search ").
	// 	SetSelectedFunc(func() { a.showSearchModal() })
	
	confirmBtn := tview.NewButton(" (F5) (C-x) Confirm ").
		SetSelectedFunc(func() { a.showConfirmModal() })
	
	exitBtn := tview.NewButton(" (F10) (C-c) Exit ").
		SetSelectedFunc(func() { a.showExitModal() })
	
	buttons.AddItem(helpBtn, 0, 1, false)
	// buttons.AddItem(searchBtn, 0, 1, false)
	buttons.AddItem(confirmBtn, 0, 1, false)
	buttons.AddItem(exitBtn, 0, 1, false)
	
	buttonBar.AddItem(buttons, 0, 1, true)
	return buttonBar
}

func (a *App) setupGlobalKeyBindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			a.showHelpModal(nil, nil)
			return nil
		// case tcell.KeyF3:
		// 	a.showSearchModal()
		// 	return nil
		case tcell.KeyF5, tcell.KeyCtrlX:
			a.showConfirmModal()
			return nil
		case tcell.KeyF10, tcell.KeyCtrlC:
			// names := a.pages.GetPageNames(true)
			// if slices.Contains(names, "tegnsettList") || slices.Contains(names, "tegnList") || slices.Contains(names, "parameterList") {
			// 	// TODO: wtf, thats baaad
			// 	return event
			// }
			a.showExitModal()
			return nil
		case tcell.KeyRune:
			// if event.Rune() == '?' {
			// 	a.showHelpModal(nil)
			// 	return nil
			// }
			// if event.Rune() == '/' {
			// 	a.showSearchModal()
			// 	return nil
			// }
		}
		return event
	})
}

func (a *App) updateStatusBar() {
	currentView := a.State.CurrentView
	statusText := ""
	
	switch currentView {
	case ViewTegnsettList:
		statusText = "View: Tegnsetts | [yellow]Space[white]: Toggle | [yellow]Enter[white]: Open | [yellow]F1[white]: Help"
	case ViewTegnList:
		statusText = "View: Tegns | [yellow]Space[white]: Toggle | [yellow]Enter[white]: Properties | [yellow]F1[white]: Help"
	case ViewParameterList:
		statusText = "View: Parameters | [yellow]Enter[white]: Edit | [yellow]Esc[white]: Back"
	}
	
	a.statusBar.SetText(statusText)
}

func (a *App) showTegnsettList() {
	order, err := tb.GetTegnsettsOrder(a.State.InitResult.TegnsettByID)
	if err != nil {
		a.State.ExitError = err
		a.app.Stop()
		return
	}

	// TODO: pick already installed 

	availability := tb.GetTegnsettsAvailability(
		a.State.OSInfExt,
		*order,
		a.State.InitResult.TegnsettByID,
		a.State.InitResult.TegnByID,
		a.State.EnabledIDsMap,
		a.State.InstalledFeatures,
	)

	list := NewTegnsettList(a.State, order.Tegnsett, availability, a)
	a.pages.AddPage("tegnsettList", list, true, true)
	removeAllPagesExceptOne(a.pages, "tegnsettList")
	a.State.CurrentView = ViewTegnsettList
	a.updateStatusBar()
}

func (a *App) showTegnList(tegnsettID string) {
	order, err := tb.GetTegnsettsOrder(a.State.InitResult.TegnsettByID)
	if err != nil {
		a.State.ExitError = err
		a.app.Stop()
		return
	}

	availability := tb.GetTegnsettsAvailability(
		a.State.OSInfExt,
		*order,
		a.State.InitResult.TegnsettByID,
		a.State.InitResult.TegnByID,
		a.State.EnabledIDsMap,
		a.State.InstalledFeatures,
	)

	list := NewTegnList(a.State, tegnsettID, order.TegnByTegnsettID[tegnsettID], availability, a)
	a.pages.AddPage("tegnList", list, true, true)
	removeAllPagesExceptOne(a.pages, "tegnList")
	a.State.CurrentView = ViewTegnList
	a.State.CurrentTegnsettID = tegnsettID
	a.updateStatusBar()
}

func (a *App) showParameterList(tegnID string) {
	list := NewParameterList(a.State, tegnID, a)
	a.pages.AddPage("parameterList", list, true, true)
	removeAllPagesExceptOne(a.pages, "parameterList")
	a.State.CurrentView = ViewParameterList
	a.State.CurrentTegnID = tegnID
	a.updateStatusBar()
}

func (a *App) showHelpModal(tegnGeneral *tb.TegnGeneral, availability *tb.TegnAvailability) {
	modal := NewHelpModal(a.State.InstalledCache, a.app, tegnGeneral, availability, func () {a.closeModals()})
	a.pages.AddPage("helpModal", modal, true, true)
}

func (a *App) showSearchModal() {
	modal := NewSearchModal(a.State, a.app, func () {a.closeModals()})
	a.pages.AddPage("searchModal", modal, true, true)
}

func (a *App) showConfirmModal() {
	modal := NewConfirmModal(a.State, a.app, func () {a.closeModals()})
	a.pages.AddPage("confirmModal", modal, true, true)
}

func (a *App) showExitModal() {
	modal := NewExitModal(a.State, a.app, func () {a.closeModals()})
	a.pages.AddPage("exitModal", modal, true, true)
}

func (a *App) navigateBack() {
	switch a.State.CurrentView {
	case ViewParameterList:
		a.showTegnList(a.State.CurrentTegnsettID)
	case ViewTegnList:
		a.showTegnsettList()
	}
	a.updateStatusBar()
}

func removeAllPagesExceptOne(pages *tview.Pages, nameToLeave string) {
	allNames := pages.GetPageNames(false)
	namesToDelete := make([]string, 0, len(allNames))
	for _, name := range namesToDelete {
		if name == nameToLeave {
			continue
		}

		namesToDelete = append(namesToDelete, name)
	}
	
	for _, name := range namesToDelete {
		pages.RemovePage(name)
	}
	pages.ShowPage(nameToLeave)
	pages.SwitchToPage(nameToLeave)
}


func (a *App) closeModals() {
	switch a.State.CurrentView {
	case ViewTegnsettList:
		removeAllPagesExceptOne(a.pages, "tegnsettList")
	case ViewTegnList:
		removeAllPagesExceptOne(a.pages, "tegnList")
	case ViewParameterList:
		removeAllPagesExceptOne(a.pages, "parameterList")
	}
}

// Run starts the application
func (a *App) Run() error {
	a.showTegnsettList()
	if a.State.ExitError != nil {
		return a.State.ExitError
	}
	return a.app.Run()
}