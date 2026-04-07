package main

import (
	"fmt"
	. "github.com/emad-elsaid/debugger/structs"
	"log/slog"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ConsolePanel struct {
	*gtk.ScrolledWindow
	debugger   *Debugger
	textView   *gtk.TextView
	textBuffer *gtk.TextBuffer
	errorTag   *gtk.TextTag
	warningTag *gtk.TextTag
}

func NewConsolePanel(debugger *Debugger) *ConsolePanel {
	slog.Debug("NewConsolePanel: loading UI from XML")
	builder := gtk.NewBuilderFromFile("ui/console.ui")

	slog.Debug("NewConsolePanel: getting widgets from builder")
	scrolledWindow := builder.GetObject("console_panel").Cast().(*gtk.ScrolledWindow)
	textView := builder.GetObject("console_text_view").Cast().(*gtk.TextView)

	panel := &ConsolePanel{
		ScrolledWindow: scrolledWindow,
		debugger:       debugger,
		textView:       textView,
		textBuffer:     textView.Buffer(),
	}

	// Create text tags for styling
	panel.initTextTags()

	// Subscribe to events
	panel.subscribeToEvents()

	// Log initial message
	panel.LogInfo("Console initialized")

	return panel
}

func (c *ConsolePanel) initTextTags() {
	tagTable := c.textBuffer.TagTable()

	// Error tag - red text
	c.errorTag = gtk.NewTextTag("error")
	c.errorTag.SetObjectProperty("foreground", "red")
	tagTable.Add(c.errorTag)

	// Warning tag - orange text
	c.warningTag = gtk.NewTextTag("warning")
	c.warningTag.SetObjectProperty("foreground", "#FF8C00")
	tagTable.Add(c.warningTag)
}

func (c *ConsolePanel) subscribeToEvents() {
	c.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		if data.State != nil {
			if data.State.Exited {
				c.LogInfo(fmt.Sprintf("Process exited with code %d", data.State.ExitStatus))
			} else if data.State.Running {
				c.LogInfo("State changed: Running")
			} else {
				c.LogInfo("State changed: Paused")
				// Check if breakpoint was hit
				if data.State.CurrentThread != nil && data.State.CurrentThread.Breakpoint != nil {
					bp := data.State.CurrentThread.Breakpoint
					c.LogInfo(fmt.Sprintf("Breakpoint hit at %s:%d", bp.File, bp.Line))
				}
			}
		}
	})

	c.debugger.Signals.Connect(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
		switch data.BreakpointOperation {
		case BreakpointCreated:
			if data.Breakpoint != nil {
				c.LogInfo(fmt.Sprintf("Breakpoint created at %s:%d", data.Breakpoint.File, data.Breakpoint.Line))
			}
		case BreakpointRemoved:
			if data.Breakpoint != nil {
				c.LogInfo(fmt.Sprintf("Breakpoint removed from %s:%d", data.Breakpoint.File, data.Breakpoint.Line))
			}
		case BreakpointEnabled:
			if data.Breakpoint != nil {
				c.LogInfo(fmt.Sprintf("Breakpoint enabled at %s:%d", data.Breakpoint.File, data.Breakpoint.Line))
			}
		case BreakpointDisabled:
			if data.Breakpoint != nil {
				c.LogInfo(fmt.Sprintf("Breakpoint disabled at %s:%d", data.Breakpoint.File, data.Breakpoint.Line))
			}
		case BreakpointsAllEnabled:
			c.LogInfo("All breakpoints enabled")
		case BreakpointsAllDisabled:
			c.LogInfo("All breakpoints disabled")
		case BreakpointsAllCleared:
			c.LogInfo("All breakpoints cleared")
		}
	})

	c.debugger.Signals.Connect(EventProcessExited, func(event DebuggerEvent, data EventData) {
		c.LogError("Process exited")
	})

	c.debugger.Signals.Connect(EventProcessRestarted, func(event DebuggerEvent, data EventData) {
		c.LogSuccess("Process restarted")
	})

	c.debugger.Signals.Connect(EventContinue, func(event DebuggerEvent, data EventData) {
		c.LogInfo("Process continued")
	})

	c.debugger.Signals.Connect(EventStep, func(event DebuggerEvent, data EventData) {
		c.LogInfo("Step into")
	})

	c.debugger.Signals.Connect(EventStepOver, func(event DebuggerEvent, data EventData) {
		c.LogInfo("Step over")
	})

	c.debugger.Signals.Connect(EventStepOut, func(event DebuggerEvent, data EventData) {
		c.LogInfo("Step out")
	})
}

func (c *ConsolePanel) LogInfo(message string) {
	c.appendText(fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05"), message))
}

func (c *ConsolePanel) LogWarning(message string) {
	text := fmt.Sprintf("[%s] WARNING: %s\n", time.Now().Format("15:04:05"), message)
	endIter := c.textBuffer.EndIter()
	offset := endIter.Offset()
	c.textBuffer.Insert(endIter, text)

	// Apply warning tag to the inserted text
	startIter := c.textBuffer.IterAtOffset(offset)
	endIter = c.textBuffer.EndIter()
	c.textBuffer.ApplyTag(c.warningTag, startIter, endIter)

	// Auto-scroll to end
	endIter = c.textBuffer.EndIter()
	mark := c.textBuffer.CreateMark("end", endIter, false)
	c.textView.ScrollToMark(mark, 0.0, false, 0.0, 1.0)
}

func (c *ConsolePanel) LogError(message string) {
	text := fmt.Sprintf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), message)
	endIter := c.textBuffer.EndIter()
	offset := endIter.Offset()
	c.textBuffer.Insert(endIter, text)

	// Apply error tag to the inserted text
	startIter := c.textBuffer.IterAtOffset(offset)
	endIter = c.textBuffer.EndIter()
	c.textBuffer.ApplyTag(c.errorTag, startIter, endIter)

	// Auto-scroll to end
	endIter = c.textBuffer.EndIter()
	mark := c.textBuffer.CreateMark("end", endIter, false)
	c.textView.ScrollToMark(mark, 0.0, false, 0.0, 1.0)
}

func (c *ConsolePanel) LogSuccess(message string) {
	c.appendText(fmt.Sprintf("[%s] ✓ %s\n", time.Now().Format("15:04:05"), message))
}

func (c *ConsolePanel) appendText(text string) {
	endIter := c.textBuffer.EndIter()
	c.textBuffer.Insert(endIter, text)

	// Auto-scroll to end
	endIter = c.textBuffer.EndIter()
	mark := c.textBuffer.CreateMark("end", endIter, false)
	c.textView.ScrollToMark(mark, 0.0, false, 0.0, 1.0)
}
