// Package view contains the UI components for the countdown timer.
package view

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"evangelion-timer/internal/controller"
)

// InputScreen represents the time input screen for the countdown timer.
type InputScreen struct {
	// Theme is the alert visual theme
	Theme *AlertTheme

	// Validator validates user input
	Validator *controller.InputValidator

	// Input editors for minutes and seconds
	minutesEditor widget.Editor
	secondsEditor widget.Editor

	// Start button
	startButton widget.Clickable

	// OnStart callback when start button is clicked with valid input
	OnStart func(minutes, seconds int)

	// Error message to display
	errorMessage string

	// Parsed values
	parsedMinutes int
	parsedSeconds int

	// Validation state
	isValid bool
}

// NewInputScreen creates a new input screen with the given theme.
func NewInputScreen(theme *AlertTheme) *InputScreen {
	s := &InputScreen{
		Theme:     theme,
		Validator: controller.NewInputValidator(),
	}

	// Configure minutes editor
	s.minutesEditor.SingleLine = true
	s.minutesEditor.Filter = "0123456789"
	s.minutesEditor.MaxLen = 2

	// Configure seconds editor
	s.secondsEditor.SingleLine = true
	s.secondsEditor.Filter = "0123456789"
	s.secondsEditor.MaxLen = 2

	return s
}

// SetOnStart sets the callback function called when start is triggered.
func (s *InputScreen) SetOnStart(callback func(minutes, seconds int)) {
	s.OnStart = callback
}

// filterDigits removes all non-digit characters from a string.
func filterDigits(input string) string {
	result := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			result = append(result, input[i])
		}
	}
	return string(result)
}

// sanitizeEditorInput removes non-digit characters from editor and updates if changed.
func (s *InputScreen) sanitizeEditorInput(editor *widget.Editor, maxLen int) {
	text := editor.Text()
	filtered := filterDigits(text)
	if len(filtered) > maxLen {
		filtered = filtered[:maxLen]
	}
	if filtered != text {
		editor.SetText(filtered)
	}
}

// validate validates the current input and updates the error message.
func (s *InputScreen) validate() {
	// Sanitize input to remove any non-digit characters (e.g., from Japanese IME)
	s.sanitizeEditorInput(&s.minutesEditor, 2)
	s.sanitizeEditorInput(&s.secondsEditor, 2)

	minutesText := s.minutesEditor.Text()
	secondsText := s.secondsEditor.Text()

	// Check if both fields are empty
	if minutesText == "" && secondsText == "" {
		s.errorMessage = ""
		s.isValid = false
		return
	}

	// Default empty fields to 0
	if minutesText == "" {
		minutesText = "0"
	}
	if secondsText == "" {
		secondsText = "0"
	}

	// Validate minutes
	minutes, err := s.Validator.ValidateMinutes(minutesText)
	if err != nil {
		s.errorMessage = err.Error()
		s.isValid = false
		return
	}

	// Validate seconds
	seconds, err := s.Validator.ValidateSeconds(secondsText)
	if err != nil {
		s.errorMessage = err.Error()
		s.isValid = false
		return
	}

	// Validate total time
	err = s.Validator.ValidateTotal(minutes, seconds)
	if err != nil {
		s.errorMessage = err.Error()
		s.isValid = false
		return
	}

	// Input is valid
	s.parsedMinutes = minutes
	s.parsedSeconds = seconds
	s.errorMessage = ""
	s.isValid = true
}

// handleStart handles the start action when input is valid.
func (s *InputScreen) handleStart() {
	if s.isValid && s.OnStart != nil {
		s.OnStart(s.parsedMinutes, s.parsedSeconds)
	}
}

// Layout renders the input screen.
func (s *InputScreen) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Handle keyboard events
	for {
		event, ok := gtx.Event(key.Filter{Name: key.NameReturn}, key.Filter{Name: key.NameEnter})
		if !ok {
			break
		}
		if e, ok := event.(key.Event); ok && e.State == key.Press {
			s.handleStart()
		}
	}

	// Check for start button click
	if s.startButton.Clicked(gtx) {
		s.handleStart()
	}

	// Validate input on every frame
	s.validate()

	// Fill background
	paint.Fill(gtx.Ops, s.Theme.BackgroundColor)

	// Center the content
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutTitle(gtx, th)
			}),
			// Spacer
			layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
			// Input fields
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutInputFields(gtx, th)
			}),
			// Error message
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutErrorMessage(gtx, th)
			}),
			// Spacer
			layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
			// Start button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutStartButton(gtx, th)
			}),
		)
	})
}

// layoutTitle renders the title text.
func (s *InputScreen) layoutTitle(gtx layout.Context, th *material.Theme) layout.Dimensions {
	title := material.H4(th, "COUNTDOWN TIMER")
	title.Color = s.Theme.PrimaryColor
	title.Font.Weight = font.Bold
	title.Alignment = text.Middle
	return title.Layout(gtx)
}

// layoutInputFields renders the minutes and seconds input fields.
func (s *InputScreen) layoutInputFields(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Minutes field
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutInputField(gtx, th, &s.minutesEditor, "分")
		}),
		// Colon separator
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutSeparator(gtx, th)
		}),
		// Seconds field
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutInputField(gtx, th, &s.secondsEditor, "秒")
		}),
	)
}

// layoutInputField renders a single input field with label.
func (s *InputScreen) layoutInputField(gtx layout.Context, th *material.Theme, editor *widget.Editor, label string) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Input box
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutInputBox(gtx, th, editor)
		}),
		// Label
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(th, label)
			lbl.Color = s.Theme.PrimaryColor
			lbl.Alignment = text.Middle
			return lbl.Layout(gtx)
		}),
	)
}

