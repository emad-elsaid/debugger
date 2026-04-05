package main

import (
	"fmt"
	. "github.com/emad-elsaid/debugger/structs"
	"html"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type BreakpointsPanel struct {
	*gtk.Box
	debugger      *Debugger
	listBox       *gtk.ListBox
	enableAllBtn  *gtk.Button
	disableAllBtn *gtk.Button
	clearAllBtn   *gtk.Button
}

func NewBreakpointsPanel(debugger *Debugger) *BreakpointsPanel {
	builder := gtk.NewBuilderFromFile("ui/breakpoints.ui")

	bp := &BreakpointsPanel{
		Box:           builder.GetObject("breakpoints_panel").Cast().(*gtk.Box),
		debugger:      debugger,
		listBox:       builder.GetObject("breakpoints_list").Cast().(*gtk.ListBox),
		enableAllBtn:  builder.GetObject("enable_all_btn").Cast().(*gtk.Button),
		disableAllBtn: builder.GetObject("disable_all_btn").Cast().(*gtk.Button),
		clearAllBtn:   builder.GetObject("clear_all_btn").Cast().(*gtk.Button),
	}

	bp.connectButtons()
	bp.subscribeToEvents()
	// Don't populate during construction - will update when breakpoints change

	return bp
}

func (bp *BreakpointsPanel) connectButtons() {
	bp.enableAllBtn.ConnectClicked(func() {
		if err := bp.debugger.EnableAllBreakpoints(); err != nil {
			// Error already logged by debugger
			return
		}
	})

	bp.disableAllBtn.ConnectClicked(func() {
		if err := bp.debugger.DisableAllBreakpoints(); err != nil {
			// Error already logged by debugger
			return
		}
	})

	bp.clearAllBtn.ConnectClicked(func() {
		if err := bp.debugger.ClearAllBreakpoints(); err != nil {
			// Error already logged by debugger
			return
		}
	})
}

func (bp *BreakpointsPanel) subscribeToEvents() {
	bp.debugger.EventBus.Subscribe(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			bp.updateBreakpoints()
			return false
		})
	})
}

func (bp *BreakpointsPanel) updateBreakpoints() {
	for {
		row := bp.listBox.FirstChild()
		if row == nil {
			break
		}
		bp.listBox.Remove(row)
	}

	for _, breakpoint := range bp.debugger.Breakpoints() {
		row := adw.NewActionRow()
		// Escape HTML to prevent markup parsing errors with "<multiple locations>"
		title := html.EscapeString(fmt.Sprintf("%s:%d", breakpoint.File, breakpoint.Line))
		row.SetTitle(title)

		funcName := breakpoint.FunctionName
		if funcName == "" {
			funcName = "unknown"
		}
		subtitle := html.EscapeString(fmt.Sprintf("%s (hit: %d)", funcName, breakpoint.TotalHitCount))
		row.SetSubtitle(subtitle)

		toggle := gtk.NewSwitch()
		toggle.SetActive(!breakpoint.Disabled)
		toggle.SetVAlign(gtk.AlignCenter)

		bpCopy := breakpoint
		toggle.ConnectStateSet(func(state bool) bool {
			if state {
				bpCopy.Disabled = false
			} else {
				bpCopy.Disabled = true
			}
			bp.debugger.AmendBreakpoint(bpCopy)
			return false
		})

		deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
		deleteBtn.SetVAlign(gtk.AlignCenter)
		deleteBtn.AddCSSClass("flat")
		deleteBtn.ConnectClicked(func() {
			bp.debugger.ClearBreakpoint(bpCopy)
		})

		row.AddSuffix(toggle)
		row.AddSuffix(deleteBtn)

		bp.listBox.Append(row)
	}
}
