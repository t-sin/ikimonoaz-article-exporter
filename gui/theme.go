package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.LightTheme().Color(name, theme.VariantLight)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return theme.LightTheme().Font(style)
	}
	if style.Italic {
		return theme.LightTheme().Font(style)
	}
	if style.Bold {
		return resourceRoundedMplus1pBoldTtf
	}
	return resourceRoundedMplus1pRegularTtf
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(name)
}
