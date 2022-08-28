package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// Insets
var (
	Inset05 = LayoutToWrapper(layout.UniformInset(SpaceUnit * 0.5).Layout)
	Inset1  = LayoutToWrapper(layout.UniformInset(SpaceUnit).Layout)
	Inset2  = LayoutToWrapper(layout.UniformInset(SpaceUnit * 2).Layout)
	Inset3  = LayoutToWrapper(layout.UniformInset(SpaceUnit * 3).Layout)
	Inset4  = LayoutToWrapper(layout.UniformInset(SpaceUnit * 4).Layout)
	Inset5  = LayoutToWrapper(layout.UniformInset(SpaceUnit * 5).Layout)
	Inset6  = LayoutToWrapper(layout.UniformInset(SpaceUnit * 6).Layout)
)

// Spaces
var (
	WSpacer1 = layout.Spacer{Width: SpaceUnit}.Layout
	WSpacer2 = layout.Spacer{Width: SpaceUnit * 2}.Layout
	WSpacer3 = layout.Spacer{Width: SpaceUnit * 3}.Layout
	WSpacer4 = layout.Spacer{Width: SpaceUnit * 4}.Layout
	WSpacer5 = layout.Spacer{Width: SpaceUnit * 5}.Layout
	WSpacer6 = layout.Spacer{Width: SpaceUnit * 6}.Layout
)

var (
	HSpacer1 = layout.Spacer{Height: SpaceUnit}.Layout
	HSpacer2 = layout.Spacer{Height: SpaceUnit * 2}.Layout
	HSpacer3 = layout.Spacer{Height: SpaceUnit * 3}.Layout
	HSpacer4 = layout.Spacer{Height: SpaceUnit * 4}.Layout
	HSpacer5 = layout.Spacer{Height: SpaceUnit * 5}.Layout
	HSpacer6 = layout.Spacer{Height: SpaceUnit * 6}.Layout
)

func Margin(t, r, b, l unit.Dp) Wrapper {
	return LayoutToWrapper(layout.Inset{Top: t, Right: r, Bottom: b, Left: l}.Layout)
}
