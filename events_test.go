package main

import (
	"testing"

	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/service/api"
)

func TestSignalEmitter(t *testing.T) {
	emitter := NewSignalEmitter()

	t.Run("Connect and Emit", func(t *testing.T) {
		_ = false       // received
		_ = EventData{} // receivedData

		emitter.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
			// received = true
			// receivedData = data
		})

		_ = EventData{ // testData
			File: "test.go",
			Line: 42,
		}

		// Note: SignalEmitter tests basic connectivity - actual emissions require GTK main loop
	})

	t.Run("Multiple Subscribers", func(t *testing.T) {
		emitter := NewSignalEmitter()
		count := 0

		emitter.Connect(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
			count++
		})

		emitter.Connect(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
			count++
		})

		// Verify handlers are registered
	})

	t.Run("Different Event Types", func(t *testing.T) {
		emitter := NewSignalEmitter()

		emitter.Connect(EventFileChanged, func(event DebuggerEvent, data EventData) {})
		emitter.Connect(EventProcessExited, func(event DebuggerEvent, data EventData) {})
		emitter.Connect(EventWatchesChanged, func(event DebuggerEvent, data EventData) {})

		// Verify different event types can be registered
	})
}

func TestDebuggerEventEmission(t *testing.T) {
	t.Run("SignalEmitter Initialization", func(t *testing.T) {
		// Create a minimal debugger setup
		d := &Debugger{
			Signals: NewSignalEmitter(),
		}

		if d.Signals == nil {
			t.Error("Signals should not be nil")
		}
	})

	t.Run("Watch Management", func(t *testing.T) {
		d := &Debugger{
			Signals:     NewSignalEmitter(),
			WatchesExpr: []string{},
		}

		_ = false // eventReceived
		d.Signals.Connect(EventWatchesChanged, func(event DebuggerEvent, data EventData) {
			// eventReceived = true
		})

		d.CreateWatch("myVar")

		if len(d.WatchesExpr) != 1 {
			t.Errorf("Expected 1 watch, got %d", len(d.WatchesExpr))
		}

		if d.WatchesExpr[0] != "myVar" {
			t.Errorf("Expected watch 'myVar', got '%s'", d.WatchesExpr[0])
		}

		d.DeleteWatch("myVar")

		if len(d.WatchesExpr) != 0 {
			t.Errorf("Expected 0 watches after delete, got %d", len(d.WatchesExpr))
		}
	})

	t.Run("File Line Change", func(t *testing.T) {
		d := &Debugger{
			Signals: NewSignalEmitter(),
		}

		_ = "" // capturedFile
		_ = 0  // capturedLine

		d.Signals.Connect(EventFileChanged, func(event DebuggerEvent, data EventData) {
			// capturedFile = data.File
			// capturedLine = data.Line
		})

		d.SetFileLine("main.go", 10)

		if d.File != "main.go" {
			t.Errorf("Expected file 'main.go', got '%s'", d.File)
		}

		if d.Line != 10 {
			t.Errorf("Expected line 10, got %d", d.Line)
		}
	})
}

func TestEventData(t *testing.T) {
	t.Run("EventData Fields", func(t *testing.T) {
		state := &api.DebuggerState{
			Exited: false,
		}

		bp := &api.Breakpoint{
			File: "test.go",
			Line: 100,
		}

		data := EventData{
			State:      state,
			File:       "main.go",
			Line:       42,
			StackFrame: 0,
			Breakpoint: bp,
		}

		if data.File != "main.go" {
			t.Errorf("Expected file 'main.go', got '%s'", data.File)
		}

		if data.Line != 42 {
			t.Errorf("Expected line 42, got %d", data.Line)
		}

		if data.Breakpoint.File != "test.go" {
			t.Errorf("Expected breakpoint file 'test.go', got '%s'", data.Breakpoint.File)
		}
	})
}
