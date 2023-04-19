package main

import (
	"fmt"
	"log"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
)

type LogKind int

const (
	LogInfo LogKind = iota
	LogError
	LogCmd
)

func (k LogKind) String() string {
	switch k {
	case LogInfo:
		return "INFO"
	case LogError:
		return "ERR"
	case LogCmd:
		return "CMD"
	default:
		return "UNK"
	}
}

type Log struct {
	Kind LogKind
	Log  string
}

func (log *Log) Layout(c C) D {
	fg := Theme.TextColor
	bg := BackgroundColor

	switch log.Kind {
	case LogInfo:
		bg = BLUEBERRY_100
		fg = BLUEBERRY_700
	case LogError:
		bg = STRAWBERRY_100
		fg = STRAWBERRY_700
	}

	return Columns(
		Rigid(
			Background(bg,
				Inset05(
					AlignMiddle(
						TextColor(fg)(
							Label(log.Kind.String()),
						),
					),
				),
			),
		),
		ColSpacer1,
		Flexed(1, Label(log.Log)),
	)(c)
}

type ConsolePanel struct {
	Editor widget.Editor
	List   List
	Logs   []W
}

func NewConsolePanel() ConsolePanel {
	return ConsolePanel{
		Editor: LineEditor(),
		List:   NewVerticalList(),
		Logs:   []W{},
	}
}

func (cp *ConsolePanel) Layout(d *Debugger) W {
	cp.List.ScrollToEnd = true
	cp.HandleEvent(d)

	logsPanel := ZebraList(&cp.List, len(cp.Logs), func(c C, i int) D {
		return Inset05(cp.Logs[i])(c)
	})

	return Rows(
		Flexed(1, logsPanel),
		Rigid(TextInput(&cp.Editor, "Write commands here")),
	)
}

func (cp *ConsolePanel) HandleEvent(d *Debugger) {
	evts := cp.Editor.Events()
	for _, evt := range evts {
		switch v := evt.(type) {
		case widget.SubmitEvent:
			text := v.Text
			if len(text) == 0 {
				continue
			}
			cp.Run(d, text)
			cp.Editor.SetText("")
		}
	}
}

func (cp *ConsolePanel) Run(d *Debugger, expr string) {
	cp.Log(LogCmd, expr)

	v, err := d.ExecuteExpr(expr)
	if err != nil {
		log.Printf("Executing command: %s resulted in error\n%s", expr, err.Error())
		return
	}

	widget := NewVarWidget(v)
	cp.Logs = append(cp.Logs, widget.Layout)
}

func (cp *ConsolePanel) Log(kind LogKind, log string, params ...interface{}) {
	l := Log{kind, fmt.Sprintf(log, params...)}
	cp.Logs = append(cp.Logs, l.Layout)
}
