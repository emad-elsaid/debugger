package ui

import (
	"image/color"

	"gioui.org/widget"
)

func Icon(b []byte, color color.NRGBA) W {
	i, _ := widget.NewIcon(b)

	return func(c C) D {
		c.Constraints.Min.X = int(Theme.FontSize)
		return i.Layout(c, color)
	}
}

func Panel(title string, w W) W {
	return Rows(
		Rigid(
			Inset1(
				Bold(
					Label(title),
				),
			),
		),
		Rigid(w),
	)
}
