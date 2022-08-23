package main

import (
	"fmt"
	"sort"

	. "github.com/emad-elsaid/debugger/ui"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	BreakPointColor         = DangerColor
	BreakPointDisabledColor = MixColor(BreakPointColor, WHITE, 50)

	IconBreakPoint         = Icon(icons.AVFiberManualRecord, BreakPointColor)
	IconBreakPointDisabled = Icon(icons.AVFiberManualRecord, BreakPointDisabledColor)
	IconDelete             = Icon(icons.ContentClear, DangerColor)
)

type BreakpointsPanel struct {
	BreakpointsList ClickableList
	EnableAll       Clickable
	DisableAll      Clickable
	ClearAll        Clickable
}

func NewBreakpointsPanel() BreakpointsPanel {
	return BreakpointsPanel{
		BreakpointsList: NewClickableList(),
	}
}

func (b *BreakpointsPanel) Layout(d *Debugger) W {
	return Rows(
		Rigid(b.Toolbar(d)),
		Flexed(1, b.BreakpointsPanel(d)),
	)
}

func (b *BreakpointsPanel) BreakpointsPanel(d *Debugger) W {
	breakpoints := d.Breakpoints()

	sort.Slice(breakpoints, func(i, j int) bool {
		return breakpoints[i].ID < breakpoints[j].ID
	})

	IconFont := FontEnlarge(2)

	elem := func(c C, i int) D {
		r := breakpoints[i]
		id := fmt.Sprintf("breakpoint-%d", r.ID)
		delId := fmt.Sprintf("breakpoint-del-%d", r.ID)

		click := func() { d.ToggleBreakpoint(r) }
		delClick := func() { d.ClearBreakpoint(r) }

		return Columns(
			Rigid(IconFont(Checkbox(id, !r.Disabled, click))),
			Flexed(1,
				Rows(
					Rigid(Label(r.Name)),
					Rigid(Wrap(Text(fmt.Sprintf("%s:%d", r.File, r.Line)), TextColor(SecondaryTextColor), MaxLines(3))),
				),
			),
			Rigid(OnClick(delId, IconFont(IconDelete), delClick)),
		)(c)
	}

	click := func(i int) {
		bp := breakpoints[i]
		d.SetFileLine(bp.File, bp.Line)
	}

	return b.BreakpointsList.Layout(len(breakpoints), elem, click)
}

func (b *BreakpointsPanel) Toolbar(d *Debugger) W {
	if b.EnableAll.Clicked() {
		d.EnableAllBreakpoints()
	}

	if b.DisableAll.Clicked() {
		d.DisableAllBreakpoints()
	}

	if b.ClearAll.Clicked() {
		d.ClearAllBreakpoints()
	}

	IconFont := FontEnlarge(2)

	return Background(ToolbarBgColor,
		Columns(
			Rigid(ToolbarButton(&b.EnableAll, IconFont(IconBreakPoint), "Enable All")),
			Rigid(ToolbarButton(&b.DisableAll, IconFont(IconBreakPointDisabled), "Disable All")),
			Rigid(ToolbarButton(&b.ClearAll, IconFont(IconDelete), "Clear All")),
		),
	)
}
