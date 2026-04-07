package main

import (
	. "github.com/emad-elsaid/debugger/structs"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/pkg/proc"
)

type SourcesPanel struct {
	*gtk.Box
	debugger *Debugger

	searchEntry *gtk.SearchEntry
	treeView    *gtk.TreeView
	treeStore   *gtk.TreeStore
	sources     []string
}

func NewSourcesPanel(debugger *Debugger) *SourcesPanel {
	builder := gtk.NewBuilderFromFile("ui/sources.ui")

	sp := &SourcesPanel{
		Box:         builder.GetObject("sources_panel").Cast().(*gtk.Box),
		debugger:    debugger,
		searchEntry: builder.GetObject("sources_search_entry").Cast().(*gtk.SearchEntry),
		treeView:    builder.GetObject("sources_tree_view").Cast().(*gtk.TreeView),
	}

	sp.searchEntry.ConnectSearchChanged(func() {
		sp.filter()
	})

	sp.treeStore = gtk.NewTreeStore([]glib.Type{glib.TypeString, glib.TypeString})
	sp.treeView.SetModel(&sp.treeStore.TreeModel)

	renderer := gtk.NewCellRendererText()
	column := gtk.NewTreeViewColumn()
	column.PackStart(renderer, true)
	column.AddAttribute(renderer, "text", 0)
	sp.treeView.AppendColumn(column)

	sp.treeView.ConnectRowActivated(func(path *gtk.TreePath, column *gtk.TreeViewColumn) {
		iter, ok := sp.treeStore.Iter(path)
		if !ok || iter == nil {
			return
		}

		value := sp.treeStore.Value(iter, 1)
		filePath := value.GoValue()
		if filePathStr, ok := filePath.(string); ok && filePathStr != "" {
			sp.debugger.SetFileLine(filePathStr, 1)
		}
	})

	sp.subscribeToEvents()
	// Don't populate during construction - will update when process stops

	return sp
}

func (sp *SourcesPanel) subscribeToEvents() {
	sp.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		sp.loadSources()
	})

	sp.debugger.Signals.Connect(EventProcessRestarted, func(event DebuggerEvent, data EventData) {
		sp.loadSources()
	})
}

func (sp *SourcesPanel) loadSources() {
	sp.debugger.WithTarget(func(target *proc.Target) {
		binInfo := target.BinInfo()
		if binInfo == nil {
			slog.Debug("Sources: binInfo is nil")
			return
		}

		slog.Debug("Sources: loading", "count", len(binInfo.Sources))

		sourceMap := make(map[string]bool)
		for _, source := range binInfo.Sources {
			if source != "" {
				sourceMap[source] = true
			}
		}

		sp.sources = make([]string, 0, len(sourceMap))
		for source := range sourceMap {
			sp.sources = append(sp.sources, source)
		}

		sort.Strings(sp.sources)
		slog.Debug("Sources: updating tree", "sources", len(sp.sources))
		sp.updateTree()
	})
}

func (sp *SourcesPanel) gatherSourcesData(target *proc.Target) {
	binInfo := target.BinInfo()
	if binInfo == nil {
		slog.Debug("Sources: binInfo is nil")
		return
	}

	slog.Debug("Sources: loading", "count", len(binInfo.Sources))

	sourceMap := make(map[string]bool)
	for _, source := range binInfo.Sources {
		if source != "" {
			sourceMap[source] = true
		}
	}

	sp.sources = make([]string, 0, len(sourceMap))
	for source := range sourceMap {
		sp.sources = append(sp.sources, source)
	}

	sort.Strings(sp.sources)
	slog.Debug("Sources: data gathered", "sources", len(sp.sources))
}

func (sp *SourcesPanel) updateTree() {
	sp.treeStore.Clear()

	query := strings.ToLower(sp.searchEntry.Text())

	tree := make(map[string]map[string]bool)
	var packages []string

	for _, source := range sp.sources {
		if query != "" && !strings.Contains(strings.ToLower(source), query) {
			continue
		}

		pkg := filepath.Dir(source)
		if tree[pkg] == nil {
			tree[pkg] = make(map[string]bool)
			packages = append(packages, pkg)
		}
		tree[pkg][source] = true
	}

	sort.Strings(packages)

	for _, pkg := range packages {
		pkgIter := sp.treeStore.Append(nil)
		sp.treeStore.Set(pkgIter, []int{0, 1}, []glib.Value{
			*glib.NewValue(pkg),
			*glib.NewValue(""),
		})

		var files []string
		for file := range tree[pkg] {
			files = append(files, file)
		}
		sort.Strings(files)

		for _, file := range files {
			fileIter := sp.treeStore.Append(pkgIter)
			fileName := filepath.Base(file)
			sp.treeStore.Set(fileIter, []int{0, 1}, []glib.Value{
				*glib.NewValue(fileName),
				*glib.NewValue(file),
			})
		}
	}

	sp.treeView.ExpandAll()
}

func (sp *SourcesPanel) filter() {
	sp.updateTree()
}
