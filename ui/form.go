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
	border := Border

	if editor.Focused() {
		border = BorderActive
	}

	return Background(BackgroundColor,
		Wrap(
			material.Editor(th, editor, hint).Layout,
			border, Inset1,
		),
	)
}

func Button(clickable *Clickable, label string) W {
	col := ACCENT_COLOR
	if clickable.Hovered() {
		col = ACCENT_COLOR_500
	}

	return LayoutToWidget(clickable.Layout)(
		RoundedCorners(
			Background(col, Wrap(Label(label), Inset1, AlignMiddle, TextColor(WHITE))),
		),
	)
}

var (
	IconCheckbox       = Inset05(Icon(icons.ToggleCheckBoxOutlineBlank, CheckboxColor))
	IconCheckboxActive = Inset05(Icon(icons.ToggleCheckBox, CheckboxColor))
	checkboxBtns       = map[string]*Clickable{}
)

func CheckboxBtn(value bool, btn *Clickable) W {
	var icon W
	if value {
		icon = IconCheckboxActive
	} else {
		icon = IconCheckbox
	}

	return LayoutToWidget(btn.Layout)(icon)
}

func Checkbox(id string, value bool, onclick func()) W {
	var btn *Clickable
	var ok bool
	if btn, ok = checkboxBtns[id]; !ok {
		btn = new(Clickable)
		checkboxBtns[id] = btn
	}

	if btn.Clicked() {
		onclick()
	}

	var icon W
	if value {
		icon = IconCheckboxActive
	} else {
		icon = IconCheckbox
	}

	return LayoutToWidget(btn.Layout)(icon)
}

func CheckboxBool(id string, value *bool) W {
	return Checkbox(id, *value, func() { *value = !*value })
}

var onClickBtns = map[string]*Clickable{}

func OnClick(id string, w W, onclick func()) W {
	var btn *Clickable
	var ok bool
	if btn, ok = onClickBtns[id]; !ok {
		btn = new(Clickable)
		onClickBtns[id] = btn
	}

	if btn.Clicked() {
		onclick()
	}

	return func(c C) D {
		return btn.Layout(c, w)
	}
}
