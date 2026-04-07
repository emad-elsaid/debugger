package main

import (
	"fmt"
	"html"
	"log/slog"
	"sort"
	"strings"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/pkg/proc"
)

type FunctionsPanel struct {
	*gtk.Box
	debugger *Debugger

	searchEntry *gtk.SearchEntry
	listBox     *gtk.ListBox
	scrolled    *gtk.ScrolledWindow

	allFunctions []FunctionInfo
}

func NewFunctionsPanel(debugger *Debugger) *FunctionsPanel {
	builder := gtk.NewBuilderFromFile("ui/functions.ui")

	fp := &FunctionsPanel{
		Box:         builder.GetObject("functions_panel").Cast().(*gtk.Box),
		debugger:    debugger,
		searchEntry: builder.GetObject("functions_search_entry").Cast().(*gtk.SearchEntry),
		listBox:     builder.GetObject("functions_list").Cast().(*gtk.ListBox),
		scrolled:    builder.GetObject("functions_scrolled").Cast().(*gtk.ScrolledWindow),
	}

	fp.listBox.SetFilterFunc(func(row *gtk.ListBoxRow) bool {
		searchText := strings.ToLower(fp.searchEntry.Text())
		if searchText == "" {
			return true
		}

		widget := row.Child()
		if actionRow, ok := widget.(*adw.ActionRow); ok {
			title := strings.ToLower(actionRow.Title())
			subtitle := strings.ToLower(actionRow.Subtitle())
			return strings.Contains(title, searchText) || strings.Contains(subtitle, searchText)
		}

		return true
	})

	fp.searchEntry.ConnectSearchChanged(func() {
		fp.listBox.InvalidateFilter()
	})

	fp.debugger.Signals.Connect(EventProcessRestarted, func(event DebuggerEvent, data EventData) {
		fp.loadFunctions()
	})

	// Don't populate during construction - will update when process stops

	return fp
}

func (fp *FunctionsPanel) loadFunctions() {
	fp.allFunctions = []FunctionInfo{}

	fp.debugger.WithTarget(func(target *proc.Target) {
		bi := target.BinInfo()
		if bi == nil {
			slog.Debug("Functions: binInfo is nil")
			return
		}

		funcs := bi.Functions
		slog.Debug("Functions: loading", "count", len(funcs))
		for _, fn := range funcs {
			if fn.Entry == 0 {
				continue
			}

			file, line, _ := bi.PCToLine(fn.Entry)
			if file == "" {
				continue
			}

			fp.allFunctions = append(fp.allFunctions, FunctionInfo{
				Name: fn.Name,
				File: file,
				Line: line,
			})
		}
	})

	sort.Slice(fp.allFunctions, func(i, j int) bool {
		return fp.allFunctions[i].Name < fp.allFunctions[j].Name
	})

	slog.Debug("Functions: updating list", "functions", len(fp.allFunctions))
	fp.updateList()
}

func (fp *FunctionsPanel) gatherFunctionsData(target *proc.Target) {
	bi := target.BinInfo()
	if bi == nil {
		slog.Debug("Functions: binInfo is nil")
		return
	}

	funcs := bi.Functions
	slog.Debug("Functions: loading", "count", len(funcs))
	for _, fn := range funcs {
		if fn.Entry == 0 {
			continue
		}

		file, line, _ := bi.PCToLine(fn.Entry)
		if file == "" {
			continue
		}

		fp.allFunctions = append(fp.allFunctions, FunctionInfo{
			Name: fn.Name,
			File: file,
			Line: line,
		})
	}

	sort.Slice(fp.allFunctions, func(i, j int) bool {
		return fp.allFunctions[i].Name < fp.allFunctions[j].Name
	})

	slog.Debug("Functions: data gathered", "functions", len(fp.allFunctions))
}

func (fp *FunctionsPanel) updateList() {
	for {
		row := fp.listBox.FirstChild()
		if row == nil {
			break
		}
		fp.listBox.Remove(row)
	}

	for _, fnInfo := range fp.allFunctions {
		row := adw.NewActionRow()
		row.SetTitle(html.EscapeString(fnInfo.Name))
		row.SetSubtitle(html.EscapeString(fmt.Sprintf("%s:%d", fnInfo.File, fnInfo.Line)))

		fnName := fnInfo.Name
		row.SetActivatable(true)
		row.ConnectActivated(func() {
			fp.debugger.JumpToFunction(fnName)
		})

		fp.listBox.Append(row)
	}
}
