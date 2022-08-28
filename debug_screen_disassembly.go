package main

import (
	"fmt"

	. "github.com/emad-elsaid/debugger/ui"
	"github.com/emad-elsaid/delve/pkg/proc"
)

type DisassemblyPanel struct {
	AssemblyList ClickableList
}

func NewDisassemblyPanel() DisassemblyPanel {
	return DisassemblyPanel{
		AssemblyList: NewClickableList(),
	}
}

func (s *DisassemblyPanel) Layout(d *Debugger) W {
	g, err := d.SelectedGoRoutine()
	if err != nil {
		return Centered(Label(err.Error()))
	}

	fn := g.CurrentLoc.Fn
	asm, err := d.Disassemble(g.ID, fn.Entry, fn.End)
	if err != nil {
		return Centered(Label(err.Error()))
	}

	ele := func(c C, i int) D {
		l := asm[i]
		w := Inset1(Label(l.Text(proc.GoFlavour, d.Target().BinInfo())))
		if i == 0 || l.Loc.Line != asm[i-1].Loc.Line {
			w = Rows(
				Rigid(
					Inset1(
						Bold(
							Label(fmt.Sprintf("%s:%d", l.Loc.File, l.Loc.Line)),
						),
					),
				),
				Rigid(w),
			)
		}
		return w(c)
	}

	click := func(i int) {
		l := asm[i]
		loc := l.Loc
		d.SetFileLine(loc.File, loc.Line)
	}

	return Rows(
		Rigid(
			Bold(
				Label(fmt.Sprintf("Function: %s PC:0x%x:0x%x", fn.Name, fn.Entry, fn.End)),
			),
		),
		RowSpacer1,
		Flexed(1, s.AssemblyList.Layout(len(asm), ele, click)),
	)
}
