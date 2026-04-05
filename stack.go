package main

import (
	"fmt"
	. "github.com/emad-elsaid/debugger/structs"
	"html"
	"path/filepath"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/pkg/proc"
)

// StackPanel displays the current goroutine stack frames
type StackPanel struct {
	*gtk.ScrolledWindow
	debugger *Debugger

	listBox *gtk.ListBox
}

// NewStackPanel creates a new stack panel
func NewStackPanel(debugger *Debugger) *StackPanel {
	builder := gtk.NewBuilderFromFile("ui/stack.ui")

	sp := &StackPanel{
		ScrolledWindow: builder.GetObject("stack_panel").Cast().(*gtk.ScrolledWindow),
		debugger:       debugger,
		listBox:        builder.GetObject("stack_list").Cast().(*gtk.ListBox),
	}

	// Handle frame selection
	sp.listBox.ConnectRowActivated(func(row *gtk.ListBoxRow) {
		frameIndex := row.Index()
		sp.onFrameSelected(frameIndex)
	})

	// Subscribe to debugger events
	sp.subscribeToEvents()

	// Don't populate during construction - will update when process stops via events

	return sp
}

func (sp *StackPanel) subscribeToEvents() {
	sp.debugger.EventBus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			sp.updateStack()
			return false
		})
	})
}

func (sp *StackPanel) updateStack() {
	// Clear existing rows
	for child := sp.listBox.FirstChild(); child != nil; child = sp.listBox.FirstChild() {
		sp.listBox.Remove(child)
	}

	// Get stack frames
	frames, err := sp.debugger.Stacktrace()
	if err != nil {
		sp.showErrorState(err)
		return
	}

	// Handle empty state with context-aware messages
	if len(frames) == 0 {
		sp.showEmptyState()
		return
	}

	// Populate list with stack frames
	for i, frame := range frames {
		row := sp.createFrameRow(i, frame)
		sp.listBox.Append(row)

		// Select current frame - need to find the ListBoxRow that contains this widget
		if i == sp.debugger.StackFrame {
			// After appending, find the ListBoxRow at this index
			if listBoxRow := sp.listBox.RowAtIndex(i); listBoxRow != nil {
				sp.listBox.SelectRow(listBoxRow)
			}
		}
	}

	// Show truncation warning if we hit the stack limit
	if len(frames) == STACK_LIMIT {
		sp.showTruncationWarning()
	}
}

func (sp *StackPanel) showErrorState(err error) {
	row := adw.NewActionRow()
	row.SetTitle("Unable to retrieve stack trace")
	row.SetSubtitle(html.EscapeString(err.Error()))
	row.AddCSSClass("error")

	// Add error icon as prefix
	icon := gtk.NewImageFromIconName("dialog-error-symbolic")
	icon.AddCSSClass("error")
	row.AddPrefix(icon)

	sp.listBox.Append(row)
}

func (sp *StackPanel) showEmptyState() {
	row := adw.NewActionRow()

	// Differentiate between running and exited states
	if sp.debugger.LastState != nil && sp.debugger.LastState.Exited {
		row.SetTitle("Stack trace unavailable")
		row.SetSubtitle("Process has exited")
	} else if sp.debugger.IsRunning() {
		row.SetTitle("Stack trace unavailable")
		row.SetSubtitle("Process is running")
	} else {
		row.SetTitle("No stack frames")
		row.SetSubtitle("Process is stopped but no frames available")
	}

	row.AddCSSClass("dim-label")

	// Add info icon as prefix
	icon := gtk.NewImageFromIconName("dialog-information-symbolic")
	icon.AddCSSClass("dim-label")
	row.AddPrefix(icon)

	sp.listBox.Append(row)
}

func (sp *StackPanel) showTruncationWarning() {
	row := adw.NewActionRow()
	row.SetTitle("Stack trace truncated")
	row.SetSubtitle(fmt.Sprintf("Showing %d frames (limit reached, more frames may exist)", STACK_LIMIT))
	row.AddCSSClass("warning")

	// Add warning icon as prefix
	icon := gtk.NewImageFromIconName("dialog-warning-symbolic")
	icon.AddCSSClass("warning")
	row.AddPrefix(icon)

	sp.listBox.Append(row)
}

func (sp *StackPanel) createFrameRow(index int, frame proc.Stackframe) gtk.Widgetter {
	row := adw.NewActionRow()

	// Get function name, file, and line from the frame
	var funcName string
	file := frame.Current.File
	line := frame.Current.Line

	// Get function name
	if frame.Current.Fn != nil {
		funcName = frame.Current.Fn.Name
	} else {
		funcName = "<unknown>"
	}

	// Set title to function name
	row.SetTitle(html.EscapeString(funcName))
	row.AddCSSClass("stack-frame-row")

	// Set subtitle to file:line
	if file != "" {
		basename := filepath.Base(file)
		subtitle := fmt.Sprintf("%s:%d", basename, line)
		row.SetSubtitle(html.EscapeString(subtitle))
		row.SetTooltipText(html.EscapeString(fmt.Sprintf("%s:%d", file, line)))
	}

	// Frame number prefix with enhanced styling
	prefix := gtk.NewLabel(fmt.Sprintf("#%d", index))
	prefix.AddCSSClass("dim-label")
	prefix.AddCSSClass("stack-frame-number")
	row.AddPrefix(prefix)

	// Highlight current frame (frame 0 or selected frame)
	if index == 0 || index == sp.debugger.StackFrame {
		row.AddCSSClass("stack-frame-current")
	}

	return row
}

func (sp *StackPanel) onFrameSelected(frameIndex int) {
	// Get the selected frame details
	frames, err := sp.debugger.Stacktrace()
	if err != nil || frameIndex >= len(frames) {
		return
	}

	// Update debugger's current stack frame
	sp.debugger.StackFrame = frameIndex

	// Extract file and line from the frame
	frame := frames[frameIndex]
	file := frame.Current.File
	line := frame.Current.Line

	// Emit events to notify other components
	sp.debugger.EventBus.Emit(EventStackFrameChanged, EventData{
		StackFrame: frameIndex,
		File:       file,
		Line:       line,
	})
}
