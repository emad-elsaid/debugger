package ui

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func Chart[T Numeric](ds []T, height float32) W {
	return func(c C) D {

		max := float32(Max(ds))
		if max == 0 {
			return EmptyWidget(c)
		}

		capds := float32(cap(ds))
		lends := float32(len(ds))
		width := float32(c.Constraints.Max.X)
		center := Pt(width/2, height/2)
		xunit := width / capds
		yunit := height / max

		matOp := op.Affine(
			f32.Affine2D{}.
				Rotate(center, math.Pi).
				Scale(center, Pt(-1, 1)).
				Offset(Pt((capds-lends)*xunit, 0)),
		).Push(c.Ops)

		p := new(clip.Path)
		p.Begin(c.Ops)
		p.MoveTo(Pt(0, 0))

		for i, v := range ds {
			x := xunit * float32(i)
			y := yunit * float32(v)
			p.LineTo(Pt(x, y))
		}

		p.LineTo(
			Pt(xunit*lends, 0),
		)
		p.Close()

		defer clip.Outline{Path: p.End()}.Op().Push(c.Ops).Pop()
		paint.Fill(c.Ops, BLUEBERRY_100)
		matOp.Pop()

		return D{Size: P{X: int(width), Y: int(height)}}
	}
}
