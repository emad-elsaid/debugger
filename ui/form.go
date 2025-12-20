package ui

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func LineEditor() widget.Editor {
	return widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
}

func TextInput(editor *widget.Editor, hint string) W {
	return func(c C) D {
		border := Border

		if c.Source.Focused(editor) {
			border = BorderActive
		}

		return Background(BackgroundColor,
			Wrap(
				material.Editor(th, editor, hint).Layout,
				border, Inset1,
			),
		)(c)
	}
}

func Button(clickable *Clickable, label string) W {
	col := ACCENT_COLOR
	if clickable.Hovered() {
		col = ACCENT_COLOR_500
	}

	return LayoutToWidget(
		clickable.Layout,
		RoundedCorners(
			Background(col, Wrap(Label(label), Inset1, AlignMiddle, TextColor(WHITE))),
		),
	)
}

var (
	IconCheckbox       = Inset05(Icon(icons.ToggleCheckBoxOutlineBlank, CheckboxColor))
	IconCheckboxActive = Inset05(Icon(icons.ToggleCheckBox, CheckboxColor))
)

func CheckboxBtn(value bool, btn *Clickable) W {
	var icon W
	if value {
		icon = IconCheckboxActive
	} else {
		icon = IconCheckbox
	}

	return LayoutToWidget(btn.Layout, icon)
}

func Checkbox(btn *Clickable, value bool, onclick func()) W {
	return func(c C) D {
		if btn.Clicked(c) {
			onclick()
		}

		var icon W
		if value {
			icon = IconCheckboxActive
		} else {
			icon = IconCheckbox
		}

		return btn.Layout(c, icon)
	}
}

func CheckboxBool(btn *Clickable, value *bool) W {
	return Checkbox(btn, *value, func() { *value = !*value })
}

func OnClick(btn *Clickable, w W, onclick func()) W {
	return func(c C) D {
		if btn.Clicked(c) {
			onclick()
		}

		return btn.Layout(c, w)
	}
}
