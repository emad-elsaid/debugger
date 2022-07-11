package main

import (
	"fmt"
	"regexp"
	"sort"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type TypesPanel struct {
	List   List
	Filter widget.Editor
}

func NewTypesPanel() TypesPanel {
	return TypesPanel{
		List:   NewVerticalList(),
		Filter: LineEditor(),
	}
}

func (t *TypesPanel) Layout(d *Debugger) W {
	filter := t.Filter.Text()
	regex, err := regexp.Compile(filter)
	if err != nil {
		regex, _ = regexp.Compile(regexp.QuoteMeta(filter))
	}

	allTypes, err := d.Target().BinInfo().Types()
	if err != nil {
		allTypes = []string{}
	}

	types := make([]string, 0, len(allTypes))
	for _, typ := range allTypes {
		if regex.Match([]byte(typ)) {
			types = append(types, typ)
		}
	}

	sort.Strings(types)

	ele := func(c C, i int) D {
		return Inset1(Label(types[i]))(c)
	}

	return Rows(
		RowSpacer1,
		Rigid(
			FormRow(
				Rigid(Label(fmt.Sprintf(" %d", len(types)))),
				ColSpacer1,
				Flexed(1, TextInput(&t.Filter, "Search types...")),
			),
		),
		RowSpacer1,
		Flexed(1, ZebraList(&t.List, len(types), ele)),
	)
}
