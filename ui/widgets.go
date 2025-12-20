package ui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type (
	W       = layout.Widget
	C       = layout.Context
	D       = layout.Dimensions
	P       = image.Point
	DP      = unit.Dp
	SP      = unit.Sp
	Wrapper = func(W) W
	List    = layout.List
)

var (
	Pt            = f32.Pt
	SpaceUnit  DP = 8
	BorderSize DP = 1

	th = material.NewTheme()
)

func init() {
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.TextSize = Theme.FontSize
	th.Fg = Theme.TextColor
}

func EmptyWidget(c C) D { return D{} }

func Wrap(w W, wrappers ...Wrapper) W {
	for i := len(wrappers) - 1; i >= 0; i-- {
		w = wrappers[i](w)
	}

	return w
}
