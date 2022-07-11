package main

import (
	"fmt"
	"regexp"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type FunctionsPanel struct {
	List   ClickableList
	Filter widget.Editor
}

func NewFunctionsPanel() FunctionsPanel {
	return FunctionsPanel{
		List:   NewClickableList(),
		Filter: LineEditor(),
	}
}

func (f *FunctionsPanel) Layout(d *Debugger) W {
	filter := f.Filter.Text()
	regex, err := regexp.Compile(filter)
	if err != nil {
		regex, _ = regexp.Compile(regexp.QuoteMeta(filter))
	}

	fs := []string{}
	for _, f := range d.Target().BinInfo().Functions {
		if regex.MatchString(f.Name) {
			fs = append(fs, f.Name)
		}
	}

	ele := func(c C, i int) D {
		return Inset1(Label(fs[i]))(c)
	}

	click := func(i int) {
		d.JumpToFunction(fs[i])
	}

	return Rows(
		RowSpacer1,
		Rigid(
			FormRow(
				Rigid(Label(fmt.Sprintf(" %d", len(fs)))),
				ColSpacer1,
				Flexed(1, TextInput(&f.Filter, "Search Functions...")),
			),
		),
		RowSpacer1,
		Flexed(1, f.List.Layout(len(fs), ele, click)),
	)
}
