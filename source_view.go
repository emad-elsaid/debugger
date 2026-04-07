package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	. "github.com/emad-elsaid/debugger/structs"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/service/api"
	"libdb.so/gotk4-sourceview/pkg/gtksource/v5"
)

const (
	breakpointEnabledCategory  = "breakpoint-enabled"
	breakpointDisabledCategory = "breakpoint-disabled"
)

// SourceView provides a code editor with syntax highlighting and breakpoint support
type SourceView struct {
	*gtk.Box
	debugger *Debugger

	scrollWin  *gtk.ScrolledWindow
	sourceView *gtksource.View
	sourceBuf  *gtksource.Buffer

	currentFile string
	currentLine int

	// Breakpoint tracking
	breakpointStates map[int]*api.Breakpoint

	// Gutter renderer for breakpoints
	gutterRenderer *gtksource.GutterRendererPixbuf
}

// NewSourceView creates a new source code viewer
func NewSourceView(debugger *Debugger) *SourceView {
	sv := &SourceView{
		Box:              gtk.NewBox(gtk.OrientationVertical, 0),
		debugger:         debugger,
		breakpointStates: make(map[int]*api.Breakpoint),
	}

	sv.sourceBuf = gtksource.NewBuffer(nil)
	sv.sourceView = gtksource.NewViewWithBuffer(sv.sourceBuf)
	sv.sourceView.SetEditable(false)
	sv.sourceView.SetShowLineNumbers(true)
	sv.sourceView.SetHighlightCurrentLine(true)
	sv.sourceView.SetMonospace(true)

	langMgr := gtksource.LanguageManagerGetDefault()
	goLang := langMgr.Language("go")
	if goLang != nil {
		sv.sourceBuf.SetLanguage(goLang)
	}

	sv.setupGutter()
	sv.setupMarkAttributes()
	sv.setupCSS()

	sv.scrollWin = gtk.NewScrolledWindow()
	sv.scrollWin.SetChild(&sv.sourceView.Widget)
	sv.scrollWin.SetVExpand(true)
	sv.scrollWin.SetHExpand(true)

	sv.Append(sv.scrollWin)

	sv.setupEventHandlers()
	sv.subscribeToEvents()

	return sv
}

func (sv *SourceView) setupGutter() {
	sv.gutterRenderer = gtksource.NewGutterRendererPixbuf()
	sv.gutterRenderer.SetXAlign(0.5)
	sv.gutterRenderer.SetYAlign(0.5)
	sv.gutterRenderer.SetXpad(2)

	gutter := sv.sourceView.Gutter(gtk.TextWindowLeft)
	gutter.Insert(sv.gutterRenderer, -10)

	sv.gutterRenderer.ConnectActivate(func(iter *gtk.TextIter, area *gdk.Rectangle, button uint, state gdk.ModifierType, nPresses int) {
		if button == gdk.BUTTON_PRIMARY {
			line := iter.Line() + 1
			sv.cycleBreakpoint(line)
		}
	})
}

func (sv *SourceView) setupMarkAttributes() {
	enabledAttrs := gtksource.NewMarkAttributes()
	enabledAttrs.SetIconName("media-record")
	bg := gdk.NewRGBA(1.0, 0.0, 0.0, 0.3)
	enabledAttrs.SetBackground(&bg)
	sv.sourceView.SetMarkAttributes(breakpointEnabledCategory, enabledAttrs, 1)

	disabledAttrs := gtksource.NewMarkAttributes()
	disabledAttrs.SetIconName("media-record")
	bgDisabled := gdk.NewRGBA(0.5, 0.5, 0.5, 0.3)
	disabledAttrs.SetBackground(&bgDisabled)
	sv.sourceView.SetMarkAttributes(breakpointDisabledCategory, disabledAttrs, 1)
}

