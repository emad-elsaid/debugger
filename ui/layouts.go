package ui

import (
	"math"

	"gioui.org/layout"
)

type (
	FlexChild = layout.FlexChild
)

var (
	Flexed = layout.Flexed
	Rigid  = layout.Rigid
)

func LayoutToWidget(r func(C, W) D, w W) W {
	return func(c C) D {
		return r(c, w)
	}
}

func LayoutToWrapper(r func(C, W) D) func(w W) W {
	return func(w W) W {
		return func(c C) D {
			return r(c, w)
		}
	}
}

func Rows(children ...layout.FlexChild) W {
	return func(c C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(c, children...)
	}
}

var (
	RowSpacer1 = Rigid(HSpacer1)
	RowSpacer2 = Rigid(HSpacer2)
	RowSpacer3 = Rigid(HSpacer3)
	RowSpacer4 = Rigid(HSpacer4)
	RowSpacer5 = Rigid(HSpacer5)
	RowSpacer6 = Rigid(HSpacer6)
)

func Columns(children ...layout.FlexChild) W {
	return func(c C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(c, children...)
	}
}

func ColumnsVCentered(children ...layout.FlexChild) W {
	return func(c C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(c, children...)
	}
}

var (
	ColSpacer1 = Rigid(WSpacer1)
	ColSpacer2 = Rigid(WSpacer2)
	ColSpacer3 = Rigid(WSpacer3)
	ColSpacer4 = Rigid(WSpacer4)
	ColSpacer5 = Rigid(WSpacer5)
	ColSpacer6 = Rigid(WSpacer6)
)

func FormRow(children ...layout.FlexChild) W {
	return func(c C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(c, children...)
	}
}

func Grid(l *List, count int, minItemSize int, ele layout.ListElement) W {
	return func(c C) D {
		perRow := c.Constraints.Max.X / minItemSize

		row := func(c C, i int) D {
			children := []FlexChild{}
			start := i * int(perRow)
			end := (i + 1) * int(perRow)
			for f := start; f < end; f++ {
				if f < count {
					children = append(children, Flexed(1, func(f int) W {
						return func(c C) D {
							return ele(c, f)
						}
					}(f)))
				} else {
					children = append(children, Flexed(1, EmptyWidget))
				}
			}

			return Columns(children...)(c)
		}

		return l.Layout(c, int(math.Ceil(float64(count)/float64(perRow))), row)
	}
}

func Centered(w W) W {
	return func(c C) D {
		v := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}
		h := layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}
		return h.Layout(c, Rigid(func(c C) D {
			return v.Layout(c, Rigid(w))
		}))
	}
}

func Constraint(width, height int, w W) W {
	return func(c C) D {
		c.Constraints.Max = P{width, height}
		return w(c)
	}
}
