package ui

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
)

type ThemeStyle struct {
	FontSize      SP
	FontFamily    *text.Shaper
	FontWeight    font.Weight
	TextAlignment text.Alignment
	TextColor     color.NRGBA
	MaxLines      int
}

var Theme = ThemeStyle{
	FontSize:      13,
	FontFamily:    fontShaper,
	FontWeight:    font.Normal,
	TextAlignment: text.Start,
	TextColor:     BLACK_500,
	MaxLines:      0,
}

func FontSize(s SP) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.FontSize
			Theme.FontSize = s
			d := w(c)
			Theme.FontSize = old
			return d
		}
	}
}

func FontEnlarge(s float32) Wrapper {
	return FontSize(SP(s) * Theme.FontSize)
}

func Font(f *text.Shaper) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.FontFamily
			Theme.FontFamily = f
			d := w(c)
			Theme.FontFamily = old
			return d
		}
	}
}

func FontWeight(f font.Weight) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.FontWeight
			Theme.FontWeight = f
			d := w(c)
			Theme.FontWeight = old
			return d
		}
	}
}

func TextAlignment(a text.Alignment) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.TextAlignment
			Theme.TextAlignment = a
			d := w(c)
			Theme.TextAlignment = old
			return d
		}
	}
}

var AlignStart = TextAlignment(text.Start)
var AlignMiddle = TextAlignment(text.Middle)
var AlignEnd = TextAlignment(text.End)

func TextColor(col color.NRGBA) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.TextColor
			Theme.TextColor = col
			d := w(c)
			Theme.TextColor = old
			return d
		}
	}
}

func MaxLines(i int) Wrapper {
	return func(w W) W {
		return func(c C) D {
			old := Theme.MaxLines
			Theme.MaxLines = i
			d := w(c)
			Theme.MaxLines = old
			return d
		}
	}
}
