package main

import (
	"fmt"
	"os"

	. "github.com/emad-elsaid/debugger/ui"
	"github.com/emad-elsaid/delve/pkg/proc"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	IconStop     = Icon(icons.AVStop, Theme.TextColor)
	IconRestart  = Icon(icons.NavigationRefresh, Theme.TextColor)
	IconContinue = Icon(icons.AVPlayArrow, SuccessColor)
	IconStepOut  = Icon(icons.NavigationArrowBack, Theme.TextColor)
	IconNext     = Icon(icons.NavigationArrowDownward, Theme.TextColor)
	IconStep     = Icon(icons.NavigationArrowForward, Theme.TextColor)
	IconProcess  = Icon(icons.AVPlayCircleFilled, SecondaryTextColor)
)

type Toolbar struct {
	RestartBtn  Clickable
	StopBtn     Clickable
	ContinueBtn Clickable
	NextBtn     Clickable
	StepBtn     Clickable
	StepOutBtn  Clickable
}

func (t *Toolbar) Layout(d *Debugger) W {
	target := d.Target()
	pid := target.Pid()

	valid := d.LastState == nil || !d.LastState.Exited
	showStop := d.IsRunning()
	showCont := !d.IsRunning() && valid
	showRestart := showStop || showCont || !valid

	showStep := !d.IsRunning() && valid && d.StopReason() != proc.StopLaunched
	showNext := !d.IsRunning() && valid && d.StopReason() != proc.StopLaunched
	showStepOut := !d.IsRunning() && valid && d.StopReason() != proc.StopLaunched

	IconSize := FontEnlarge(2)

	btns := []FlexChild{
		Flexed(1,
			Inset1(
				Columns(
					Rigid(IconSize(IconFolder)),
					ColSpacer1,
					Rigid(
						Rows(
							Rigid(Bold(Label(d.Path))),
							Rigid(
								WidgetIf(
									d.File != "",
									Wrap(
										Text(fmt.Sprintf("%s:%d", d.File, d.Line)),
										AlignStart,
										TextColor(SecondaryTextColor),
										MaxLines(3),
									),
								),
							),
						),
					),
					ColSpacer3,
					Rigid(IconSize(IconProcess)),
					ColSpacer1,
					Rigid(
						Rows(
							Rigid(Label(fmt.Sprintf("pid: %d", pid))),
							Rigid(Label(fmt.Sprintf("cwd: %s", t.CWD(pid)))),
						),
					),
				),
			),
		),
	}

	if t.StopBtn.Clicked() {
		d.Stop()
	}
	if showStop {
		btns = append(btns,
			Rigid(ToolbarButton(&t.StopBtn, IconSize(IconStop), "Stop")),
		)
	}

	if t.RestartBtn.Clicked() {
		d.Restart()
	}
	if showRestart {
		btns = append(btns,
			Rigid(ToolbarButton(&t.RestartBtn, IconSize(IconRestart), "Restart")),
		)
	}

	if t.ContinueBtn.Clicked() {
		d.Continue()
	}
	if showCont {
		btns = append(btns,
			Rigid(ToolbarButton(&t.ContinueBtn, IconSize(IconContinue), "Continue")),
		)
	}

	if t.StepOutBtn.Clicked() {
		d.StepOut()
	}
	if showStepOut {
		btns = append(btns,
			Rigid(ToolbarButton(&t.StepOutBtn, IconSize(IconStepOut), "Step out")),
		)
	}

	if t.NextBtn.Clicked() {
		d.Next()
	}
	if showNext {
		btns = append(btns,
			Rigid(ToolbarButton(&t.NextBtn, IconSize(IconNext), "Next")),
		)
	}

	if t.StepBtn.Clicked() {
		d.Step()
	}
	if showStep {
		btns = append(btns,
			Rigid(ToolbarButton(&t.StepBtn, IconSize(IconStep), "Step")),
		)
	}

	return Background(ToolbarBgColor,
		ColumnsVCentered(btns...),
	)
}

func (t *Toolbar) CWD(pid int) string {
	cwd, _ := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	return string(cwd)
}
