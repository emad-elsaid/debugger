package main

import (
	"fmt"
	"log"
	"strings"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type SourcesPanel struct {
	List   ClickableList
	Filter widget.Editor
}

func NewSourcesPanel() SourcesPanel {
	return SourcesPanel{
		List:   NewClickableList(),
		Filter: LineEditor(),
	}
}

func (s *SourcesPanel) Layout(d *Debugger) W {
	f := s.Filter.Text()
	ss := d.Target().BinInfo().Sources
	filtered := []string{}
	for i := range ss {
		if strings.Contains(ss[i], f) {
			filtered = append(filtered, ss[i])
		}
	}

	ele := func(c C, i int) D {
		w := Inset1(Label(filtered[i]))

		if filtered[i] == d.File {
			w = Background(HighlightColor, w)
		}

		return w(c)
	}

	click := func(i int) {
		f := filtered[i]
		log.Printf("Switching file: %s", f)
		d.SetFileLine(f, 0)
	}

	return Rows(
		RowSpacer1,
		Rigid(
			FormRow(
				Rigid(Label(fmt.Sprintf(" %d", len(filtered)))),
				ColSpacer1,
				Flexed(1, TextInput(&s.Filter, "Search Sources...")),
			),
		),
		Rigid(HSpacer1),
		Flexed(1, s.List.Layout(len(filtered), ele, click)),
	)
}
