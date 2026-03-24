// Package view contains the UI components for the countdown timer.
package view

import (
	"image"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"countdown-timer/internal/controller"
)

// CountdownScreen represents the countdown display screen.
type CountdownScreen struct {
	// Theme is the alert visual theme
	Theme *AlertTheme

	// Controller is the timer controller
	Controller *controller.TimerController

	// Buttons
	toggleButton widget.Clickable
	resetButton  widget.Clickable

	// OnReset callback when reset button is clicked
	OnReset func()

	// Millisecond animation state
	lastFrameTime time.Time
	milliseconds  int64 // Current milliseconds for animation (0-999)

	// Animation frame tracking for ~60fps
	frameInterval time.Duration
}

// NewCountdownScreen creates a new countdown screen with the given theme and controller.
func NewCountdownScreen(theme *AlertTheme, ctrl *controller.TimerController) *CountdownScreen {
	return &CountdownScreen{
		Theme:         theme,
		Controller:    ctrl,
		lastFrameTime: time.Now(),
		frameInterval: 16 * time.Millisecond, // ~60fps
	}
}

// SetOnReset sets the callback function called when reset is triggered.
func (s *CountdownScreen) SetOnReset(callback func()) {
	s.OnReset = callback
}

// handleToggle handles the toggle (stop/resume) action.
func (s *CountdownScreen) handleToggle() {
	s.Controller.Toggle()
}

// handleReset handles the reset action.
func (s *CountdownScreen) handleReset() {
	if s.OnReset != nil {
		s.OnReset()
	}
}

// updateMilliseconds updates the millisecond animation.
// Returns true if a redraw is needed.
func (s *CountdownScreen) updateMilliseconds() bool {
	now := time.Now()
	elapsed := now.Sub(s.lastFrameTime)

	if elapsed >= s.frameInterval {
		s.lastFrameTime = now

		model := s.Controller.GetModel()
		if model.IsRunning() {
			// Update milliseconds based on elapsed time
			s.milliseconds = (s.milliseconds + int64(elapsed.Milliseconds())) % 1000
		}
		return true
	}
	return false
}

// Layout renders the countdown screen.
func (s *CountdownScreen) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Handle keyboard events
	for {
		event, ok := gtx.Event(
			key.Filter{Name: key.NameSpace},
			key.Filter{Name: key.NameEscape},
		)
		if !ok {
			break
		}
		if e, ok := event.(key.Event); ok && e.State == key.Press {
			switch e.Name {
			case key.NameSpace:
				s.handleToggle()
			case key.NameEscape:
				s.handleReset()
			}
		}
	}

	// Check for button clicks
	if s.toggleButton.Clicked(gtx) {
		s.handleToggle()
	}
	if s.resetButton.Clicked(gtx) {
		s.handleReset()
	}

	// Update millisecond animation
	s.updateMilliseconds()

	// Request continuous animation for millisecond display
	model := s.Controller.GetModel()
	if model.IsRunning() {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(s.frameInterval)})
	}

	// Fill background
	paint.Fill(gtx.Ops, s.Theme.BackgroundColor)

	// Check if time is negative
	remainingSeconds := model.GetRemainingSeconds()
	isNegative := remainingSeconds < 0

	// Calculate scale based on window size
	windowWidth := float32(gtx.Constraints.Max.X)
	windowHeight := float32(gtx.Constraints.Max.Y)
	
	// Base scale on the smaller dimension to ensure it fits
	// Reference size: 800x500
	scaleX := windowWidth / 800.0
	scaleY := windowHeight / 500.0
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	// Clamp scale to reasonable range
	if scale < 0.5 {
		scale = 0.5
	}
	if scale > 2.0 {
		scale = 2.0
	}

	// Center the content
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Timer display
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutTimerDisplay(gtx, th, isNegative, scale)
			}),
			// Spacer
			layout.Rigid(layout.Spacer{Height: unit.Dp(40 * scale)}.Layout),
			// Control buttons
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutControlButtons(gtx, th, isNegative, scale)
			}),
		)
	})
}

