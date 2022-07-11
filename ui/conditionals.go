package ui

func WidgetIf(cond bool, w W) W {
	if cond {
		return w
	} else {
		return EmptyWidget
	}
}
