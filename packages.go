package main

import (
	"html"
	"log/slog"
	"sort"
	"strings"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/pkg/proc"
)

// PackagesPanel displays all packages in the target binary
type PackagesPanel struct {
	*gtk.Box
	debugger *Debugger

	searchEntry  *gtk.SearchEntry
	packagesList *gtk.ListBox
}

// NewPackagesPanel creates a new packages panel
func NewPackagesPanel(debugger *Debugger) *PackagesPanel {
	builder := gtk.NewBuilderFromFile("ui/packages.ui")

	pp := &PackagesPanel{
		Box:          builder.GetObject("packages_panel").Cast().(*gtk.Box),
		debugger:     debugger,
		searchEntry:  builder.GetObject("packages_search_entry").Cast().(*gtk.SearchEntry),
		packagesList: builder.GetObject("packages_list").Cast().(*gtk.ListBox),
	}

	// Setup search filter
	pp.packagesList.SetFilterFunc(func(row *gtk.ListBoxRow) bool {
		searchText := strings.ToLower(pp.searchEntry.Text())
		if searchText == "" {
			return true
		}

		widget := row.Child()
		if actionRow, ok := widget.(*adw.ActionRow); ok {
			title := strings.ToLower(actionRow.Title())
			return strings.Contains(title, searchText)
		}

		return true
	})

	pp.searchEntry.ConnectSearchChanged(func() {
		pp.packagesList.InvalidateFilter()
	})

	// Subscribe to events
	pp.subscribeToEvents()

	// Don't populate during construction - will update when process stops

	return pp
}

func (pp *PackagesPanel) subscribeToEvents() {
	pp.debugger.EventBus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
		glib.IdleAdd(func() bool {
			pp.updatePackages()
			return false
		})
	})
}

func (pp *PackagesPanel) loadPackages() {
	pp.updatePackages()
}

func (pp *PackagesPanel) updatePackages() {
	// Clear existing
	for {
		row := pp.packagesList.FirstChild()
		if row == nil {
			break
		}
		pp.packagesList.Remove(row)
	}

	pp.debugger.WithTarget(func(target *proc.Target) {
		bi := target.BinInfo()
		if bi == nil {
			slog.Debug("Packages: binInfo is nil")
			return
		}

		// Get all packages
		packages := bi.PackageMap
		slog.Debug("Packages: loading", "count", len(packages))
		if len(packages) == 0 {
			// Add placeholder if no packages
			row := adw.NewActionRow()
			row.SetTitle("No packages found")
			row.SetSensitive(false)
			pp.packagesList.Append(row)
			return
		}

		// Extract package names and sort alphabetically
		pkgNames := make([]string, 0, len(packages))
		for pkgName := range packages {
			pkgNames = append(pkgNames, pkgName)
		}
		sort.Strings(pkgNames)

		// Add each package in sorted order
		for _, pkgName := range pkgNames {
			row := adw.NewActionRow()
			row.SetTitle(html.EscapeString(pkgName))
			pp.packagesList.Append(row)
		}
	})
}
