package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
	"github.com/go-delve/delve/pkg/proc"
)

type StackPanel struct {
	Filter        widget.Editor
	List          ClickableList
	RoutineList   ClickableList
	StackVarsList List
	Vars          []W

	// Cache
	GoRoutines    []*proc.G
	DebuggerState DebuggerState
}

func NewStackPanel() StackPanel {
	return StackPanel{
		Filter:        LineEditor(),
		List:          NewClickableList(),
		RoutineList:   NewClickableList(),
		StackVarsList: NewVerticalList(),
		Vars:          []W{},
		GoRoutines:    []*proc.G{},
	}
}

func (sp *StackPanel) Layout(d *Debugger) W {
	if d.State != sp.DebuggerState {
		sp.DebuggerState = d.State
		sp.SetVariables(d)
		sp.SetGoRoutines(d)
	}

	return Columns(
		Flexed(1, sp.GoroutinesPanel(d)),
		ColSpacer2,
		Flexed(4, sp.StackPanel(d)),
		ColSpacer2,
		Flexed(2, sp.VariablesPanel),
	)
}

func (sp *StackPanel) GoroutinesPanel(d *Debugger) W {
	count := len(sp.GoRoutines)

	elem := func(c C, i int) D {
		if i >= count {
			return D{}
		}

		r := sp.GoRoutines[i]
		g, _ := d.SelectedGoRoutine()
		label := ""
		if r.CurrentLoc.Fn != nil {
			label = r.CurrentLoc.Fn.Name
		} else {
			label = fmt.Sprintf("%s:%d", r.CurrentLoc.File, r.CurrentLoc.Line)
		}

		w := Inset1(Label(fmt.Sprintf("%d: %s", r.ID, label)))

		if g != nil && g.ID == r.ID {
			w = Background(HighlightColor, w)
		}

		return w(c)
	}

	onClick := func(i int) {
		if i >= count {
			return
		}

		r := sp.GoRoutines[i]
		log.Printf("Switching Go routine to: %d: %s", r.ID, r.Go().File)
		d.Target().SwitchGoroutine(r)
	}

	return Panel(
		fmt.Sprintf("Go Routines: (%d)", count),
		sp.RoutineList.Layout(count, elem, onClick),
	)
}

func (sp *StackPanel) StackPanel(d *Debugger) W {
	stack, err := d.Stacktrace()
	if err != nil {
		return Centered(Label(err.Error()))
	}

	filter := sp.Filter.Text()
	stk := make([]proc.Stackframe, 0, len(stack))
	for _, i := range stack {
		l := ""
		if i.Current.Fn != nil {
			l = i.Current.Fn.Name
		} else {
			l = fmt.Sprintf("%s:%d", i.Current.File, i.Current.Line)
		}

		if strings.Contains(l, filter) {
			stk = append(stk, i)
		}
	}

	ele := func(c C, i int) D {
		s := stk[i]
		l := ""
		if s.Current.Fn != nil {
			l = s.Current.Fn.Name
		} else {
			l = fmt.Sprintf("%s:%d", s.Current.File, s.Current.Line)
		}

		w := Inset1(Label(l))

		if d.File == s.Current.File && d.Line == s.Current.Line {
			w = Background(HighlightColor, w)
		}

		return w(c)
	}

	click := func(i int) {
		d.StackFrame = i
		s := stk[i]
		d.SetFileLine(s.Current.File, s.Current.Line)
	}

	return Rows(
		RowSpacer1,
		Rigid(
			FormRow(
				Rigid(Label(fmt.Sprintf("%d", len(stk)))),
				ColSpacer1,
				Flexed(1, TextInput(&sp.Filter, "Search Stack...")),
			),
		),
		RowSpacer1,
		Flexed(1, sp.List.Layout(len(stk), ele, click)),
	)
}

func (sp *StackPanel) SetVariables(d *Debugger) {
	sp.Vars = []W{}

	args, err := d.FunctionArguments()
	if err != nil {
		return
	}
	sp.varsToVarWidget("Arguments", args)

	vars, err := d.LocalVariables()
	if err != nil {
		return
	}
	sp.varsToVarWidget("Local Variables", vars)

	g, err := d.SelectedGoRoutine()
	if err != nil {
		return
	}
	sp.GoRoutineToWidgets(g)

}

func (sp *StackPanel) SetGoRoutines(d *Debugger) {
	routines, err := d.GoRoutines()
	if err == nil {
		sp.GoRoutines = routines
	}
}

func (sp *StackPanel) varsToVarWidget(title string, args []*proc.Variable) {
	sp.Vars = append(sp.Vars, Bold(Text(title)))
	for _, v := range args {
		w := NewVarWidget(v)
		sp.Vars = append(sp.Vars, w.Layout)
	}
}

func (sp *StackPanel) GoRoutineToWidgets(g *proc.G) {
	sp.Vars = append(sp.Vars,
		Bold(Text("Go Routine properties")),
		Text(fmt.Sprintf("Status: %s", GoRoutineStatus(g.Status))),
		Text(fmt.Sprintf("SystemStack: %v", g.SystemStack)),
		Text(fmt.Sprintf("Wait Since: %s", time.Duration(g.WaitSince))),
		Text(fmt.Sprintf("Wait Reason: %s", waitReason(g.WaitReason))),
		Text(fmt.Sprintf("Go statement: %s:%d", g.Go().File, g.Go().Line)),
	)
}

func (sp *StackPanel) VariablesPanel(c C) D {
	ele := func(c C, i int) D {
		item := sp.Vars[i]
		return item(c)
	}
	return sp.StackVarsList.Layout(c, len(sp.Vars), ele)
}
