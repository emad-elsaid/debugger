package main

import (
	"fmt"

	. "github.com/emad-elsaid/debugger/structs"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/pkg/proc"
)

type DisassemblyPanel struct {
	*gtk.ScrolledWindow
	debugger *Debugger
	textView *gtk.TextView
	textBuf  *gtk.TextBuffer
}

func NewDisassemblyPanel(debugger *Debugger) *DisassemblyPanel {
	builder := gtk.NewBuilderFromFile("ui/disassembly.ui")

	panel := &DisassemblyPanel{
		ScrolledWindow: builder.GetObject("disassembly_panel").Cast().(*gtk.ScrolledWindow),
		debugger:       debugger,
		textView:       builder.GetObject("disassembly_text_view").Cast().(*gtk.TextView),
	}

	panel.textBuf = panel.textView.Buffer()

	// Subscribe to events
	panel.subscribeToEvents()

	return panel
}

func (dp *DisassemblyPanel) subscribeToEvents() {
	dp.debugger.EventBus.Subscribe(EventFileChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			dp.disassemble()
			return false
		})
	})

	dp.debugger.EventBus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			dp.disassemble()
			return false
		})
	})

	dp.debugger.EventBus.Subscribe(EventStackFrameChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			dp.disassemble()
			return false
		})
	})
}

func (dp *DisassemblyPanel) disassemble() {
	dp.debugger.WithTarget(func(target *proc.Target) {
		if ok, _ := target.Valid(); !ok {
			dp.textBuf.SetText("No valid target")
			return
		}

		// Get selected goroutine directly from target (avoid nested WithTarget)
		g := target.SelectedGoroutine()
		if g == nil {
			routines, _, _ := proc.GoroutinesInfo(target, 0, 0)
			for _, g = range routines {
				break
			}
		}
		if g == nil {
			dp.textBuf.SetText("No goroutines available")
			return
		}

		frames, err := proc.GoroutineStacktrace(target, g, 50, 0)
		if err != nil || len(frames) == 0 {
			dp.textBuf.SetText("No stack frames available")
			return
		}

		// Use selected stack frame's PC
		frameIndex := dp.debugger.StackFrame
		if frameIndex >= len(frames) {
			frameIndex = 0
		}
		pc := frames[frameIndex].Current.PC

		// Get thread for registers
		threads := target.ThreadList()
		if len(threads) == 0 {
			dp.textBuf.SetText("No threads")
			return
		}

		thread := threads[0]
		mem := target.Memory()
		regs, err := thread.Registers()
		if err != nil {
			dp.textBuf.SetText(fmt.Sprintf("Error getting registers: %s", err))
			return
		}
		breakpoints := target.Breakpoints()
		binInfo := target.BinInfo()

		// Disassemble range centered around PC
		const rangeSize = 1024
		startPC := pc - rangeSize/2
		endPC := pc + rangeSize/2

		insts, err := proc.Disassemble(mem, regs, breakpoints, binInfo, startPC, endPC)
		if err != nil {
			dp.textBuf.SetText(fmt.Sprintf("Error disassembling: %s", err))
			return
		}

		// Build output and track current instruction position
		output := ""
		currentInstLine := -1
		lineNum := 0

		for _, inst := range insts {
			prefix := "  "
			if inst.Loc.PC == pc {
				prefix = "=>"
				currentInstLine = lineNum
			}
			output += fmt.Sprintf("%s %#x: %s\n", prefix, inst.Loc.PC, inst.Text(proc.AssemblyFlavour(0), binInfo))
			lineNum++
		}

		dp.textBuf.SetText(output)

		// Auto-scroll to current instruction
		if currentInstLine >= 0 {
			iter, _ := dp.textBuf.IterAtLine(currentInstLine)
			mark := dp.textBuf.CreateMark("current", iter, true)
			dp.textView.ScrollToMark(mark, 0.25, true, 0.0, 0.5)
			dp.textBuf.DeleteMark(mark)
		}
	})
}
