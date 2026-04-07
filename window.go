package main

import (
	. "github.com/emad-elsaid/debugger/structs"
	"log/slog"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/pkg/proc"
)

type DebugWindow struct {
	*adw.ApplicationWindow
	app      *adw.Application
	debugger *Debugger

	// Widgets
	headerBar  *adw.HeaderBar
	sourceView *SourceView
	sidebar    *Sidebar
	mainPane   *gtk.Paned
	bottomPane *gtk.Paned
	tabView    *adw.TabView
	tabBar     *adw.TabBar

	// Panels
	consolePanel     *ConsolePanel
	stackPanel       *StackPanel
	typesPanel       *TypesPanel
	functionsPanel   *FunctionsPanel
	packagesPanel    *PackagesPanel
	sourcesPanel     *SourcesPanel
	disassemblyPanel *DisassemblyPanel
	breakpointsPanel *BreakpointsPanel
	openFilesPanel   *OpenFilesPanel
	memoryPanel      *MemoryPanel

	// Toolbar buttons
	continueBtn *gtk.Button
	pauseBtn    *gtk.Button
	restartBtn  *gtk.Button
	stepBtn     *gtk.Button
	nextBtn     *gtk.Button
	stepOutBtn  *gtk.Button
}

func NewDebugWindow(app *adw.Application, debugger *Debugger) *DebugWindow {
	slog.Debug("NewDebugWindow: loading UI from XML")
	builder := gtk.NewBuilderFromFile("ui/window.ui")

	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromPath("ui/style.css")
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	slog.Debug("NewDebugWindow: getting window object")
	obj := builder.GetObject("main_window")
	appWin := obj.Cast().(*adw.ApplicationWindow)

	slog.Debug("NewDebugWindow: setting window title")
	title := "Debugger"
	if debugger.BinName != "" {
		title = "Debugger - " + debugger.BinName
	}
	appWin.SetTitle(title)

	slog.Debug("NewDebugWindow: setting application on window")
	appWin.SetApplication(&app.Application)

	win := &DebugWindow{
		ApplicationWindow: appWin,
		app:               app,
		debugger:          debugger,
	}

	slog.Debug("NewDebugWindow: getting UI widgets from builder")
	win.headerBar = builder.GetObject("header_bar").Cast().(*adw.HeaderBar)
	win.mainPane = builder.GetObject("main_pane").Cast().(*gtk.Paned)
	win.bottomPane = builder.GetObject("bottom_pane").Cast().(*gtk.Paned)
	win.tabView = builder.GetObject("tab_view").Cast().(*adw.TabView)
	win.tabBar = builder.GetObject("tab_bar").Cast().(*adw.TabBar)

	// Get tab pages and set closable property (cannot be set in XML)
	pages := []string{
		"console_page", "stack_page", "types_page", "functions_page",
		"packages_page", "sources_page", "disassembly_page",
		"breakpoints_page", "openfiles_page", "memory_page",
	}
	for _, pageID := range pages {
		page := builder.GetObject(pageID).Cast().(*adw.TabPage)
		page.Object.SetObjectProperty("closable", false)
	}

	// Get placeholder boxes
	consolePlaceholder := builder.GetObject("console_placeholder").Cast().(*gtk.Box)
	stackPlaceholder := builder.GetObject("stack_placeholder").Cast().(*gtk.Box)
	typesPlaceholder := builder.GetObject("types_placeholder").Cast().(*gtk.Box)
	functionsPlaceholder := builder.GetObject("functions_placeholder").Cast().(*gtk.Box)
	packagesPlaceholder := builder.GetObject("packages_placeholder").Cast().(*gtk.Box)
	sourcesPlaceholder := builder.GetObject("sources_placeholder").Cast().(*gtk.Box)
	disassemblyPlaceholder := builder.GetObject("disassembly_placeholder").Cast().(*gtk.Box)
	breakpointsPlaceholder := builder.GetObject("breakpoints_placeholder").Cast().(*gtk.Box)
	openfilesPlaceholder := builder.GetObject("openfiles_placeholder").Cast().(*gtk.Box)
	memoryPlaceholder := builder.GetObject("memory_placeholder").Cast().(*gtk.Box)

	slog.Debug("NewDebugWindow: getting toolbar buttons from builder")
	win.continueBtn = builder.GetObject("continue_btn").Cast().(*gtk.Button)
	win.pauseBtn = builder.GetObject("pause_btn").Cast().(*gtk.Button)
	win.restartBtn = builder.GetObject("restart_btn").Cast().(*gtk.Button)
	win.stepBtn = builder.GetObject("step_btn").Cast().(*gtk.Button)
	win.nextBtn = builder.GetObject("next_btn").Cast().(*gtk.Button)
	win.stepOutBtn = builder.GetObject("step_out_btn").Cast().(*gtk.Button)

	slog.Debug("NewDebugWindow: connecting button signals")
	win.connectButtonSignals()

	slog.Debug("NewDebugWindow: connecting close request handler")
	win.ConnectCloseRequest(func() bool {
		slog.Info("Window close requested")
		glib.IdleAdd(func() {
			slog.Info("Calling app.Quit()")
			win.app.Quit()
			slog.Info("app.Quit() returned")
		})
		slog.Info("Close handler returning true")
		return true
	})

	slog.Debug("NewDebugWindow: setting up layout")
	win.setupLayout(
		consolePlaceholder,
		stackPlaceholder,
		typesPlaceholder,
		functionsPlaceholder,
		packagesPlaceholder,
		sourcesPlaceholder,
		disassemblyPlaceholder,
		breakpointsPlaceholder,
		openfilesPlaceholder,
		memoryPlaceholder,
	)

	slog.Debug("NewDebugWindow: setting up shortcuts")
	win.setupShortcuts()

	slog.Debug("NewDebugWindow: subscribing to events")
	win.subscribeToEvents()

	slog.Debug("NewDebugWindow: scheduling initial panel data load")
	glib.IdleAdd(func() bool {
		slog.Debug("NewDebugWindow: loading initial panel data (async)")
		win.refreshPanels()
		return false
	})

	slog.Debug("NewDebugWindow: complete")
	return win
}

