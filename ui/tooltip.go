package ui

import (
	"image"

	"gioui.org/op"
)

var (
	TooltipFgColor = WHITE
	TooltipBgColor = Alpha(BLACK_500, 242)
)

func Tooltip(attachTo W, s string) W {
	return func(c C) D {
		attachToD := attachTo(c)
		tooltipMacro := op.Record(c.Ops)
		d := RoundedCorners(
			Background(TooltipBgColor,
				Inset1(
					TextColor(TooltipFgColor)(Label(s)),
				),
			),
		)(c)
		tooltipOp := tooltipMacro.Stop()

		macro := op.Record(c.Ops)
		trans := op.Offset(image.Pt(attachToD.Size.X/2-d.Size.X/2, attachToD.Size.Y+int(SpaceUnit/2))).Push(c.Ops)
		tooltipOp.Add(c.Ops)
		trans.Pop()
		op.Defer(c.Ops, macro.Stop())

		return attachToD
	}
}
