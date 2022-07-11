package ui

var (
	ToolbarBgColor = SILVER_100
)

func ToolbarButton(clickable *Clickable, icon W, desc string) W {
	bg := ToolbarBgColor
	hovered := clickable.Hovered()
	if hovered {
		bg = MixColor(ToolbarBgColor, BLACK_900, 80)
	}

	btn := func(c C) D {
		return clickable.Layout(c,
			Background(bg,
				Inset1(icon),
			),
		)
	}

	if hovered {
		btn = Tooltip(btn, desc)
	}

	return btn
}
