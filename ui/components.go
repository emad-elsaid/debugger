package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

func Icon(b []byte, color color.NRGBA) W {
	i, _ := widget.NewIcon(b)

	return func(c C) D {
		c.Constraints.Min.X = int(Theme.FontSize)
		return i.Layout(c, color)
	}
}

func DataSheet(data ...[]string) W {
	rows := make([]layout.FlexChild, 0, len(data))
	for i := range data {
		rows = append(rows, Rigid(dataSheetLine(data[i])))
	}

	return Rows(rows...)
}

func dataSheetLine(datum []string) W {
	row := make([]layout.FlexChild, 0, len(datum))
	for j := range datum {
		row = append(row, Rigid(Label(datum[j])))
	}

	return func(c C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(c, row...)
	}
}

func Card(w W) W {
	return func(c C) D {
		return layout.Inset{Bottom: SpaceUnit}.Layout(c,
			Background(CardColor,
				Inset3(w),
			),
		)
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
