package structs

import (
	"github.com/go-delve/delve/service/api"
)

// DebuggerEvent represents types of debugger state changes
type DebuggerEvent int

const (
	EventStateChanged DebuggerEvent = iota
	EventBreakpointsChanged
	EventFileChanged
	EventProcessExited
	EventProcessRestarted
	EventStackFrameChanged
	EventWatchesChanged
	EventError
	EventContinue
	EventStep
	EventStepOver
	EventStepOut
)

// BreakpointOperation represents the type of breakpoint operation
type BreakpointOperation int

const (
	BreakpointCreated BreakpointOperation = iota
	BreakpointRemoved
	BreakpointEnabled
	BreakpointDisabled
	BreakpointsAllEnabled
	BreakpointsAllDisabled
	BreakpointsAllCleared
)

// EventData holds the data associated with an event
type EventData struct {
	State               *api.DebuggerState
	File                string
	Line                int
	StackFrame          int
	Breakpoint          *api.Breakpoint
	BreakpointOperation BreakpointOperation
	Error               error
}

// EventHandler is a function that handles debugger events
type EventHandler func(event DebuggerEvent, data EventData)
