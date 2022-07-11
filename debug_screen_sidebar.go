package main

import (
	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var lfmt = message.NewPrinter(language.English)

type SideBar struct {
	watches           []VarWidget
	watchesExpr       []string
	watchesList       List
	watchEditor       widget.Editor
	debuggerLastState DebuggerState
}

func NewSideBar() SideBar {
	return SideBar{
		watches:     []VarWidget{},
		watchesList: NewVerticalList(),
		watchEditor: LineEditor(),
	}
}

func (s *SideBar) Layout(d *Debugger) W {
	return Background(SidebarBgColor, Rows(
		Rigid(s.Watches(d)),
	))
}

func (s *SideBar) Watches(d *Debugger) W {
	evts := s.watchEditor.Events()
	for _, evt := range evts {
		switch v := evt.(type) {
		case widget.SubmitEvent:
			text := v.Text
			if len(text) == 0 {
				continue
			}

			s.watchesExpr = append(s.watchesExpr, text)
			s.watchEditor.SetText("")
		}
	}

	s.ExecuteWatches(d)

	ele := func(c C, i int) D {
		if i >= len(s.watches) {
			return D{}
		}
		expr := s.watchesExpr[i]
		varW := s.watches[i]
		var w W
		if varW == nil {
			w = Label(expr)
		} else {
			w = varW.Layout
		}

		del := func() { s.DeleteWatch(expr) }

		return Inset05(
			Columns(
				Flexed(1, w),
				Rigid(OnClick(expr, IconDelete, del)),
			),
		)(c)
	}

	return Panel("Watches",
		Rows(
			Rigid(ZebraList(&s.watchesList, len(s.watchesExpr), ele)),
			Rigid(TextInput(&s.watchEditor, "New watch expression...")),
		),
	)
}

func (s *SideBar) ExecuteWatches(d *Debugger) {
	if d.State == s.debuggerLastState && len(s.watches) == len(s.watchesExpr) {
		return
	}

	s.debuggerLastState = d.State

	results := []VarWidget{}
	for _, expr := range s.watchesExpr {
		val, err := d.ExecuteExpr(expr)
		if err != nil {
			results = append(results, nil)
		} else {
			results = append(results, NewVarWidget(val))
		}
	}

	s.watches = results
}

func (s *SideBar) DeleteWatch(expr string) {
	nWatches := []string{}
	for _, v := range s.watchesExpr {
		if v != expr {
			nWatches = append(nWatches, v)
		}
	}
	s.watchesExpr = nWatches
}
