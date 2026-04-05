package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	. "github.com/emad-elsaid/debugger/structs"
	. "github.com/emad-elsaid/types"
	"github.com/go-delve/delve/pkg/gobuild"
	"github.com/go-delve/delve/pkg/proc"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
)

type Debugger struct {
	*debugger.Debugger
	State     DebuggerState
	LastState *api.DebuggerState
	EventBus  *UIEventBus

	// Project props
	Path    string
	BinName string
	Args    []string
	Test    bool

	// Run state
	StackFrame int
	File       string
	Line       int

	// Cached items
	breakpoints []*api.Breakpoint

	// Watches
	WatchesExpr []string
}

func NewDebugger(bin string, args []string, runImmediately bool, test bool) (*Debugger, error) {
	wd, _ := os.Getwd()

	// If a directory is specified in args, use it as the working directory
	if len(args) > 0 && args[0] != "." {
		dir := args[0]
		// Expand ~ to home directory
		if strings.HasPrefix(dir, "~/") {
			home, _ := os.UserHomeDir()
			dir = filepath.Join(home, dir[2:])
		}
		// Convert to absolute path
		if absPath, err := filepath.Abs(dir); err == nil {
			wd = absPath
		}
	}

	d := Debugger{
		Path:        wd,
		BinName:     bin,
		Args:        args,
		breakpoints: []*api.Breakpoint{},
		Test:        test,
		EventBus:    NewUIEventBus(),
	}

	if err := d.InitDebugger(); err != nil {
		return nil, err
	}

	// Don't start running immediately - let the UI load panels first
	// The UI will call Continue() after panels are populated
	_ = runImmediately // Will be used by UI

	return &d, nil
}

func (d *Debugger) WithTarget(fn func(*proc.Target)) {
	tg, unlock := d.Debugger.LockTargetGroup()
	defer unlock()
	fn(tg.Selected)
}

func (d *Debugger) compileArgs() []string {
	var a Slice[string] = d.Args
	sep := a.Index("--")

	var args []string
	if sep == -1 {
		args = d.Args
	} else if sep == 0 {
		return []string{}
	} else {
		args = d.Args[:sep]
	}

	// Skip the first argument if it's a directory path (already used as d.Path)
	if len(args) > 0 && args[0] != "." {
		return args[1:]
	}

	return args
}

func (d *Debugger) runArgs() []string {
	var a Slice[string] = d.Args
	sep := a.Index("--")

	if sep == -1 {
		return []string{}
	}

	return d.Args[sep+1:]
}

func (d *Debugger) InitDebugger() error {
	config := d.DebugConfig()

	if err := d.Compile(); err != nil {
		return err
	}

	binPath := filepath.Join(d.Path, d.BinName)
	deb, err := debugger.New(config, append([]string{binPath}, d.runArgs()...))
	if err != nil {
		return err
	}

	d.Debugger = deb
	d.cacheBreakpoints()

	return nil
}

func (d *Debugger) DebugConfig() *debugger.Config {
	executeKind := debugger.ExecutingGeneratedFile

	if d.Test {
		executeKind = debugger.ExecutingGeneratedTest
	}

	return &debugger.Config{
		WorkingDir:  d.Path,
		Backend:     "default",
		Foreground:  true,
		ExecuteKind: executeKind,
		Packages:    []string{},
		BuildFlags:  strings.Join(d.compileArgs(), " "),
	}
}

func (d *Debugger) Stacktrace() ([]proc.Stackframe, error) {
	if d.LastState != nil && d.LastState.Exited {
		return nil, ErrProcessExited
	}

	var frames []proc.Stackframe
	d.WithTarget(func(target *proc.Target) {
		// Get selected goroutine directly within the lock
		g := target.SelectedGoroutine()
		if g == nil {
			routines, _, _ := proc.GoroutinesInfo(target, 0, 0)
			for _, g = range routines {
				break
			}
		}
		if g != nil {
			frames, _ = proc.GoroutineStacktrace(target, g, STACK_LIMIT, 0)
		}
	})
	return frames, nil
}

func (d *Debugger) CreateBreakpoint(bp *api.Breakpoint) (*api.Breakpoint, error) {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	b, err := d.Debugger.CreateBreakpoint(bp, "", nil, false)
	d.cacheBreakpoints()
	d.EventBus.Emit(EventBreakpointsChanged, EventData{
		Breakpoint:          b,
		BreakpointOperation: BreakpointCreated,
	})
	return b, err
}

func (d *Debugger) AmendBreakpoint(bp *api.Breakpoint) error {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	// Determine operation type based on the disabled state
	operation := BreakpointEnabled
	if bp.Disabled {
		operation = BreakpointDisabled
	}

	err := d.Debugger.AmendBreakpoint(bp)
	d.cacheBreakpoints()
	d.EventBus.Emit(EventBreakpointsChanged, EventData{
		Breakpoint:          bp,
		BreakpointOperation: operation,
	})
	return err
}