// layoutTimerDisplay renders the main timer display with segment-style font.
func (s *CountdownScreen) layoutTimerDisplay(gtx layout.Context, th *material.Theme, isNegative bool, scale float32) layout.Dimensions {
	model := s.Controller.GetModel()
	remainingSeconds := model.GetRemainingSeconds()

	// Calculate total milliseconds for display
	var totalMs int64
	if model.IsRunning() {
		totalMs = int64(remainingSeconds)*1000 + s.milliseconds
		if remainingSeconds < 0 {
			totalMs = int64(remainingSeconds)*1000 - s.milliseconds
		}
	} else {
		totalMs = int64(remainingSeconds) * 1000
	}

	// Format the time
	timeStr := FormatTime(totalMs)

	// Get the appropriate color based on whether time is negative
	textColor := s.Theme.GetTextColorForState(isNegative)

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Start, // Align to left for the label
	}.Layout(gtx,
		// Label "活動限界まであと"
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutLabel(gtx, th, textColor, scale)
		}),
		// Spacer
		layout.Rigid(layout.Spacer{Height: unit.Dp(8 * scale)}.Layout),
		// Main time display
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutSegmentDisplay(gtx, th, timeStr, textColor, isNegative, scale)
		}),
		// Status indicator
		layout.Rigid(layout.Spacer{Height: unit.Dp(16 * scale)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutStatusIndicator(gtx, th, isNegative, scale)
		}),
	)
}

// layoutLabel renders the "活動限界まであと" label above the timer.
func (s *CountdownScreen) layoutLabel(gtx layout.Context, th *material.Theme, textColor color.NRGBA, scale float32) layout.Dimensions {
	label := material.Label(th, unit.Sp(18*scale), "活動限界まであと")
	label.Color = textColor
	// Use serif font for Japanese aesthetic
	label.Font.Typeface = "serif"
	return label.Layout(gtx)
}

// layoutSegmentDisplay renders the segment-style timer display using 7-segment LED style.
// MM:SS is displayed at full size, :mm (milliseconds) is displayed at 70% size.
func (s *CountdownScreen) layoutSegmentDisplay(gtx layout.Context, th *material.Theme, timeStr string, textColor color.NRGBA, isNegative bool, scale float32) layout.Dimensions {
	// Create segment displays with scaled size
	mainDisplay := NewSevenSegmentDisplayWithSize(s.Theme, 1.5*scale)
	millisDisplay := NewSevenSegmentDisplayWithSize(s.Theme, 1.05*scale)

	// Dim color for unlit segments (creates LED glow effect)
	dimColor := DimColor(textColor, 0.15)

	// Split time string into main part (MM:SS) and milliseconds part (:mm)
	// Format is "MM:SS:mm" or "-MM:SS:mm"
	var mainPart, millisPart string
	var hasMinusSign bool

	if len(timeStr) > 0 && timeStr[0] == '-' {
		hasMinusSign = true
		timeStr = timeStr[1:] // Remove minus sign for parsing
	}

	// Find the last colon to split MM:SS from mm
	lastColonIdx := -1
	for i := len(timeStr) - 1; i >= 0; i-- {
		if timeStr[i] == ':' {
			lastColonIdx = i
			break
		}
	}

	if lastColonIdx > 0 {
		mainPart = timeStr[:lastColonIdx]
		millisPart = timeStr[lastColonIdx:] // includes the colon
	} else {
		mainPart = timeStr
		millisPart = ""
	}

	if hasMinusSign {
		mainPart = "-" + mainPart
	}

	// Create a container with padding for the display
	return layout.Inset{
		Top:    unit.Dp(20 * scale),
		Bottom: unit.Dp(20 * scale),
		Left:   unit.Dp(30 * scale),
		Right:  unit.Dp(30 * scale),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Calculate display size for background
		mainWidth := mainDisplay.GetTotalWidth(mainPart)
		millisWidth := millisDisplay.GetTotalWidth(millisPart)
		displayWidth := mainWidth + millisWidth
		displayHeight := mainDisplay.GetHeight()

		// Add padding for background panel
		paddingH := unit.Dp(30 * scale)
		paddingV := unit.Dp(25 * scale)
		totalWidth := displayWidth + paddingH*2
		totalHeight := displayHeight + paddingV*2

		// Set constraints for the background
		gtx.Constraints.Min = image.Point{
			X: gtx.Dp(totalWidth),
			Y: gtx.Dp(totalHeight),
		}
		gtx.Constraints.Max = gtx.Constraints.Min

		// Draw background panel for segment display
		macro := op.Record(gtx.Ops)

		dims := layout.Stack{}.Layout(gtx,
			// Background panel
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return s.layoutDisplayBackground(gtx, isNegative, scale)
			}),
			// 7-segment LED display
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    paddingV,
					Bottom: paddingV,
					Left:   paddingH,
					Right:  paddingH,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// Layout main part and milliseconds part side by side
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.End, // Align to bottom
					}.Layout(gtx,
						// Main time (MM:SS)
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return mainDisplay.LayoutTime(gtx, mainPart, textColor, dimColor)
						}),
						// Milliseconds (:mm) - smaller
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return millisDisplay.LayoutTime(gtx, millisPart, textColor, dimColor)
						}),
					)
				})
			}),
		)

		call := macro.Stop()
		call.Add(gtx.Ops)

		return dims
	})
}