func (sv *SourceView) setupCSS() {
	provider := gtk.NewCSSProvider()
	provider.LoadFromString(`
		.current-line {
			background-color: rgba(144, 238, 144, 0.3);
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		provider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}

func (sv *SourceView) setupEventHandlers() {
	gesture := gtk.NewGestureClick()
	gesture.ConnectPressed(func(n int, x, y float64) {
		if n == 1 {
			bufX, bufY := sv.sourceView.WindowToBufferCoords(gtk.TextWindowWidget, int(x), int(y))
			iter, _ := sv.sourceView.IterAtLocation(bufX, bufY)
			line := iter.Line() + 1

			rect := sv.sourceView.VisibleRect()
			if int(x) < 50 || int(x) > rect.Width()-10 {
				sv.cycleBreakpoint(line)
			}
		}
	})
	sv.sourceView.AddController(gesture)
}

func (sv *SourceView) cycleBreakpoint(line int) {
	if sv.currentFile == "" {
		return
	}

	for _, bp := range sv.debugger.Breakpoints() {
		if bp.File == sv.currentFile && bp.Line == line {
			if bp.Disabled {
				sv.debugger.ClearBreakpoint(bp)
			} else {
				bp.Disabled = true
				sv.debugger.AmendBreakpoint(bp)
			}
			return
		}
	}

	bp := &api.Breakpoint{
		File: sv.currentFile,
		Line: line,
	}
	sv.debugger.CreateBreakpoint(bp)
}

func (sv *SourceView) subscribeToEvents() {
	sv.debugger.Signals.Connect(EventFileChanged, func(event DebuggerEvent, data EventData) {
		sv.LoadFile(data.File, data.Line)
	})

	sv.debugger.Signals.Connect(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
		sv.updateBreakpoints()
	})

	sv.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		if data.State != nil && !data.State.Exited {
			sv.highlightCurrentLine()
		}
	})
}

func (sv *SourceView) LoadFile(file string, line int) {
	slog.Debug("LoadFile: starting", "file", file, "line", line)
	if file == "" {
		return
	}

	sv.currentFile = file
	sv.currentLine = line

	slog.Debug("LoadFile: opening file")
	f, err := os.Open(file)
	if err != nil {
		sv.sourceBuf.SetText(fmt.Sprintf("Error opening file: %s", err))
		return
	}
	defer f.Close()

	slog.Debug("LoadFile: reading file")
	content, err := io.ReadAll(f)
	if err != nil {
		sv.sourceBuf.SetText(fmt.Sprintf("Error reading file: %s", err))
		return
	}

	slog.Debug("LoadFile: setting text buffer")
	sv.sourceBuf.SetText(string(content))
	slog.Debug("LoadFile: highlighting current line")
	sv.highlightCurrentLine()
	slog.Debug("LoadFile: updating breakpoints")
	sv.updateBreakpoints()
	slog.Debug("LoadFile: scrolling to line")
	sv.scrollToLine(line)
	slog.Debug("LoadFile: completed")
}

func (sv *SourceView) highlightCurrentLine() {
	if sv.currentLine <= 0 {
		return
	}

	// Ensure the tag exists before attempting to remove it
	tagTable := sv.sourceBuf.TagTable()
	tag := tagTable.Lookup("current-line")
	if tag == nil {
		tag = gtk.NewTextTag("current-line")
		tagTable.Add(tag)
	}

	start, end := sv.sourceBuf.Bounds()
	sv.sourceBuf.RemoveTagByName("current-line", start, end)

	lineIter, _ := sv.sourceBuf.IterAtLine(sv.currentLine - 1)
	lineEnd := lineIter
	lineEnd.ForwardToLineEnd()

	sv.sourceBuf.ApplyTag(tag, lineIter, lineEnd)
}

func (sv *SourceView) scrollToLine(line int) {
	if line <= 0 {
		return
	}

	iter, _ := sv.sourceBuf.IterAtLine(line - 1)
	mark := sv.sourceBuf.CreateMark("", iter, false)
	sv.sourceView.ScrollToMark(mark, 0.25, true, 0.0, 0.5)
}

func (sv *SourceView) updateBreakpoints() {
	sv.removeAllMarks()

	sv.breakpointStates = make(map[int]*api.Breakpoint)
	for _, bp := range sv.debugger.Breakpoints() {
		if bp.File == sv.currentFile {
			sv.breakpointStates[bp.Line] = bp

			iter, _ := sv.sourceBuf.IterAtLine(bp.Line - 1)

			category := breakpointEnabledCategory
			if bp.Disabled {
				category = breakpointDisabledCategory
			}

			sv.sourceBuf.CreateSourceMark(fmt.Sprintf("bp-%d", bp.Line), category, iter)
		}
	}

	sv.updateGutterIcons()
}

func (sv *SourceView) removeAllMarks() {
	start, end := sv.sourceBuf.Bounds()

	sv.sourceBuf.RemoveSourceMarks(start, end, breakpointEnabledCategory)
	sv.sourceBuf.RemoveSourceMarks(start, end, breakpointDisabledCategory)
}

func (sv *SourceView) updateGutterIcons() {
	for line, bp := range sv.breakpointStates {
		iter, _ := sv.sourceBuf.IterAtLine(line - 1)

		iconName := "media-record"
		if bp.Disabled {
			sv.gutterRenderer.SetIconName("")
		} else {
			sv.gutterRenderer.SetIconName(iconName)
		}

		_ = iter
	}
}