func (d *Debugger) ClearBreakpoint(bp *api.Breakpoint) (*api.Breakpoint, error) {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	b, err := d.Debugger.ClearBreakpoint(bp)
	d.cacheBreakpoints()
	d.EventBus.Emit(EventBreakpointsChanged, EventData{
		Breakpoint:          b,
		BreakpointOperation: BreakpointRemoved,
	})
	return b, err
}

func (d *Debugger) ClearAllBreakpoints() error {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	for _, bp := range d.Breakpoints() {
		if _, err := d.Debugger.ClearBreakpoint(bp); err != nil {
			return err
		}
	}

	d.cacheBreakpoints()
	d.EventBus.Emit(EventBreakpointsChanged, EventData{
		BreakpointOperation: BreakpointsAllCleared,
	})

	return nil
}

func (d *Debugger) ToggleBreakpoint(bp *api.Breakpoint) error {
	bp.Disabled = !bp.Disabled
	return d.AmendBreakpoint(bp)
}

func (d *Debugger) EnableAllBreakpoints() error {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	modified := false
	for _, bp := range d.Breakpoints() {
		if !bp.Disabled {
			continue
		}
		bp.Disabled = false
		if err := d.Debugger.AmendBreakpoint(bp); err != nil {
			return err
		}
		modified = true
	}

	if modified {
		d.cacheBreakpoints()
		d.EventBus.Emit(EventBreakpointsChanged, EventData{BreakpointOperation: BreakpointsAllEnabled})
	}

	return nil
}

func (d *Debugger) DisableAllBreakpoints() error {
	if d.IsRunning() {
		d.Stop()
		defer d.Continue()
	}

	modified := false
	for _, bp := range d.Breakpoints() {
		if bp.Disabled {
			continue
		}
		bp.Disabled = true
		if err := d.Debugger.AmendBreakpoint(bp); err != nil {
			return err
		}
		modified = true
	}

	if modified {
		d.cacheBreakpoints()
		d.EventBus.Emit(EventBreakpointsChanged, EventData{BreakpointOperation: BreakpointsAllDisabled})
	}

	return nil
}

func (d *Debugger) ExecuteExpr(expr string) (*proc.Variable, error) {
	var result *proc.Variable
	d.WithTarget(func(target *proc.Target) {
		scope, err := proc.ConvertEvalScope(target, -1, d.StackFrame, 0)
		if err != nil {
			return
		}
		result, _ = scope.EvalExpression(expr, proc.LoadConfig{
			FollowPointers:     true,
			MaxVariableRecurse: 1,
			MaxStringLen:       64,
			MaxArrayValues:     64,
			MaxStructFields:    -1,
		})
	})
	return result, nil
}

func (d *Debugger) SetFileLine(file string, line int) {
	d.File = file
	d.Line = line
	d.InvalidateState()
	d.EventBus.Emit(EventFileChanged, EventData{File: file, Line: line})
}

func (d *Debugger) InvalidateState() {
	d.State++
}

func (d *Debugger) Command(cmd *api.DebuggerCommand) {
	state, _ := d.Debugger.Command(cmd, nil, nil, nil)
	if state == nil {
		return
	}

	d.LastState = state
	d.InvalidateState()

	if state.NextInProgress {
		if err := d.CancelNext(); err != nil {
			log.Printf("Error cancelling next: %s", err)
		}
	}

	for _, thread := range state.Threads {
		if thread.Breakpoint != nil {
			file := thread.Breakpoint.File
			line := thread.Breakpoint.Line
			d.SetFileLine(file, line)
			break
		}
	}

	if cmd.Name == api.Step ||
		cmd.Name == api.StepOut ||
		cmd.Name == api.Next {
		g, err := d.SelectedGoRoutine()
		if err != nil {
			return
		}

		file := g.CurrentLoc.File
		line := g.CurrentLoc.Line
		d.SetFileLine(file, line)
	}

	// Emit state changed event
	d.EventBus.Emit(EventStateChanged, EventData{State: state})

	// If the process exited load it again into the debugger.
	// this allow the user to manipulate the breakpoints after exiting the program
	if state.Exited {
		d.Debugger.Restart(false, "", false, []string{}, [3]string{}, false)
		d.EventBus.Emit(EventProcessExited, EventData{State: state})
	}
}

func (d *Debugger) SelectedGoRoutine() (*proc.G, error) {
	var g *proc.G
	d.WithTarget(func(target *proc.Target) {
		g = target.SelectedGoroutine()
		if g == nil {
			routines, _, _ := proc.GoroutinesInfo(target, 0, 0)
			for _, g = range routines {
				break
			}
		}
	})
	return g, nil
}

func (d *Debugger) GoRoutines() ([]*proc.G, error) {
	var routines []*proc.G
	d.WithTarget(func(target *proc.Target) {
		routines, _, _ = proc.GoroutinesInfo(target, 0, 0)
	})
	return routines, nil
}

func (d *Debugger) Continue() {
	if d.IsRunning() {
		return
	}

	d.EventBus.Emit(EventContinue, EventData{})
	go d.Command(&api.DebuggerCommand{Name: api.Continue})
}

func (d *Debugger) Stop() {
	d.Command(&api.DebuggerCommand{Name: api.Halt})
}

