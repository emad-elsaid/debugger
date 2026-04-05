package structs

import (
	"sync"

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

// EventBus manages event subscriptions and emissions
type EventBus struct {
	handlers map[DebuggerEvent][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[DebuggerEvent][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (eb *EventBus) Subscribe(event DebuggerEvent, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[event] = append(eb.handlers[event], handler)
}

// Emit triggers all handlers for a specific event
// Note: Handlers are called synchronously. UI code should use glib.IdleAdd internally.
func (eb *EventBus) Emit(event DebuggerEvent, data EventData) {
	eb.mu.RLock()
	handlers := make([]EventHandler, len(eb.handlers[event]))
	copy(handlers, eb.handlers[event])
	eb.mu.RUnlock()

	for _, handler := range handlers {
		handler(event, data)
	}
}

// HandlerCount returns the number of handlers for a specific event (for testing)
func (eb *EventBus) HandlerCount(event DebuggerEvent) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.handlers[event])
}

// EventTypeCount returns the total number of event types with handlers (for testing)
func (eb *EventBus) EventTypeCount() int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.handlers)
}
