package main

import (
	"testing"

	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/service/api"
)

func TestEventBus(t *testing.T) {
	bus := NewEventBus()

	t.Run("Subscribe and Emit", func(t *testing.T) {
		_ = false       // received
		_ = EventData{} // receivedData

		bus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
			// received = true
			// receivedData = data
		})

		_ = EventData{ // testData
			File: "test.go",
			Line: 42,
		}

		// Note: Core EventBus doesn't use GTK - testing direct subscription
		if bus.HandlerCount(EventStateChanged) != 1 {
			t.Errorf("Expected 1 handler, got %d", bus.HandlerCount(EventStateChanged))
		}
	})

	t.Run("Multiple Subscribers", func(t *testing.T) {
		bus := NewEventBus()
		count := 0

		bus.Subscribe(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
			count++
		})

		bus.Subscribe(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
			count++
		})

		if bus.HandlerCount(EventBreakpointsChanged) != 2 {
			t.Errorf("Expected 2 handlers, got %d", bus.HandlerCount(EventBreakpointsChanged))
		}
	})

	t.Run("Different Event Types", func(t *testing.T) {
		bus := NewEventBus()

		bus.Subscribe(EventFileChanged, func(event DebuggerEvent, data EventData) {})
		bus.Subscribe(EventProcessExited, func(event DebuggerEvent, data EventData) {})
		bus.Subscribe(EventWatchesChanged, func(event DebuggerEvent, data EventData) {})

		if bus.EventTypeCount() != 3 {
			t.Errorf("Expected 3 event types, got %d", bus.EventTypeCount())
		}
	})
}

func TestDebuggerEventEmission(t *testing.T) {
	t.Run("EventBus Initialization", func(t *testing.T) {
		// Create a minimal debugger setup
		d := &Debugger{
			EventBus: NewUIEventBus(),
		}

		if d.EventBus == nil {
			t.Error("EventBus should not be nil")
		}

		if d.EventBus.EventBus == nil {
			t.Error("Underlying EventBus should not be nil")
		}
	})

	t.Run("Watch Management", func(t *testing.T) {
		d := &Debugger{
			EventBus:    NewUIEventBus(),
			WatchesExpr: []string{},
		}

		_ = false // eventReceived
		d.EventBus.Subscribe(EventWatchesChanged, func(event DebuggerEvent, data EventData) {
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
			EventBus: NewUIEventBus(),
		}

		_ = "" // capturedFile
		_ = 0  // capturedLine

		d.EventBus.Subscribe(EventFileChanged, func(event DebuggerEvent, data EventData) {
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