func (d *Debugger) Restart() error {
	d.Stop()

	// First try rebuilding manually using our Compile method
	if err := d.Compile(); err != nil {
		log.Printf("Error compiling during restart: %s", err)
		d.EventBus.Emit(EventError, EventData{Error: err})
		return err
	}

	// Now restart without rebuild (rebuild=false) since we already compiled
	_, err := d.Debugger.Restart(false, "", false, []string{}, [3]string{}, false)
	if err != nil {
		log.Printf("Error restarting process: %s", err)
		d.EventBus.Emit(EventError, EventData{Error: err})
		return err
	}

	d.EventBus.Emit(EventProcessRestarted, EventData{})
	go d.Continue()
	return nil
}

func (d *Debugger) StepOut() {
	if d.IsRunning() {
		return
	}

	d.EventBus.Emit(EventStepOut, EventData{})
	go d.Command(&api.DebuggerCommand{Name: api.StepOut})
}

func (d *Debugger) Next() {
	if d.IsRunning() {
		return
	}

	d.EventBus.Emit(EventStepOver, EventData{})
	go d.Command(&api.DebuggerCommand{Name: api.Next})
}

func (d *Debugger) Step() {
	if d.IsRunning() {
		return
	}

	d.EventBus.Emit(EventStep, EventData{})
	go d.Command(&api.DebuggerCommand{Name: api.Step})
}

func (d *Debugger) JumpToFunction(name string) {
	d.WithTarget(func(target *proc.Target) {
		bi := target.BinInfo()
		pcs, err := proc.FindFunctionLocation(target.Process, name, 0)
		if err != nil {
			slog.Debug("JumpToFunction: error finding function", "name", name, "error", err)
			return
		}

		for _, pc := range pcs {
			file, line, _ := bi.PCToLine(pc)
			if file != "" {
				d.SetFileLine(file, line)
				return
			}
		}
	})
}

func (d *Debugger) ShowMainFunction() {
	d.JumpToFunction("main.main")
}

func (d *Debugger) Compile() error {
	config := d.DebugConfig()

	// Change to target directory for compilation
	originalDir, _ := os.Getwd()
	if err := os.Chdir(d.Path); err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	// Build binary path relative to target directory
	binPath := filepath.Join(d.Path, d.BinName)

	if _, err := os.Stat(binPath); err == nil {
		gobuild.Remove(binPath)
	}

	var err error
	if d.Test {
		err = gobuild.GoTestBuild(binPath, config.Packages, config.BuildFlags)
	} else {
		err = gobuild.GoBuild(binPath, config.Packages, config.BuildFlags)
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *Debugger) Breakpoints() []*api.Breakpoint {
	return d.breakpoints
}

func (d *Debugger) cacheBreakpoints() {
	d.breakpoints = []*api.Breakpoint{}
	d.WithTarget(func(target *proc.Target) {
		for _, lbp := range target.Breakpoints().Logical {
			abp := api.ConvertLogicalBreakpoint(lbp)
			d.breakpoints = append(d.breakpoints, abp)
		}
	})
}

func (d *Debugger) FunctionArguments() ([]*proc.Variable, error) {
	var vars []*proc.Variable
	d.WithTarget(func(target *proc.Target) {
		s, err := proc.ConvertEvalScope(target, -1, d.StackFrame, 0)
		if err != nil {
			return
		}
		vars, _ = s.FunctionArguments(ProcLoadConfig)
	})
	return vars, nil
}

func (d *Debugger) LocalVariables() ([]*proc.Variable, error) {
	var vars []*proc.Variable
	d.WithTarget(func(target *proc.Target) {
		s, err := proc.ConvertEvalScope(target, -1, d.StackFrame, 0)
		if err != nil {
			return
		}
		vars, _ = s.LocalVariables(ProcLoadConfig)
	})
	return vars, nil
}

func (d *Debugger) PackageVariables() ([]*proc.Variable, error) {
	var vars []*proc.Variable
	d.WithTarget(func(target *proc.Target) {
		s, err := proc.ConvertEvalScope(target, -1, d.StackFrame, 0)
		if err != nil {
			return
		}
		vars, _ = s.PackageVariables(ProcLoadConfig)
	})
	return vars, nil
}

func (d *Debugger) RunningLines() []int {
	lines := []int{}

	d.WithTarget(func(target *proc.Target) {
		if ok, _ := target.Valid(); ok {
			threads := target.ThreadList()
			for _, thread := range threads {
				loc, err := proc.ThreadLocation(thread)
				if err == nil && loc.File == d.File {
					lines = append(lines, loc.Line)
				}
			}
		}
	})

	return lines
}

func (d *Debugger) CreateWatch(expr string) {
	d.WatchesExpr = append(d.WatchesExpr, expr)
	d.EventBus.Emit(EventWatchesChanged, EventData{})
}

func (d *Debugger) DeleteWatch(expr string) {
	nWatches := []string{}
	for _, v := range d.WatchesExpr {
		if v != expr {
			nWatches = append(nWatches, v)
		}
	}
	d.WatchesExpr = nWatches
	d.EventBus.Emit(EventWatchesChanged, EventData{})
}
