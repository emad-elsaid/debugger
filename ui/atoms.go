package ui

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func RoundedCorners(w W) W {
	return func(c C) D {
		macro := op.Record(c.Ops)
		d := w(c)
		macroOp := macro.Stop()

		defer clip.UniformRRect(image.Rect(0, 0, d.Size.X, d.Size.Y), 5).Push(c.Ops).Pop()
		macroOp.Add(c.Ops)
		return d
	}
}

func Border(w W) W {
	return func(c C) D {
		return widget.Border{Color: BorderColor, Width: BorderSize, CornerRadius: 0}.Layout(c, w)
	}
}

func BorderActive(w W) W {
	return func(c C) D {
		return widget.Border{Color: ActiveBorderColor, Width: BorderSize, CornerRadius: 0}.Layout(c, w)
	}
}

func Outline(width int, col color.NRGBA, w W) W {
	return func(c C) (d D) {
		wMacro := op.Record(c.Ops)
		d = w(c)
		wOp := wMacro.Stop()

		sz := d.Size
		width := c.Dp(DP(width))
		sz.X += width
		sz.Y += width
		r := image.Rectangle{Max: sz}
		op.Offset(P{-width / 2, -width / 2}).Add(c.Ops)

		paint.FillShape(c.Ops,
			col,
			clip.Stroke{
				Path:  clip.Rect(r).Path(),
				Width: float32(width),
			}.Op(),
		)

		op.Offset(P{width / 2, width / 2}).Add(c.Ops)
		wOp.Add(c.Ops)

		return
	}
}

func Background(background color.NRGBA, w W) W {
	return func(c C) D {
		macro := op.Record(c.Ops)
		d := w(c)
		path := macro.Stop()

		cl := clip.Rect{Max: d.Size}.Push(c.Ops)
		paint.Fill(c.Ops, background)
		cl.Pop()

		path.Add(c.Ops)
		return d
	}
}

func HR(sz int) W {
	return func(c C) D {
		cl := clip.Path{}
		cl.Begin(c.Ops)
		cl.MoveTo(Pt(0, 0))
		cl.Line(Pt(float32(c.Constraints.Max.X), 0))

		defer clip.Stroke{
			Path:  cl.End(),
			Width: float32(sz),
		}.Op().Push(c.Ops).Pop()

		paint.Fill(c.Ops, BorderColor)

		return D{Size: P{c.Constraints.Min.X, sz}}
	}
}

func VR(sz int) W {
	return func(c C) D {
		cl := clip.Path{}
		cl.Begin(c.Ops)
		cl.MoveTo(Pt(0, 0))
		cl.Line(Pt(0, float32(c.Constraints.Max.Y)))

		defer clip.Stroke{
			Path:  cl.End(),
			Width: float32(sz),
		}.Op().Push(c.Ops).Pop()

		paint.Fill(c.Ops, BorderColor)

		return D{Size: P{sz, c.Constraints.Min.Y}}
	}
}
