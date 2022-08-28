package main

import (
	"fmt"
	"os"

	"gioui.org/op"
	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type OpenScreen struct {
	Path             string
	BinaryEditor     widget.Editor
	Arguments        widget.Editor
	Test             bool
	RunImmediately   bool
	Clickables       Clickables
	Error            error
	DirectoryBrowser DirectoryBrowser
	Sessions         []Session
	SessionsList     List
	Tabs             Tabs
}

func (o *OpenScreen) Layout(c C) D {
	return o.Tabs.Layout(
		&TabChild{Name: "Open Project", Panel: func(c C) D { return o.Browser(c) }},
		&TabChild{Name: "Recent Projects", Panel: func(c C) D { return o.SessionsListWidget(c) }},
	)(c)
}

func (o *OpenScreen) Browser(c C) D {
	if o.Clickables.Get("debug").Clicked() {
		o.OpenPath()
	}

	checkboxFont := FontEnlarge(2)

	children := []FlexChild{
		Flexed(2, o.DirectoryBrowser.Layout(&o.Path)),
		RowSpacer3,
		Rigid(
			Inset1(
				Rows(
					Rigid(
						FormRow(
							Flexed(1, Label("Output binary:")),
							ColSpacer1,
							Flexed(5, TextInput(&o.BinaryEditor, "")),
						),
					),
					RowSpacer1,
					Rigid(
						FormRow(
							Flexed(1, Label("Arguments:")),
							ColSpacer1,
							Flexed(5, TextInput(&o.Arguments, "")),
						),
					),
					RowSpacer1,
					Rigid(
						FormRow(
							Flexed(1, Label("Debug Tests")),
							ColSpacer1,
							Flexed(5,
								checkboxFont(
									CheckboxBool(
										o.Clickables.Get("debug-tests"), &o.Test),
								),
							),
						),
					),
					RowSpacer1,
					Rigid(
						FormRow(
							Flexed(1, Label("Run immediately")),
							ColSpacer1,
							Flexed(5,
								checkboxFont(
									CheckboxBool(
										o.Clickables.Get("run-immediatly"), &o.RunImmediately),
								),
							),
						),
					),
					RowSpacer3,
					Rigid(Button(o.Clickables.Get("debug"), "Debug")),
				),
			),
		),
	}

	if o.Error != nil {
		children = append(children,
			RowSpacer3,
			Rigid(Label(o.ErrorString())),
		)
	}

	return Rows(children...)(c)
}

func (o *OpenScreen) SessionsListWidget(c C) D {
	if len(o.Sessions) == 0 {
		return Centered(Label("No recent projects yet."))(c)
	}

	ele := func(c C, i int) D {
		if i >= len(o.Sessions) {
			return D{}
		}

		s := o.Sessions[i]

		del := func() {
			DeleteSession(s.ID)
			o.Sessions = ListSessions()
			op.InvalidateOp{}.Add(c.Ops)
		}

		click := func() {
			debugger, err := FromSession(s)
			if err != nil {
				o.Error = err
				return
			}

			d, err := NewDebugScreen(debugger)
			if err != nil {
				o.Error = err
				return
			}

			tree = d.Layout
			win.Invalidate()
		}

		sessionItem := Columns(
			Rigid(FontEnlarge(2.5)(IconFolder)),
			ColSpacer3,
			Flexed(1,
				Rows(
					Rigid(Bold(Label(s.Path))),
					Rigid(Label(s.BinName+" "+s.Args)),
					Rigid(Label(fmt.Sprintf("Run immediately: %t", s.RunImmediately))),
				),
			),
		)

		delItem := OnClick(o.Clickables.Get(s.ID+"-del"), FontEnlarge(2.5)(IconDelete), del)

		return Inset1(
			Columns(
				Flexed(1, OnClick(o.Clickables.Get(s.ID), sessionItem, click)),
				Rigid(delItem),
			),
		)(c)
	}

	return ZebraList(&o.SessionsList, len(o.Sessions), ele)(c)
}

func (o *OpenScreen) OpenPath() {
	debugger, err := NewDebugger(o.Path, o.BinaryEditor.Text(), o.Arguments.Text(), o.RunImmediately, o.Test)
	if err != nil {
		o.Error = err
		return
	}

	d, err := NewDebugScreen(debugger)
	if err != nil {
		o.Error = err
		return
	}

	s := ToSession(debugger)
	if err := SaveSession(s); err != nil {
		d.Log(LogError, "Can't save session: %s", err)
	}

	tree = d.Layout
}

func (o *OpenScreen) ErrorString() (s string) {
	if o.Error != nil {
		s = o.Error.Error()
	}
	return
}

func NewOpenScreen() (o *OpenScreen) {
	wd, _ := os.Getwd()
	o = &OpenScreen{
		Path:             wd,
		DirectoryBrowser: NewDirectoryBrowser(),
		Sessions:         ListSessions(),
		SessionsList:     NewVerticalList(),
		Clickables:       NewClickables(),
		BinaryEditor:     LineEditor(),
		Arguments:        LineEditor(),
	}

	o.BinaryEditor.SetText("debug")
	return
}
