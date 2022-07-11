package main

import (
	"fmt"
	"sort"
	"strings"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type PackagesPanel struct {
	Filter widget.Editor
	List   ClickableList
}

func NewPackagesPanel() PackagesPanel {
	return PackagesPanel{
		Filter: LineEditor(),
		List:   NewClickableList(),
	}
}

func (p *PackagesPanel) Layout(d *Debugger) W {
	filter := p.Filter.Text()

	pkgsMap := map[string]bool{}
	fs := d.Target().BinInfo().Functions
	for _, i := range fs {
		pkgsMap[i.PackageName()] = true
	}

	pkgs := make([]string, 0, len(pkgsMap))
	for i := range pkgsMap {
		if strings.Contains(i, filter) {
			pkgs = append(pkgs, i)
		}
	}

	sort.Strings(pkgs)

	ele := func(c C, i int) D {
		return Inset1(Label(pkgs[i]))(c)
	}

	click := func(i int) {
		OpenBrowser(fmt.Sprintf("https://pkg.go.dev/%s", pkgs[i]))
	}

	return Rows(
		RowSpacer1,
		Rigid(
			FormRow(
				Rigid(Label(fmt.Sprintf(" %d", len(pkgs)))),
				ColSpacer1,
				Flexed(1, TextInput(&p.Filter, "Search Packages...")),
			),
		),
		RowSpacer1,
		Flexed(1,
			p.List.Layout(len(pkgs), ele, click),
		),
	)
}