func (w *DebugWindow) refreshPanels() {
	// Manually trigger initial loads for panels that depend on debugger state
	var mainFile string
	var mainLine int

	// Gather all data in a single WithTarget call to avoid nested locks
	w.debugger.WithTarget(func(target *proc.Target) {
		// Gather sources data (no GTK)
		if w.sourcesPanel != nil {
			w.sourcesPanel.gatherSourcesData(target)
		}

		// Gather functions data (no GTK)
		if w.functionsPanel != nil {
			w.functionsPanel.gatherFunctionsData(target)
		}

		// Gather types data (no GTK)
		if w.typesPanel != nil {
			w.typesPanel.gatherTypesData(target)
		}

		// Find main function location
		bi := target.BinInfo()
		pcs, err := proc.FindFunctionLocation(target.Process, "main.main", 0)
		if err == nil && len(pcs) > 0 {
			mainFile, mainLine, _ = bi.PCToLine(pcs[0])
		}
	})

	// Update UI outside the lock
	if w.sourcesPanel != nil {
		w.sourcesPanel.updateTree()
	}

	if w.functionsPanel != nil {
		w.functionsPanel.updateList()
	}

	if w.typesPanel != nil {
		w.typesPanel.updateList()
	}

	// Load panels that don't need target access
	if w.packagesPanel != nil {
		w.packagesPanel.loadPackages()
	}

	// Show main function in source view (no lock needed - just emits event)
	if mainFile != "" {
		w.debugger.SetFileLine(mainFile, mainLine)
	}
}

func (w *DebugWindow) connectButtonSignals() {
	w.continueBtn.ConnectClicked(func() {
		w.debugger.Continue()
	})

	w.pauseBtn.ConnectClicked(func() {
		w.debugger.Stop()
	})

	w.restartBtn.ConnectClicked(func() {
		go w.debugger.Restart()
	})

	w.stepBtn.ConnectClicked(func() {
		w.debugger.Step()
	})

	w.nextBtn.ConnectClicked(func() {
		w.debugger.Next()
	})

	w.stepOutBtn.ConnectClicked(func() {
		w.debugger.StepOut()
	})
}

