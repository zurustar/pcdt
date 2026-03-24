// Package view contains the UI components for the countdown timer.
package view

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/unit"
)

// AlertTheme defines the visual theme for the countdown timer.
// It uses a red and black color scheme inspired by the warning displays.
type AlertTheme struct {
	// Background colors
	BackgroundColor color.NRGBA

	// Primary colors (alert red)
	PrimaryColor   color.NRGBA
	SecondaryColor color.NRGBA

	// Text colors for different states
	TextColorNormal   color.NRGBA // Yellow-green for positive time
	TextColorNegative color.NRGBA // Red for negative time

	// Accent colors
	AccentColor  color.NRGBA
	WarningColor color.NRGBA

	// Button colors
	ButtonBackground       color.NRGBA
	ButtonBackgroundHover  color.NRGBA
	ButtonText             color.NRGBA
	ButtonDisabled         color.NRGBA
	ButtonDisabledText     color.NRGBA

	// Input field colors
	InputBackground color.NRGBA
	InputBorder     color.NRGBA
	InputText       color.NRGBA
	InputError      color.NRGBA

	// Font settings
	FontFace       font.Font
	FontSizeSmall  unit.Sp
	FontSizeMedium unit.Sp
	FontSizeLarge  unit.Sp
	FontSizeTimer  unit.Sp

	// Segment display settings
	SegmentFontWeight font.Weight
}

// NewAlertTheme creates a new theme with alert-style colors.
// The theme uses yellow-green for normal countdown and red for overtime.
func NewAlertTheme() *AlertTheme {
	// Normal color: #EAD700 (yellow-green)
	normalColor := color.NRGBA{R: 234, G: 215, B: 0, A: 255}
	// Negative color: bright red
	negativeColor := color.NRGBA{R: 255, G: 50, B: 50, A: 255}

	return &AlertTheme{
		// Background: Deep black
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},

		// Primary colors match the timer color
		PrimaryColor:   normalColor,
		SecondaryColor: DimColor(normalColor, 0.7),

		// Text colors for different states
		TextColorNormal:   normalColor,
		TextColorNegative: negativeColor,

		// Accent colors
		AccentColor:  normalColor,
		WarningColor: color.NRGBA{R: 255, G: 200, B: 0, A: 255},

		// Button colors - match timer color
		ButtonBackground:      DimColor(normalColor, 0.3),
		ButtonBackgroundHover: DimColor(normalColor, 0.5),
		ButtonText:            normalColor,
		ButtonDisabled:        color.NRGBA{R: 80, G: 80, B: 80, A: 255},
		ButtonDisabledText:    color.NRGBA{R: 120, G: 120, B: 120, A: 255},

		// Input field colors - match timer color
		InputBackground: color.NRGBA{R: 20, G: 20, B: 20, A: 255},
		InputBorder:     normalColor,
		InputText:       normalColor,
		InputError:      negativeColor,

		// Font settings for segment display style
		FontFace: font.Font{
			Typeface: "monospace",
		},
		FontSizeSmall:  unit.Sp(14),
		FontSizeMedium: unit.Sp(18),
		FontSizeLarge:  unit.Sp(24),
		FontSizeTimer:  unit.Sp(72),

		// Bold weight for segment display emphasis
		SegmentFontWeight: font.Bold,
	}
}

// GetTextColorForState returns the appropriate text color based on whether time is negative.
func (t *AlertTheme) GetTextColorForState(isNegative bool) color.NRGBA {
	if isNegative {
		return t.TextColorNegative
	}
	return t.TextColorNormal
}

// GetBlinkColor returns a color that alternates between visible and dimmed
// based on the blink state. Used for implementing the blink animation effect.
func (t *AlertTheme) GetBlinkColor(baseColor color.NRGBA, visible bool) color.NRGBA {
	if visible {
		return baseColor
	}
	// Return a dimmed version of the color when not visible
	return color.NRGBA{
		R: baseColor.R / 4,
		G: baseColor.G / 4,
		B: baseColor.B / 4,
		A: baseColor.A,
	}
}

// GetTimerFont returns the font configuration for the main timer display.
func (t *AlertTheme) GetTimerFont() font.Font {
	return font.Font{
		Typeface: t.FontFace.Typeface,
		Weight:   t.SegmentFontWeight,
	}
}

// DimColor returns a dimmed version of the given color.
// The factor should be between 0.0 (fully dimmed) and 1.0 (no change).
func DimColor(c color.NRGBA, factor float64) color.NRGBA {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return color.NRGBA{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: c.A,
	}
}

// LerpColor linearly interpolates between two colors.
// The t parameter should be between 0.0 (returns a) and 1.0 (returns b).
func LerpColor(a, b color.NRGBA, t float64) color.NRGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return color.NRGBA{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: uint8(float64(a.A) + (float64(b.A)-float64(a.A))*t),
	}
}
