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
	watchesList       List
	watchEditor       widget.Editor
	clickables        Clickables
	debuggerLastState DebuggerState
}

func NewSideBar() SideBar {
	return SideBar{
		watches:     []VarWidget{},
		watchesList: NewVerticalList(),
		watchEditor: LineEditor(),
		clickables:  NewClickables(),
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

			d.CreateWatch(text)
			s.watchEditor.SetText("")
		}
	}

	s.ExecuteWatches(d)

	ele := func(c C, i int) D {
		if i >= len(s.watches) {
			return D{}
		}
		expr := d.WatchesExpr[i]
		varW := s.watches[i]
		var w W
		if varW == nil {
			w = Label(expr)
		} else {
			w = varW.Layout
		}

		del := func() { d.DeleteWatch(expr) }

		return Inset05(
			Columns(
				Flexed(1, w),
				Rigid(OnClick(s.clickables.Get(expr), IconDelete, del)),
			),
		)(c)
	}

	return Panel("Watches",
		Rows(
			Rigid(ZebraList(&s.watchesList, len(d.WatchesExpr), ele)),
			Rigid(TextInput(&s.watchEditor, "New watch expression...")),
		),
	)
}

func (s *SideBar) ExecuteWatches(d *Debugger) {
	if d.State == s.debuggerLastState && len(s.watches) == len(d.WatchesExpr) {
		return
	}

	s.debuggerLastState = d.State

	results := []VarWidget{}
	for _, expr := range d.WatchesExpr {
		val, err := d.ExecuteExpr(expr)
		if err != nil {
			results = append(results, nil)
		} else {
			results = append(results, NewVarWidget(val))
		}
	}

	s.watches = results
}