func (w *DebugWindow) setupLayout(
	consolePlaceholder,
	stackPlaceholder,
	typesPlaceholder,
	functionsPlaceholder,
	packagesPlaceholder,
	sourcesPlaceholder,
	disassemblyPlaceholder,
	breakpointsPlaceholder,
	openfilesPlaceholder,
	memoryPlaceholder *gtk.Box,
) {
	slog.Debug("setupLayout: creating source view")
	w.sourceView = NewSourceView(w.debugger)

	slog.Debug("setupLayout: creating sidebar")
	w.sidebar = NewSidebar(w.debugger)

	slog.Debug("setupLayout: creating console panel")
	w.consolePanel = NewConsolePanel(w.debugger)

	slog.Debug("setupLayout: creating stack panel")
	w.stackPanel = NewStackPanel(w.debugger)

	slog.Debug("setupLayout: creating types panel")
	w.typesPanel = NewTypesPanel(w.debugger)

	slog.Debug("setupLayout: creating functions panel")
	w.functionsPanel = NewFunctionsPanel(w.debugger)

	slog.Debug("setupLayout: creating packages panel")
	w.packagesPanel = NewPackagesPanel(w.debugger)

	slog.Debug("setupLayout: creating sources panel")
	w.sourcesPanel = NewSourcesPanel(w.debugger)

	slog.Debug("setupLayout: creating disassembly panel")
	w.disassemblyPanel = NewDisassemblyPanel(w.debugger)

	slog.Debug("setupLayout: creating breakpoints panel")
	w.breakpointsPanel = NewBreakpointsPanel(w.debugger)

	slog.Debug("setupLayout: creating open files panel")
	w.openFilesPanel = NewOpenFilesPanel(w.debugger)

	slog.Debug("setupLayout: creating memory panel")
	w.memoryPanel = NewMemoryPanel(w.debugger)

	slog.Debug("setupLayout: inserting source view and sidebar into main pane")
	w.mainPane.SetStartChild(w.sourceView)
	w.mainPane.SetEndChild(w.sidebar)

	// Ensure sidebar respects its minimum width (300px set in sidebar.ui)
	w.mainPane.SetShrinkEndChild(false)  // Don't shrink sidebar below its width-request
	w.mainPane.SetResizeStartChild(true) // Source view is resizable

	slog.Debug("setupLayout: populating panel placeholders defined in UI")
	consolePlaceholder.Append(w.consolePanel)
	stackPlaceholder.Append(w.stackPanel)
	typesPlaceholder.Append(w.typesPanel)
	functionsPlaceholder.Append(w.functionsPanel)
	packagesPlaceholder.Append(w.packagesPanel)
	sourcesPlaceholder.Append(w.sourcesPanel)
	disassemblyPlaceholder.Append(w.disassemblyPanel)
	breakpointsPlaceholder.Append(w.breakpointsPanel)
	openfilesPlaceholder.Append(w.openFilesPanel)
	memoryPlaceholder.Append(w.memoryPanel)
}

func (w *DebugWindow) setupShortcuts() {
	evtController := gtk.NewEventControllerKey()

	evtController.ConnectKeyPressed(func(keyval uint, keycode uint, state gdk.ModifierType) bool {
		ctrl := state.Has(gdk.ControlMask)
		shift := state.Has(gdk.ShiftMask)

		switch keyval {
		case gdk.KEY_F5:
			if ctrl && shift {
				go w.debugger.Restart()
			} else {
				w.debugger.Continue()
			}
			return true
		case gdk.KEY_F10:
			w.debugger.Next()
			return true
		case gdk.KEY_F11:
			if shift {
				w.debugger.StepOut()
			} else {
				w.debugger.Step()
			}
			return true
		}
		return false
	})

	w.AddController(evtController)
}

func (w *DebugWindow) subscribeToEvents() {
	w.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		// Update UI based on state changes
		if data.State != nil {
			// Update button sensitivity based on running state
			isRunning := w.debugger.IsRunning()
			w.continueBtn.SetSensitive(!isRunning)
			w.pauseBtn.SetSensitive(isRunning)
			w.stepBtn.SetSensitive(!isRunning)
			w.nextBtn.SetSensitive(!isRunning)
			w.stepOutBtn.SetSensitive(!isRunning)
		}
	})

	w.debugger.Signals.Connect(EventProcessExited, func(event DebuggerEvent, data EventData) {
		// Show notification or update status
	})

	w.debugger.Signals.Connect(EventFileChanged, func(event DebuggerEvent, data EventData) {
		// Update source view (will be implemented in Phase 2)
	})

	w.debugger.Signals.Connect(EventError, func(event DebuggerEvent, data EventData) {
		if data.Error != nil {
			w.showErrorDialog("Error", data.Error.Error())
		}
	})
}

func (w *DebugWindow) showErrorDialog(title, message string) {
	dialog := adw.NewMessageDialog(&w.ApplicationWindow.Window, title, message)
	dialog.AddResponse("ok", "OK")
	dialog.SetDefaultResponse("ok")
	dialog.SetCloseResponse("ok")
	dialog.Show()
}
