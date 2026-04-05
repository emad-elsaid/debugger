package main

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	. "github.com/emad-elsaid/debugger/structs"
)

// UIEventBus wraps the core EventBus and ensures handlers run on GTK main thread
type UIEventBus struct {
	*EventBus
}

// NewUIEventBus creates a new UI-aware event bus
func NewUIEventBus() *UIEventBus {
	return &UIEventBus{
		EventBus: NewEventBus(),
	}
}

// Emit triggers all handlers on the GTK main thread
func (ueb *UIEventBus) Emit(event DebuggerEvent, data EventData) {
	// Schedule emission on GTK main thread
	glib.IdleAdd(func() {
		ueb.EventBus.Emit(event, data)
	})
}