// layoutInputBox renders the styled input box.
func (s *InputScreen) layoutInputBox(gtx layout.Context, th *material.Theme, editor *widget.Editor) layout.Dimensions {
	// Define box dimensions
	boxWidth := unit.Dp(80)
	boxHeight := unit.Dp(60)

	// Create a macro to measure the editor
	macro := op.Record(gtx.Ops)

	// Layout the editor
	gtx.Constraints.Min = image.Point{X: gtx.Dp(boxWidth), Y: gtx.Dp(boxHeight)}
	gtx.Constraints.Max = gtx.Constraints.Min

	// Draw background
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, s.Theme.InputBackground, clip.Rect(rect).Op())

	// Draw border
	borderWidth := gtx.Dp(unit.Dp(2))
	drawBorder(gtx.Ops, rect, s.Theme.InputBorder, borderWidth)

	// Layout editor centered in the box
	dims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		ed := material.Editor(th, editor, "00")
		ed.Color = s.Theme.InputText
		ed.HintColor = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
		ed.TextSize = s.Theme.FontSizeLarge
		ed.Font.Weight = font.Bold
		return ed.Layout(gtx)
	})

	call := macro.Stop()
	call.Add(gtx.Ops)

	return dims
}

// drawBorder draws a border around a rectangle.
func drawBorder(ops *op.Ops, rect image.Rectangle, c color.NRGBA, width int) {
	// Top border
	paint.FillShape(ops, c, clip.Rect{
		Min: rect.Min,
		Max: image.Point{X: rect.Max.X, Y: rect.Min.Y + width},
	}.Op())
	// Bottom border
	paint.FillShape(ops, c, clip.Rect{
		Min: image.Point{X: rect.Min.X, Y: rect.Max.Y - width},
		Max: rect.Max,
	}.Op())
	// Left border
	paint.FillShape(ops, c, clip.Rect{
		Min: rect.Min,
		Max: image.Point{X: rect.Min.X + width, Y: rect.Max.Y},
	}.Op())
	// Right border
	paint.FillShape(ops, c, clip.Rect{
		Min: image.Point{X: rect.Max.X - width, Y: rect.Min.Y},
		Max: rect.Max,
	}.Op())
}

// layoutSeparator renders the colon separator between input fields.
func (s *InputScreen) layoutSeparator(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Left:  unit.Dp(16),
		Right: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		sep := material.H3(th, ":")
		sep.Color = s.Theme.PrimaryColor
		sep.Font.Weight = font.Bold
		return sep.Layout(gtx)
	})
}

// layoutErrorMessage renders the error message if present.
func (s *InputScreen) layoutErrorMessage(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if s.errorMessage == "" {
		// Reserve space for error message to prevent layout shift
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	}

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		errLabel := material.Body2(th, s.errorMessage)
		errLabel.Color = s.Theme.InputError
		errLabel.Alignment = text.Middle
		return errLabel.Layout(gtx)
	})
}

// layoutStartButton renders the start button.
func (s *InputScreen) layoutStartButton(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Determine button colors based on validity
	var bgColor, textColor color.NRGBA
	if s.isValid {
		bgColor = s.Theme.ButtonBackground
		textColor = s.Theme.ButtonText
	} else {
		bgColor = s.Theme.ButtonDisabled
		textColor = s.Theme.ButtonDisabledText
	}

	// Button dimensions
	buttonWidth := unit.Dp(160)
	buttonHeight := unit.Dp(50)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{X: gtx.Dp(buttonWidth), Y: gtx.Dp(buttonHeight)}
			gtx.Constraints.Max = gtx.Constraints.Min

			// Draw button background
			rect := image.Rectangle{Max: gtx.Constraints.Max}
			paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())

			// Draw button border
			borderWidth := gtx.Dp(unit.Dp(2))
			var borderColor color.NRGBA
			if s.isValid {
				borderColor = s.Theme.PrimaryColor
			} else {
				borderColor = s.Theme.ButtonDisabled
			}
			drawBorder(gtx.Ops, rect, borderColor, borderWidth)

			// Layout button text centered
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Body1(th, "スタート")
				btn.Color = textColor
				btn.Font.Weight = font.Bold
				return btn.Layout(gtx)
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Only make clickable if valid
			if s.isValid {
				return s.startButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Min}
				})
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
	)
}

// IsValid returns whether the current input is valid.
func (s *InputScreen) IsValid() bool {
	return s.isValid
}

// GetParsedMinutes returns the parsed minutes value.
func (s *InputScreen) GetParsedMinutes() int {
	return s.parsedMinutes
}

// GetParsedSeconds returns the parsed seconds value.
func (s *InputScreen) GetParsedSeconds() int {
	return s.parsedSeconds
}

// GetErrorMessage returns the current error message.
func (s *InputScreen) GetErrorMessage() string {
	return s.errorMessage
}

// SetMinutes sets the minutes input value.
func (s *InputScreen) SetMinutes(value string) {
	s.minutesEditor.SetText(value)
}

// SetSeconds sets the seconds input value.
func (s *InputScreen) SetSeconds(value string) {
	s.secondsEditor.SetText(value)
}

// GetMinutesText returns the current minutes input text.
func (s *InputScreen) GetMinutesText() string {
	return s.minutesEditor.Text()
}

// GetSecondsText returns the current seconds input text.
func (s *InputScreen) GetSecondsText() string {
	return s.secondsEditor.Text()
}

// Reset clears the input fields and error message.
func (s *InputScreen) Reset() {
	s.minutesEditor.SetText("")
	s.secondsEditor.SetText("")
	s.errorMessage = ""
	s.isValid = false
	s.parsedMinutes = 0
	s.parsedSeconds = 0
}
