package ui

import "image/color"

func Color(hex int) color.NRGBA {
	return color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 255,
	}
}

func Alpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = alpha
	return c
}

func MixColor(c1, c2 color.NRGBA, percent int) color.NRGBA {
	p1 := float32(percent) / float32(100.0)
	p2 := 1 - p1
	return color.NRGBA{
		R: uint8(float32(c1.R)*p1 + float32(c2.R)*p2),
		G: uint8(float32(c1.G)*p1 + float32(c2.G)*p2),
		B: uint8(float32(c1.B)*p1 + float32(c2.B)*p2),
		A: uint8(float32(c1.A)*p1 + float32(c2.A)*p2),
	}
}

// ElementaryOS palette https://github.com/elementary/stylesheet/blob/master/src/gtk-4.0/_palette.scss
var (
	STRAWBERRY_100 = Color(0xff8c82)
	STRAWBERRY_300 = Color(0xed5353)
	STRAWBERRY_500 = Color(0xc6262e)
	STRAWBERRY_700 = Color(0xa10705)
	STRAWBERRY_900 = Color(0x7a0000)

	ORANGE_100 = Color(0xffc27d)
	ORANGE_300 = Color(0xffa154)
	ORANGE_500 = Color(0xf37329)
	ORANGE_700 = Color(0xcc3b02)
	ORANGE_900 = Color(0xa62100)

	BANANA_100 = Color(0xfff394)
	BANANA_300 = Color(0xffe16b)
	BANANA_500 = Color(0xf9c440)
	BANANA_700 = Color(0xd48e15)
	BANANA_900 = Color(0xad5f00)

	LIME_100 = Color(0xd1ff82)
	LIME_300 = Color(0x9bdb4d)
	LIME_500 = Color(0x68b723)
	LIME_700 = Color(0x3a9104)
	LIME_900 = Color(0x206b00)

	MINT_100 = Color(0x89ffdd)
	MINT_300 = Color(0x43d6b5)
	MINT_500 = Color(0x28bca3)
	MINT_700 = Color(0x0e9a83)
	MINT_900 = Color(0x007367)

	BLUEBERRY_100 = Color(0x8cd5ff)
	BLUEBERRY_300 = Color(0x64baff)
	BLUEBERRY_500 = Color(0x3689e6)
	BLUEBERRY_700 = Color(0x0d52bf)
	BLUEBERRY_900 = Color(0x002e99)

	BUBBLEGUM_100 = Color(0xfe9ab8)
	BUBBLEGUM_300 = Color(0xf4679d)
	BUBBLEGUM_500 = Color(0xde3e80)
	BUBBLEGUM_700 = Color(0xbc245d)
	BUBBLEGUM_900 = Color(0x910e38)

	GRAPE_100 = Color(0xe4c6fa)
	GRAPE_300 = Color(0xcd9ef7)
	GRAPE_500 = Color(0xa56de2)
	GRAPE_700 = Color(0x7239b3)
	GRAPE_900 = Color(0x452981)

	COCOA_100 = Color(0xa3907c)
	COCOA_300 = Color(0x8a715e)
	COCOA_500 = Color(0x715344)
	COCOA_700 = Color(0x57392d)
	COCOA_900 = Color(0x3d211b)

	SILVER_100 = Color(0xfafafa)
	SILVER_300 = Color(0xd4d4d4)
	SILVER_500 = Color(0xabacae)
	SILVER_700 = Color(0x7e8087)
	SILVER_900 = Color(0x555761)

	SLATE_100 = Color(0x95a3ab)
	SLATE_300 = Color(0x667885)
	SLATE_500 = Color(0x485a6c)
	SLATE_700 = Color(0x273445)
	SLATE_900 = Color(0x0e141f)

	BLACK_100 = Color(0x666666)
	BLACK_300 = Color(0x4d4d4d)
	BLACK_500 = Color(0x333333)
	BLACK_700 = Color(0x1a1a1a)
	BLACK_900 = Color(0x000000)

	WHITE = Color(0xffffff)
)

var (
	ACCENT_COLOR_100 = BLUEBERRY_100
	ACCENT_COLOR_300 = BLUEBERRY_300
	ACCENT_COLOR_500 = BLUEBERRY_500
	ACCENT_COLOR_700 = BLUEBERRY_700
	ACCENT_COLOR_900 = BLUEBERRY_900
	ACCENT_COLOR     = MixColor(BLUEBERRY_300, BLUEBERRY_500, 25)
)

var (
	BackgroundColor    = WHITE
	SecondaryTextColor = SILVER_500

	DangerColor  = STRAWBERRY_500
	SuccessColor = LIME_700
	WarningColor = BANANA_900

	InputBgColor   = SILVER_100
	ViewsBgColor   = WHITE
	SidebarBgColor = MixColor(SILVER_100, SILVER_300, 75)

	BorderColor       = MixColor(SILVER_300, WHITE, 70)
	ActiveBorderColor = SILVER_500

	CardColor     = SidebarBgColor
	CheckboxColor = ACCENT_COLOR_500
)
