// Package view contains the UI components for the countdown timer.
package view

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// HelpDialog represents the help dialog that displays keyboard shortcuts.
type HelpDialog struct {
	// Theme is the alert visual theme
	Theme *AlertTheme

	// visible indicates whether the dialog is currently shown
	visible bool

	// closeButton for dismissing the dialog
	closeButton widget.Clickable
}

// NewHelpDialog creates a new help dialog with the given theme.
func NewHelpDialog(theme *AlertTheme) *HelpDialog {
	return &HelpDialog{
		Theme:   theme,
		visible: false,
	}
}

// Show displays the help dialog.
func (h *HelpDialog) Show() {
	h.visible = true
}

// Hide hides the help dialog.
func (h *HelpDialog) Hide() {
	h.visible = false
}

// Toggle toggles the visibility of the help dialog.
func (h *HelpDialog) Toggle() {
	h.visible = !h.visible
}

// IsVisible returns whether the help dialog is currently visible.
func (h *HelpDialog) IsVisible() bool {
	return h.visible
}

// shortcutEntry represents a single keyboard shortcut entry.
type shortcutEntry struct {
	key         string
	description string
}

// getShortcuts returns the list of keyboard shortcuts to display.
func getShortcuts() []shortcutEntry {
	return []shortcutEntry{
		{key: "Enter", description: "カウントダウン開始（入力画面）"},
		{key: "Space", description: "一時停止 / 再開（カウントダウン画面）"},
		{key: "Escape", description: "入力画面に戻る"},
		{key: "H / ?", description: "ヘルプを表示"},
	}
}

// Layout renders the help dialog as an overlay.
func (h *HelpDialog) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !h.visible {
		return layout.Dimensions{}
	}

	// Handle keyboard events to close dialog
	for {
		event, ok := gtx.Event(
			key.Filter{Name: key.NameEscape},
			key.Filter{Name: "H"},
			key.Filter{Name: "?"},
		)
		if !ok {
			break
		}
		if e, ok := event.(key.Event); ok && e.State == key.Press {
			h.Hide()
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}
	}

	// Check for close button click
	if h.closeButton.Clicked(gtx) {
		h.Hide()
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}

	// Draw semi-transparent overlay background
	overlayColor := color.NRGBA{R: 0, G: 0, B: 0, A: 200}
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, overlayColor, clip.Rect(rect).Op())

	// Center the dialog
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return h.layoutDialogBox(gtx, th)
	})
}

// layoutDialogBox renders the dialog box content.
func (h *HelpDialog) layoutDialogBox(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Dialog dimensions
	dialogWidth := unit.Dp(400)
	dialogPadding := unit.Dp(24)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Set minimum width
			gtx.Constraints.Min.X = gtx.Dp(dialogWidth)

			return layout.Inset{
				Top:    dialogPadding,
				Bottom: dialogPadding,
				Left:   dialogPadding,
				Right:  dialogPadding,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Title
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return h.layoutTitle(gtx, th)
					}),
					// Spacer
					layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
					// Shortcuts list
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return h.layoutShortcutsList(gtx, th)
					}),
					// Spacer
					layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
					// Close button
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return h.layoutCloseButton(gtx, th)
					}),
				)
			})
		}),
		// Background with border (drawn behind content)
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return h.layoutDialogBackground(gtx)
		}),
	)
}

// layoutDialogBackground renders the dialog background with border.
func (h *HelpDialog) layoutDialogBackground(gtx layout.Context) layout.Dimensions {
	// Background color
	bgColor := color.NRGBA{R: 20, G: 20, B: 20, A: 255}
	rect := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())

	// Draw border
	borderWidth := gtx.Dp(unit.Dp(3))
	drawBorder(gtx.Ops, rect, h.Theme.PrimaryColor, borderWidth)

	return layout.Dimensions{Size: gtx.Constraints.Min}
}

// layoutTitle renders the dialog title.
func (h *HelpDialog) layoutTitle(gtx layout.Context, th *material.Theme) layout.Dimensions {
	title := material.H5(th, "キーボードショートカット")
	title.Color = h.Theme.PrimaryColor
	title.Font.Weight = font.Bold
	title.Alignment = text.Middle
	return title.Layout(gtx)
}

// layoutShortcutsList renders the list of keyboard shortcuts.
func (h *HelpDialog) layoutShortcutsList(gtx layout.Context, th *material.Theme) layout.Dimensions {
	shortcuts := getShortcuts()

	children := make([]layout.FlexChild, 0, len(shortcuts)*2-1)
	for i, shortcut := range shortcuts {
		s := shortcut // capture for closure
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.layoutShortcutEntry(gtx, th, s)
		}))
		// Add spacer between entries (but not after the last one)
		if i < len(shortcuts)-1 {
			children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout))
		}
	}

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Start,
	}.Layout(gtx, children...)
}

// layoutShortcutEntry renders a single shortcut entry.
func (h *HelpDialog) layoutShortcutEntry(gtx layout.Context, th *material.Theme, entry shortcutEntry) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Key badge
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.layoutKeyBadge(gtx, th, entry.key)
		}),
		// Spacer
		layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
		// Description
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			desc := material.Body1(th, entry.description)
			desc.Color = h.Theme.TextColorNormal
			return desc.Layout(gtx)
		}),
	)
}

// layoutKeyBadge renders a key badge (styled key indicator).
func (h *HelpDialog) layoutKeyBadge(gtx layout.Context, th *material.Theme, keyText string) layout.Dimensions {
	// Badge dimensions
	minWidth := unit.Dp(80)
	padding := unit.Dp(8)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    padding,
				Bottom: padding,
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Ensure minimum width
				gtx.Constraints.Min.X = gtx.Dp(minWidth) - gtx.Dp(unit.Dp(24))

				label := material.Body1(th, keyText)
				label.Color = h.Theme.ButtonText
				label.Font.Weight = font.Bold
				label.Font.Typeface = "monospace"
				label.Alignment = text.Middle
				return label.Layout(gtx)
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Badge background
			bgColor := h.Theme.SecondaryColor
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())

			// Badge border
			borderWidth := gtx.Dp(unit.Dp(1))
			drawBorder(gtx.Ops, rect, h.Theme.PrimaryColor, borderWidth)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
	)
}

// layoutCloseButton renders the close button.
func (h *HelpDialog) layoutCloseButton(gtx layout.Context, th *material.Theme) layout.Dimensions {
	buttonWidth := unit.Dp(120)
	buttonHeight := unit.Dp(40)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{X: gtx.Dp(buttonWidth), Y: gtx.Dp(buttonHeight)}
			gtx.Constraints.Max = gtx.Constraints.Min

			// Draw button background
			rect := image.Rectangle{Max: gtx.Constraints.Max}
			paint.FillShape(gtx.Ops, h.Theme.ButtonBackground, clip.Rect(rect).Op())

			// Draw button border
			borderWidth := gtx.Dp(unit.Dp(2))
			drawBorder(gtx.Ops, rect, h.Theme.PrimaryColor, borderWidth)

			// Layout button text centered
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(th, "閉じる")
				label.Color = h.Theme.ButtonText
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return h.closeButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			})
		}),
	)
}
