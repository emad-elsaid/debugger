package main

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/pkg/proc"
)

type TypesPanel struct {
	*gtk.Box
	debugger *Debugger

	searchEntry *gtk.SearchEntry
	listBox     *gtk.ListBox
	types       []string
}

func NewTypesPanel(debugger *Debugger) *TypesPanel {
	builder := gtk.NewBuilderFromFile("ui/types.ui")

	tp := &TypesPanel{
		Box:         builder.GetObject("types_panel").Cast().(*gtk.Box),
		debugger:    debugger,
		searchEntry: builder.GetObject("types_search_entry").Cast().(*gtk.SearchEntry),
		listBox:     builder.GetObject("types_list").Cast().(*gtk.ListBox),
	}

	tp.searchEntry.ConnectSearchChanged(func() {
		tp.filter()
	})

	tp.subscribeToEvents()
	// Don't populate during construction - will update when process stops

	return tp
}

func (tp *TypesPanel) subscribeToEvents() {
	tp.debugger.EventBus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			tp.loadTypes()
			return false
		})
	})

	tp.debugger.EventBus.Subscribe(EventProcessRestarted, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			tp.loadTypes()
			return false
		})
	})
}

func (tp *TypesPanel) loadTypes() {
	tp.types = []string{}

	tp.debugger.WithTarget(func(target *proc.Target) {
		binInfo := target.BinInfo()
		if binInfo == nil {
			slog.Debug("Types: binInfo is nil")
			return
		}

		// Get all types from binary info
		typeNames, err := binInfo.Types()
		if err != nil {
			slog.Debug("Types: error getting types", "error", err)
			return
		}

		slog.Debug("Types: loading", "count", len(typeNames))
		tp.types = typeNames
		sort.Strings(tp.types)
	})

	tp.updateList()
}

func (tp *TypesPanel) gatherTypesData(target *proc.Target) {
	binInfo := target.BinInfo()
	if binInfo == nil {
		slog.Debug("Types: binInfo is nil")
		return
	}

	// Get all types from binary info
	typeNames, err := binInfo.Types()
	if err != nil {
		slog.Debug("Types: error getting types", "error", err)
		return
	}

	slog.Debug("Types: loading", "count", len(typeNames))
	tp.types = typeNames
	sort.Strings(tp.types)
	slog.Debug("Types: data gathered", "types", len(tp.types))
}

func (tp *TypesPanel) loadTypesWithTarget(target *proc.Target) {
	tp.gatherTypesData(target)
	tp.updateList()
}

func (tp *TypesPanel) updateList() {
	for {
		row := tp.listBox.FirstChild()
		if row == nil {
			break
		}
		tp.listBox.Remove(row)
	}

	query := strings.ToLower(tp.searchEntry.Text())

	if len(tp.types) == 0 {
		placeholder := gtk.NewLabel("No types available")
		placeholder.SetXAlign(0.5)
		placeholder.AddCSSClass("dim-label")
		tp.listBox.Append(placeholder)
		return
	}

	hasMatches := false
	for _, typeName := range tp.types {
		if query != "" && !strings.Contains(strings.ToLower(typeName), query) {
			continue
		}
		hasMatches = true

		label := gtk.NewLabel(typeName)
		label.SetXAlign(0)
		label.AddCSSClass("types-list-label")

		tp.listBox.Append(label)
	}

	if !hasMatches && query != "" {
		placeholder := gtk.NewLabel("No types match '" + query + "'")
		placeholder.SetXAlign(0.5)
		placeholder.AddCSSClass("dim-label")
		tp.listBox.Append(placeholder)
	}
}

func (tp *TypesPanel) filter() {
	tp.updateList()
}
