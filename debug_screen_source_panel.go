package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/layout"
	. "github.com/emad-elsaid/debugger/ui"
	. "github.com/emad-elsaid/types"
	"github.com/go-delve/delve/service/api"
)

var (
	HighlightColor = BANANA_100
	RunningColor   = LIME_100
)

type SourcePanel struct {
	layout.List
	content           []string
	File              string
	Line              int
	LastUpdate        time.Time
	LinesBtns         map[int]*Clickable
	BpBtns            map[int]*Clickable
	DebuggerLastState DebuggerState
	Vars              map[int64][]VarWidget
}

func NewSourcePanel() SourcePanel {
	return SourcePanel{
		List:      NewVerticalList(),
		content:   []string{},
		LinesBtns: map[int]*Clickable{},
		BpBtns:    map[int]*Clickable{},
		Vars:      map[int64][]VarWidget{},
	}
}

func (s *SourcePanel) Layout(d *Debugger) W {
	s.LoadFile(d.File)
	s.ScrollTo(d.File, d.Line)
	s.LoadVarsWidgets(d)

	bps := map[int]bool{}
	breakpoints := d.Breakpoints()
	for _, v := range breakpoints {
		if v.File == d.File {
			bps[v.Line] = !v.Disabled
		}
	}

	running := d.RunningLines()

	click := func(i int) {
		var err error
		name := fmt.Sprintf("%s:%d", d.File, i+1)

		if d.IsRunning() {
			d.Stop()
			defer d.Continue()
		}

		bp := d.FindBreakpointByName(name)

		if bp == nil {
			bp = &api.Breakpoint{
				Name: name,
				File: d.File,
				Line: i + 1,
			}

			_, err = d.CreateBreakpoint(bp)
		} else if bp.Disabled {
			_, err = d.ClearBreakpoint(bp)
		} else {
			bp.Disabled = true
			err = d.AmendBreakpoint(bp)
		}

		if err != nil {
			log.Println(LogError, err.Error())
		}
	}

	return s.LayoutLines(d, running, bps, click)
}

func (s *SourcePanel) LoadFile(file string) {
	stat, err := os.Stat(file)
	if err != nil {
		s.SetText("")
		return
	}

	if file == s.File && stat.ModTime() == s.LastUpdate {
		return
	}

	s.File = file
	s.LastUpdate = stat.ModTime()

	c, err := os.ReadFile(file)
	if err != nil {
		s.SetText("")
		return
	}

	s.SetText(string(c))
}

func (s *SourcePanel) ScrollTo(file string, line int) {
	if file == s.File && line == s.Line {
		return
	}

	s.File = file
	s.Line = line

	topLines := s.Position.Count / 2
	if topLines <= 0 {
		topLines = 5
	}

	scrollTo := line - topLines
	if scrollTo < 0 {
		scrollTo = 0
	}

	s.Position.First = scrollTo
}

func (s *SourcePanel) LineBtn(i int) *Clickable {
	if btn, ok := s.LinesBtns[i]; ok {
		return btn
	}

	btn := Clickable{}
	s.LinesBtns[i] = &btn
	return &btn
}

func (s *SourcePanel) SetText(t string) {
	s.content = strings.Split(t, "\n")
}

func (src *SourcePanel) LayoutLines(d *Debugger, running []int, bps map[int]bool, BpClick func(int)) W {
	content := src.content
	l := len(src.content)

	elem := func(c C, i int) D {
		lineNo := i + 1

		if _, ok := src.BpBtns[i]; !ok {
			src.BpBtns[i] = new(Clickable)
		}
		btn := src.BpBtns[i]
		if btn.Clicked() {
			BpClick(i)
		}

		line := strings.ReplaceAll(content[i], "\t", "    ")

		bgCol := BackgroundColor
		fgCol := SecondaryTextColor
		if v, ok := bps[lineNo]; ok {
			fgCol = BackgroundColor
			if v {
				bgCol = BreakPointColor
			} else {
				bgCol = BreakPointDisabledColor
			}
		}

		bp := Wrap(Label(fmt.Sprintf("%03d", lineNo)), Inset05, AlignEnd, TextColor(fgCol))
		if bgCol != BackgroundColor {
			bp = Background(bgCol, bp)
		}
		bp = RoundedCorners(bp)

		bpWithButton := LayoutToWidget(btn.Layout, bp)

		lineBtn := src.LineBtn(i)
		if lineBtn.Clicked() {
			d.Line = i + 1
			src.Line = i + 1
		}

		varRows := []FlexChild{}
		if vs, ok := src.Vars[int64(lineNo)]; ok {
			for _, lineVar := range vs {
				varRows = append(varRows, Rigid(lineVar.Layout))
			}
		}

		w := Columns(
			ColSpacer1,
			Rigid(bpWithButton),
			ColSpacer1,
			Rigid(Inset05(Label(line))),
			ColSpacer1,
			Rigid(
				Wrap(Rows(varRows...), Inset05, TextColor(SecondaryTextColor)),
			),
		)

		if Slice[int](running).Include(lineNo) {
			w = Background(RunningColor, w)
		} else if d.Line == lineNo {
			w = Background(HighlightColor, w)
		}

		return lineBtn.Layout(c, w)
	}

	return func(c C) D {
		return src.List.Layout(c, l, elem)
	}
}

func (s *SourcePanel) LoadVarsWidgets(d *Debugger) {
	if d.LastState == nil || d.State == s.DebuggerLastState {
		return
	}

	s.DebuggerLastState = d.State
	s.Vars = map[int64][]VarWidget{}

	locals, _ := d.LocalVariables()
	for _, v := range locals {
		w := NewVarWidget(v)
		if _, ok := s.Vars[v.DeclLine]; !ok {
			s.Vars[v.DeclLine] = []VarWidget{w}
		} else {
			s.Vars[v.DeclLine] = append(s.Vars[v.DeclLine], w)
		}
	}

	args, _ := d.FunctionArguments()
	for _, v := range args {
		w := NewVarWidget(v)
		if _, ok := s.Vars[v.DeclLine]; !ok {
			s.Vars[v.DeclLine] = []VarWidget{w}
		} else {
			s.Vars[v.DeclLine] = append(s.Vars[v.DeclLine], w)
		}
	}
}