// layoutDisplayBackground renders the background panel for the timer display.
func (s *CountdownScreen) layoutDisplayBackground(gtx layout.Context, isNegative bool, scale float32) layout.Dimensions {
	// Background color varies by state
	var bgColor color.NRGBA
	if isNegative {
		// Darker red background for negative time
		bgColor = color.NRGBA{R: 40, G: 0, B: 0, A: 255}
	} else {
		// Dark background for normal state
		bgColor = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	}

	rect := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())

	// Draw border
	borderColor := s.Theme.PrimaryColor
	if isNegative {
		borderColor = s.Theme.TextColorNegative
	}
	borderWidth := gtx.Dp(unit.Dp(3 * scale))
	drawBorder(gtx.Ops, rect, borderColor, borderWidth)

	return layout.Dimensions{Size: gtx.Constraints.Min}
}

// layoutStatusIndicator renders the status text below the timer.
func (s *CountdownScreen) layoutStatusIndicator(gtx layout.Context, th *material.Theme, isNegative bool, scale float32) layout.Dimensions {
	model := s.Controller.GetModel()

	var statusText string
	var statusColor color.NRGBA

	switch {
	case model.IsNegative():
		statusText = "TIME OVER"
		statusColor = s.Theme.TextColorNegative
	case model.IsPaused():
		statusText = "PAUSED"
		statusColor = s.Theme.WarningColor
	case model.IsRunning():
		statusText = ""
		statusColor = s.Theme.TextColorNormal
	default:
		statusText = "READY"
		statusColor = s.Theme.TextColorNormal
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Label(th, unit.Sp(14*scale), statusText)
		label.Color = statusColor
		label.Font.Weight = font.Bold
		return label.Layout(gtx)
	})
}

// layoutControlButtons renders the stop/resume and reset buttons.
func (s *CountdownScreen) layoutControlButtons(gtx layout.Context, th *material.Theme, isNegative bool, scale float32) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEvenly,
	}.Layout(gtx,
		// Toggle button (Stop/Resume)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutToggleButton(gtx, th, isNegative, scale)
		}),
		// Spacer
		layout.Rigid(layout.Spacer{Width: unit.Dp(30 * scale)}.Layout),
		// Reset button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutResetButton(gtx, th, isNegative, scale)
		}),
	)
}

// layoutToggleButton renders the stop/resume toggle button.
func (s *CountdownScreen) layoutToggleButton(gtx layout.Context, th *material.Theme, isNegative bool, scale float32) layout.Dimensions {
	model := s.Controller.GetModel()

	// Determine button text based on state
	var buttonText string
	if model.IsRunning() {
		buttonText = "停止"
	} else {
		buttonText = "再開"
	}

	return s.layoutButton(gtx, th, &s.toggleButton, buttonText, isNegative, true, scale)
}

// layoutResetButton renders the reset button.
func (s *CountdownScreen) layoutResetButton(gtx layout.Context, th *material.Theme, isNegative bool, scale float32) layout.Dimensions {
	return s.layoutButton(gtx, th, &s.resetButton, "リセット", isNegative, true, scale)
}

// layoutButton renders a styled button.
func (s *CountdownScreen) layoutButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, text string, isNegative bool, enabled bool, scale float32) layout.Dimensions {
	// Button dimensions
	buttonWidth := unit.Dp(120 * scale)
	buttonHeight := unit.Dp(50 * scale)

	// Get colors based on negative state
	var bgColor, textColor, borderColor color.NRGBA
	if enabled {
		if isNegative {
			bgColor = DimColor(s.Theme.TextColorNegative, 0.3)
			textColor = s.Theme.TextColorNegative
			borderColor = s.Theme.TextColorNegative
		} else {
			bgColor = s.Theme.ButtonBackground
			textColor = s.Theme.ButtonText
			borderColor = s.Theme.PrimaryColor
		}
	} else {
		bgColor = s.Theme.ButtonDisabled
		textColor = s.Theme.ButtonDisabledText
		borderColor = s.Theme.ButtonDisabled
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{X: gtx.Dp(buttonWidth), Y: gtx.Dp(buttonHeight)}
			gtx.Constraints.Max = gtx.Constraints.Min

			// Draw button background
			rect := image.Rectangle{Max: gtx.Constraints.Max}
			paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())

			// Draw button border
			borderWidth := gtx.Dp(unit.Dp(2 * scale))
			drawBorder(gtx.Ops, rect, borderColor, borderWidth)

			// Layout button text centered
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(14*scale), text)
				label.Color = textColor
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			if enabled {
				return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Min}
				})
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
	)
}

// Reset resets the countdown screen state.
func (s *CountdownScreen) Reset() {
	s.milliseconds = 0
	s.lastFrameTime = time.Now()
}
