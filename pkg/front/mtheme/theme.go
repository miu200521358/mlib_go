package mtheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MTheme struct{}

var _ fyne.Theme = (*MTheme)(nil)

// return bundled font resource
func (*MTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return resourceResourcesMPLUS1RegularTtf
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return resourceResourcesMPLUS1RegularTtf
}

func (*MTheme) Color(name fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorRed:
		// #C03221
		return color.RGBA{R: 0xc0, G: 0x32, B: 0x21, A: 0xff}
	case theme.ColorOrange:
		// #F16A1B
		return color.RGBA{R: 0xf1, G: 0x6a, B: 0x1b, A: 0xff}
	case theme.ColorYellow:
		// #f66701
		return color.RGBA{R: 0xf6, G: 0x67, B: 0x01, A: 0xff}
	case theme.ColorGreen:
		// #1AA053
		return color.RGBA{R: 0x1a, G: 0xa0, B: 0x53, A: 0xff}
	case theme.ColorBlue:
		// #3A57E8
		return color.RGBA{R: 0x3a, G: 0x57, B: 0xe8, A: 0xff}
	case theme.ColorPurple:
		// #000b2a
		return color.RGBA{R: 0x00, G: 0x0b, B: 0x2a, A: 0xff}
	case theme.ColorBrown:
		// #020528
		return color.RGBA{R: 0x02, G: 0x05, B: 0x28, A: 0xff}
	case theme.ColorGray:
		// #6C757D
		return color.RGBA{R: 0x6c, G: 0x75, B: 0x7d, A: 0xff}
	case theme.ColorNameBackground:
		// #dee2e6
		return color.RGBA{R: 0xde, G: 0xe2, B: 0xe6, A: 0xff}
	case theme.ColorNameButton:
		// #8E9BA5
		return color.RGBA{R: 0x8e, G: 0x9b, B: 0xa5, A: 0xff}
	case theme.ColorNameDisabledButton:
		// #E9ECEF
		return color.RGBA{R: 0xe2, G: 0xe3, B: 0xe5, A: 0xff}
	case theme.ColorNameDisabled:
		// #8A92A6
		return color.RGBA{R: 0x8A, G: 0x92, B: 0xA6, A: 0xff}
	case theme.ColorNameError:
		// #C03221
		return color.RGBA{R: 0xc0, G: 0x32, B: 0x21, A: 0xff}
	case theme.ColorNameFocus:
		// #079AA233
		return color.RGBA{R: 0x07, G: 0x9a, B: 0xa2, A: 0x33}
	case theme.ColorNameForeground:
		// #232d42
		return color.RGBA{R: 0x23, G: 0x2d, B: 0x42, A: 0xff}
	case theme.ColorNameHeaderBackground:
		// #8A92A6
		return color.RGBA{R: 0x8a, G: 0x92, B: 0xa6, A: 0xff}
	case theme.ColorNameHover:
		// #6C757D
		return color.RGBA{R: 0x6c, G: 0x75, B: 0x7d, A: 0xff}
	case theme.ColorNameHyperlink:
		// #2E46BA
		return color.RGBA{R: 0x2e, G: 0x46, B: 0xba, A: 0xff}
	case theme.ColorNameInputBackground:
		// #F8F9FA
		return color.RGBA{R: 0xf8, G: 0xf9, B: 0xfa, A: 0xff}
	case theme.ColorNameInputBorder:
		// #85888A
		return color.RGBA{R: 0x85, G: 0x88, B: 0x8a, A: 0xff}
	case theme.ColorNameMenuBackground:
		// #DEDFDF
		return color.RGBA{R: 0xde, G: 0xdf, B: 0xdf, A: 0xff}
	case theme.ColorNameOverlayBackground:
		// #8A92A6
		return color.RGBA{R: 0xde, G: 0xdf, B: 0xdf, A: 0xff}
	case theme.ColorNamePlaceHolder:
		// #8A92A6
		return color.RGBA{R: 0x8a, G: 0x92, B: 0xa6, A: 0xff}
	case theme.ColorNamePressed:
		// #212529
		return color.RGBA{R: 0x21, G: 0x25, B: 0x29, A: 0xff}
	case theme.ColorNamePrimary:
		// #3a57e8
		return color.RGBA{R: 0x3a, G: 0x57, B: 0xe8, A: 0xff}
	case theme.ColorNameScrollBar:
		// #C4C4C4
		return color.RGBA{R: 0xc4, G: 0xc4, B: 0xc4, A: 0xff}
	case theme.ColorNameSelection:
		// #079AA2
		return color.RGBA{R: 0x07, G: 0x9a, B: 0xa2, A: 0xff}
	case theme.ColorNameSeparator:
		// #C4C4C4
		return color.RGBA{R: 0xc4, G: 0xc4, B: 0xc4, A: 0xff}
	case theme.ColorNameShadow:
		// #232D42
		return color.RGBA{R: 0x23, G: 0x2d, B: 0x42, A: 0x99}
	case theme.ColorNameSuccess:
		// #1AA053
		return color.RGBA{R: 0x1a, G: 0xa0, B: 0x53, A: 0xff}
	case theme.ColorNameWarning:
		// #F16A1B
		return color.RGBA{R: 0xF1, G: 0x6a, B: 0x1b, A: 0xff}
	}
	// #232d42 (Foreground)
	return color.RGBA{R: 0x23, G: 0x2d, B: 0x42, A: 0xff}
}

func (*MTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
