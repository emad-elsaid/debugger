package main

import (
	"sync"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	. "github.com/emad-elsaid/debugger/structs"
)

// SignalEmitter provides GTK-compatible signal-based event emission.
// All signal emissions execute on the GTK main thread, ensuring thread-safe GTK API access.
// This provides a GTK-idiomatic interface replacing the previous custom event bus.
type SignalEmitter struct {
	handlers map[DebuggerEvent][]EventHandler
	mu       sync.RWMutex
}

// NewSignalEmitter creates a new GObject-compatible signal emitter
func NewSignalEmitter() *SignalEmitter {
	return &SignalEmitter{
		handlers: make(map[DebuggerEvent][]EventHandler),
	}
}

// Connect subscribes to a specific event type.
// Handlers are invoked on the GTK main thread, allowing direct GTK API calls.
// This mirrors GTK's Object.Connect() pattern for consistency.
func (se *SignalEmitter) Connect(event DebuggerEvent, handler EventHandler) {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.handlers[event] = append(se.handlers[event], handler)
}

// Emit triggers all handlers for a specific event on the GTK main thread.
// Uses g_main_context_invoke_full so this is safe to call from any goroutine.
// Delve runs off the GTK thread, so main thread marshalling is critical.
func (se *SignalEmitter) Emit(event DebuggerEvent, data EventData) {
	glib.MainContextDefault().InvokeFull(int(glib.PriorityDefaultIdle), func() bool {
		se.mu.RLock()
		handlers := make([]EventHandler, len(se.handlers[event]))
		copy(handlers, se.handlers[event])
		se.mu.RUnlock()

		for _, handler := range handlers {
			handler(event, data)
		}
		return false
	})
}
