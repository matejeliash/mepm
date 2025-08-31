package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type mobileTheme struct{}

var _ fyne.Theme = (*mobileTheme)(nil)

func (m mobileTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m mobileTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m mobileTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m mobileTheme) Size(name fyne.ThemeSizeName) float32 {
	// Increase all default sizes by a multiplier (e.g., 1.5x)
	defaultSize := theme.DefaultTheme().Size(name)
	return defaultSize * 1.5
}
