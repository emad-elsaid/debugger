package main

import (
	"errors"
	"path"
	"time"

	. "github.com/emad-elsaid/debugger/ui"
	"github.com/fsnotify/fsnotify"
)

const STACK_LIMIT = 100000

var (
	ErrCantInspectTargetProcess = errors.New("Can't inspect target process")
	ErrProcessExited            = errors.New("Process exited")
	ErrNoSelectedGoRoutine      = errors.New("Can't find any go routines")
	ErrCantGetScope             = errors.New("Can't get current scope")
)

type DebugScreen struct {
	Watch            *fsnotify.Watcher
	Analytics        Analytics
	Debugger         *Debugger
	BottomTabs       Tabs
	Toolbar          Toolbar
	SideBar          SideBar
	FunctionsPanel   FunctionsPanel
	PackagesPanel    PackagesPanel
	SourcesPanel     SourcesPanel
	StackPanel       StackPanel
	TypesPanel       TypesPanel
	ConsolePanel     ConsolePanel
	BreakpointsPanel BreakpointsPanel
	DisassemblyPanel DisassemblyPanel
	SourcePanel      SourcePanel
	OpenFilesPanel   OpenFilesPanel
	MemoryPanel      MemoryPanel
}

func (d *DebugScreen) Layout(c C) D {
	debugger := d.Debugger

	if debugger.File == "" {
		debugger.JumpToFunction("main.main")
	}

	return Rows(
		Rigid(d.Toolbar.Layout(debugger)),
		Rigid(HR(1)),
		Flexed(2,
			Columns(
				Flexed(4, d.SourcePanel.Layout(debugger)),
				Flexed(1, d.SideBar.Layout(debugger)),
			),
		),
		Flexed(1,
			d.BottomTabs.Layout(
				&TabChild{Name: "Console", Panel: func(c C) D { return d.ConsolePanel.Layout(debugger)(c) }},
				&TabChild{Name: "Stack Trace", Panel: func(c C) D { return d.StackPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Types", Panel: func(c C) D { return d.TypesPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Functions", Panel: func(c C) D { return d.FunctionsPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Packages", Panel: func(c C) D { return d.PackagesPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Sources", Panel: func(c C) D { return d.SourcesPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Disassembly", Panel: func(c C) D { return d.DisassemblyPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Breakpoints", Panel: func(c C) D { return d.BreakpointsPanel.Layout(debugger)(c) }},
				&TabChild{Name: "Open Files", Panel: func(c C) D { return d.OpenFilesPanel.Layout(debugger, &d.Analytics)(c) }},
				&TabChild{Name: "Memory", Panel: func(c C) D { return d.MemoryPanel.Layout(&d.Analytics)(c) }},
			),
		),
	)(c)
}

func NewDebugScreen(debugger *Debugger) (*DebugScreen, error) {
	w := &DebugScreen{
		Debugger:         debugger,
		Analytics:        NewAnalytics(),
		BottomTabs:       Tabs{},
		Toolbar:          Toolbar{},
		SideBar:          NewSideBar(),
		FunctionsPanel:   NewFunctionsPanel(),
		PackagesPanel:    NewPackagesPanel(),
		SourcesPanel:     NewSourcesPanel(),
		StackPanel:       NewStackPanel(),
		TypesPanel:       NewTypesPanel(),
		ConsolePanel:     NewConsolePanel(),
		BreakpointsPanel: NewBreakpointsPanel(),
		DisassemblyPanel: NewDisassemblyPanel(),
		SourcePanel:      NewSourcePanel(),
		OpenFilesPanel:   NewOpenFilesPanel(),
		MemoryPanel:      NewMemoryPanel(),
	}

	go w.Clock()
	go w.StartWatch()

	return w, nil
}

func (d *DebugScreen) Log(kind LogKind, log string, params ...interface{}) {
	d.ConsolePanel.Log(kind, log, params...)
}

func (d *DebugScreen) StopWatch() {
	if d.Watch != nil {
		d.Watch.Close()
		d.Watch = nil
	}
}

func (d *DebugScreen) StartWatch() {
	debugger := d.Debugger
	binpath := path.Join(d.Debugger.Path, d.Debugger.BinName)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		d.Log(LogError, "Can't create watch. Error: %s", err)
		return
	}
	defer watcher.Close()

	d.Watch = watcher

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				isWrite := event.Op == fsnotify.Write
				isRename := event.Op == fsnotify.Rename
				isRemove := event.Op == fsnotify.Remove

				supportedOp := isWrite || isRename || isRemove

				if supportedOp && event.Name != binpath {
					d.Log(LogInfo, "%s file: %s", event.Op, event.Name)
					debugger.Restart()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				d.Log(LogError, "Watch error: %s", err)
			}
		}
	}()

	if err := WatchAddRecursive(watcher, debugger.ProjectPath); err != nil {
		d.Log(LogError, "Can't watch: %s Error: %s", debugger.ProjectPath, err)
	}

	<-done
}

func (d *DebugScreen) Clock() {
	for {
		if d.Debugger.IsRunning() {
			d.Analytics.Collect(d.Debugger.Target().Pid())
		}
		time.Sleep(time.Millisecond * 10)
	}
}
