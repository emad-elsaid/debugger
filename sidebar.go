package main

import (
	"fmt"
	"html"
	"log/slog"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	. "github.com/emad-elsaid/debugger/structs"
)

// Sidebar contains watches and breakpoints
type Sidebar struct {
	*gtk.Box
	debugger *Debugger

	watchInput      *gtk.Entry
	watchAddBtn     *gtk.Button
	watchesList     *gtk.ListBox
	breakpointsList *gtk.ListBox
}

// NewSidebar creates a collapsible sidebar
func NewSidebar(debugger *Debugger) *Sidebar {
	slog.Debug("NewSidebar: loading UI from XML")
	builder := gtk.NewBuilderFromFile("ui/sidebar.ui")

	slog.Debug("NewSidebar: getting widgets from builder")
	box := builder.GetObject("sidebar_box").Cast().(*gtk.Box)
	watchInput := builder.GetObject("watch_input").Cast().(*gtk.Entry)
	watchAddBtn := builder.GetObject("watch_add_button").Cast().(*gtk.Button)
	watchesList := builder.GetObject("watches_list").Cast().(*gtk.ListBox)
	breakpointsList := builder.GetObject("breakpoints_list").Cast().(*gtk.ListBox)

	sb := &Sidebar{
		Box:             box,
		debugger:        debugger,
		watchInput:      watchInput,
		watchAddBtn:     watchAddBtn,
		watchesList:     watchesList,
		breakpointsList: breakpointsList,
	}

	// Setup watch input handlers
	sb.setupWatchInput()

	// Subscribe to events
	sb.subscribeToEvents()

	// Initial population
	slog.Debug("NewSidebar: triggering initial breakpoints update")
	sb.updateBreakpoints()

	return sb
}

func (sb *Sidebar) setupWatchInput() {
	addWatch := func() {
		expr := sb.watchInput.Text()
		if expr == "" {
			return
		}
		sb.debugger.CreateWatch(expr)
		sb.watchInput.SetText("")
		sb.watchInput.GrabFocus()
	}

	sb.watchAddBtn.ConnectClicked(func() {
		addWatch()
	})

	sb.watchInput.ConnectActivate(func() {
		addWatch()
	})
}

func (sb *Sidebar) subscribeToEvents() {
	sb.debugger.Signals.Connect(EventWatchesChanged, func(event DebuggerEvent, data EventData) {
		sb.updateWatches()
	})

	sb.debugger.Signals.Connect(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
		sb.updateBreakpoints()
	})

	sb.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		sb.updateWatches()
	})

	sb.debugger.Signals.Connect(EventStackFrameChanged, func(event DebuggerEvent, data EventData) {
		sb.updateWatches()
	})
}

func (sb *Sidebar) updateWatches() {
	// Clear existing
	for {
		row := sb.watchesList.FirstChild()
		if row == nil {
			break
		}
		sb.watchesList.Remove(row)
	}

	// Show empty state if no watches
	if len(sb.debugger.WatchesExpr) == 0 {
		placeholder := adw.NewActionRow()
		placeholder.SetTitle("No watches")
		placeholder.SetSubtitle("Add expressions to monitor their values")
		placeholder.SetSensitive(false)
		placeholder.AddCSSClass("empty-state-label")
		sb.watchesList.Append(placeholder)
		return
	}

	// Add watches
	for _, expr := range sb.debugger.WatchesExpr {
		row := adw.NewActionRow()
		row.SetTitle(html.EscapeString(expr))

		// Evaluate expression
		if v, err := sb.debugger.ExecuteExpr(expr); err == nil && v != nil {
			val := "nil"
			if v.Value != nil {
				val = v.Value.ExactString()
				if len(val) > 50 {
					val = val[:47] + "..."
				}
			}
			row.SetSubtitle(html.EscapeString(val))
		} else {
			errorText := fmt.Sprintf("<span foreground=\"red\">%s</span>", html.EscapeString("Error evaluating"))
			row.SetSubtitle(errorText)
			row.SetUseMarkup(true)
		}

		// Delete button
		deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
		deleteBtn.SetVAlign(gtk.AlignCenter)
		deleteBtn.AddCSSClass("flat")
		exprCopy := expr // Capture for closure
		deleteBtn.ConnectClicked(func() {
			sb.debugger.DeleteWatch(exprCopy)
		})
		row.AddSuffix(deleteBtn)
		row.SetActivatableWidget(deleteBtn)

		sb.watchesList.Append(row)
	}
}

func (sb *Sidebar) updateBreakpoints() {
	// Clear existing
	for {
		row := sb.breakpointsList.FirstChild()
		if row == nil {
			break
		}
		sb.breakpointsList.Remove(row)
	}

	// Add breakpoints
	for _, bp := range sb.debugger.Breakpoints() {
		row := adw.NewActionRow()
		funcName := bp.FunctionName
		if funcName == "" {
			funcName = "unknown"
		}
		row.SetTitle(html.EscapeString(funcName))
		row.SetSubtitle(html.EscapeString(fmt.Sprintf("%s:%d", bp.File, bp.Line)))

		// Toggle button
		toggle := gtk.NewSwitch()
		toggle.SetActive(!bp.Disabled)
		toggle.SetVAlign(gtk.AlignCenter)
		bpCopy := bp // Capture for closure
		toggle.ConnectStateSet(func(state bool) bool {
			bpCopy.Disabled = !state
			sb.debugger.AmendBreakpoint(bpCopy)
			return false
		})

		// Delete button
		deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
		deleteBtn.SetVAlign(gtk.AlignCenter)
		deleteBtn.AddCSSClass("flat")
		deleteBtn.ConnectClicked(func() {
			sb.debugger.ClearBreakpoint(bpCopy)
		})

		row.AddSuffix(toggle)
		row.AddSuffix(deleteBtn)

		sb.breakpointsList.Append(row)
	}
}
