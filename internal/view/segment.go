// Package view contains the UI components for the countdown timer.
package view

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SevenSegmentDisplay renders numbers in 7-segment LED style.
type SevenSegmentDisplay struct {
	// Theme is the alert visual theme
	Theme *AlertTheme

	// SegmentWidth is the width of each segment
	SegmentWidth float32
	// SegmentLength is the length of each segment
	SegmentLength float32
	// DigitWidth is the total width of a single digit
	DigitWidth float32
	// DigitHeight is the total height of a single digit
	DigitHeight float32
	// DigitSpacing is the space between digits
	DigitSpacing float32
	// ColonWidth is the width of the colon separator
	ColonWidth float32
}

// NewSevenSegmentDisplay creates a new 7-segment display with default sizes.
func NewSevenSegmentDisplay(theme *AlertTheme) *SevenSegmentDisplay {
	return &SevenSegmentDisplay{
		Theme:         theme,
		SegmentWidth:  8,
		SegmentLength: 40,
		DigitWidth:    56,
		DigitHeight:   100,
		DigitSpacing:  12,
		ColonWidth:    24,
	}
}

// NewSevenSegmentDisplayWithSize creates a new 7-segment display with custom scale.
func NewSevenSegmentDisplayWithSize(theme *AlertTheme, scale float32) *SevenSegmentDisplay {
	return &SevenSegmentDisplay{
		Theme:         theme,
		SegmentWidth:  8 * scale,
		SegmentLength: 40 * scale,
		DigitWidth:    56 * scale,
		DigitHeight:   100 * scale,
		DigitSpacing:  12 * scale,
		ColonWidth:    24 * scale,
	}
}

// segmentMap defines which segments are on for each digit (0-9) and minus sign
// Segments are numbered:
//
//	 _0_
//	|   |
//	1   2
//	|_3_|
//	|   |
//	4   5
//	|_6_|
var segmentMap = map[rune][7]bool{
	'0': {true, true, true, false, true, true, true},
	'1': {false, false, true, false, false, true, false},
	'2': {true, false, true, true, true, false, true},
	'3': {true, false, true, true, false, true, true},
	'4': {false, true, true, true, false, true, false},
	'5': {true, true, false, true, false, true, true},
	'6': {true, true, false, true, true, true, true},
	'7': {true, false, true, false, false, true, false},
	'8': {true, true, true, true, true, true, true},
	'9': {true, true, true, true, false, true, true},
	'-': {false, false, false, true, false, false, false},
}

// LayoutTime renders a time string (e.g., "05:30:42" or "-00:15:73") in 7-segment style.
func (s *SevenSegmentDisplay) LayoutTime(gtx layout.Context, timeStr string, textColor color.NRGBA, dimColor color.NRGBA) layout.Dimensions {
	// Calculate total width
	var totalWidth float32
	for _, ch := range timeStr {
		if ch == ':' {
			totalWidth += s.ColonWidth
		} else {
			totalWidth += s.DigitWidth + s.DigitSpacing
		}
	}
	totalWidth -= s.DigitSpacing // Remove last spacing

	// Set constraints
	width := int(totalWidth)
	height := int(s.DigitHeight)

	gtx.Constraints.Min = image.Point{X: width, Y: height}
	gtx.Constraints.Max = gtx.Constraints.Min

	// Draw each character
	var offsetX float32
	for _, ch := range timeStr {
		if ch == ':' {
			s.drawColon(gtx.Ops, offsetX, textColor)
			offsetX += s.ColonWidth
		} else {
			s.drawDigit(gtx.Ops, offsetX, ch, textColor, dimColor)
			offsetX += s.DigitWidth + s.DigitSpacing
		}
	}

	return layout.Dimensions{Size: image.Point{X: width, Y: height}}
}

// drawDigit draws a single digit at the given x offset.
func (s *SevenSegmentDisplay) drawDigit(ops *op.Ops, offsetX float32, digit rune, onColor, offColor color.NRGBA) {
	segments, ok := segmentMap[digit]
	if !ok {
		return
	}

	// Draw all 7 segments
	// Segment 0: top horizontal
	s.drawHorizontalSegment(ops, offsetX+s.SegmentWidth, 0, segments[0], onColor, offColor)

	// Segment 1: upper left vertical
	s.drawVerticalSegment(ops, offsetX, s.SegmentWidth, segments[1], onColor, offColor)

	// Segment 2: upper right vertical
	s.drawVerticalSegment(ops, offsetX+s.DigitWidth-s.SegmentWidth, s.SegmentWidth, segments[2], onColor, offColor)

	// Segment 3: middle horizontal
	s.drawHorizontalSegment(ops, offsetX+s.SegmentWidth, s.DigitHeight/2-s.SegmentWidth/2, segments[3], onColor, offColor)

	// Segment 4: lower left vertical
	s.drawVerticalSegment(ops, offsetX, s.DigitHeight/2+s.SegmentWidth/2, segments[4], onColor, offColor)

	// Segment 5: lower right vertical
	s.drawVerticalSegment(ops, offsetX+s.DigitWidth-s.SegmentWidth, s.DigitHeight/2+s.SegmentWidth/2, segments[5], onColor, offColor)

	// Segment 6: bottom horizontal
	s.drawHorizontalSegment(ops, offsetX+s.SegmentWidth, s.DigitHeight-s.SegmentWidth, segments[6], onColor, offColor)
}

// drawHorizontalSegment draws a horizontal segment.
func (s *SevenSegmentDisplay) drawHorizontalSegment(ops *op.Ops, x, y float32, on bool, onColor, offColor color.NRGBA) {
	c := offColor
	if on {
		c = onColor
	}

	// Draw a hexagonal shape for the segment
	width := s.SegmentLength
	height := s.SegmentWidth
	inset := height / 2

	path := clip.Path{}
	path.Begin(ops)
	path.MoveTo(f32.Pt(x+inset, y))
	path.LineTo(f32.Pt(x+width-inset, y))
	path.LineTo(f32.Pt(x+width, y+height/2))
	path.LineTo(f32.Pt(x+width-inset, y+height))
	path.LineTo(f32.Pt(x+inset, y+height))
	path.LineTo(f32.Pt(x, y+height/2))
	path.Close()

	paint.FillShape(ops, c, clip.Outline{Path: path.End()}.Op())
}

// drawVerticalSegment draws a vertical segment.
func (s *SevenSegmentDisplay) drawVerticalSegment(ops *op.Ops, x, y float32, on bool, onColor, offColor color.NRGBA) {
	c := offColor
	if on {
		c = onColor
	}

	// Draw a hexagonal shape for the segment
	width := s.SegmentWidth
	height := s.SegmentLength
	inset := width / 2

	path := clip.Path{}
	path.Begin(ops)
	path.MoveTo(f32.Pt(x+width/2, y))
	path.LineTo(f32.Pt(x+width, y+inset))
	path.LineTo(f32.Pt(x+width, y+height-inset))
	path.LineTo(f32.Pt(x+width/2, y+height))
	path.LineTo(f32.Pt(x, y+height-inset))
	path.LineTo(f32.Pt(x, y+inset))
	path.Close()

	paint.FillShape(ops, c, clip.Outline{Path: path.End()}.Op())
}

// drawColon draws a colon separator at the given x offset.
func (s *SevenSegmentDisplay) drawColon(ops *op.Ops, offsetX float32, c color.NRGBA) {
	dotSize := s.SegmentWidth
	centerX := offsetX + s.ColonWidth/2 - dotSize/2

	// Upper dot
	upperY := s.DigitHeight/3 - dotSize/2
	rect1 := image.Rectangle{
		Min: image.Point{X: int(centerX), Y: int(upperY)},
		Max: image.Point{X: int(centerX + dotSize), Y: int(upperY + dotSize)},
	}
	paint.FillShape(ops, c, clip.Rect(rect1).Op())

	// Lower dot
	lowerY := s.DigitHeight*2/3 - dotSize/2
	rect2 := image.Rectangle{
		Min: image.Point{X: int(centerX), Y: int(lowerY)},
		Max: image.Point{X: int(centerX + dotSize), Y: int(lowerY + dotSize)},
	}
	paint.FillShape(ops, c, clip.Rect(rect2).Op())
}

// GetTotalWidth calculates the total width for a time string.
func (s *SevenSegmentDisplay) GetTotalWidth(timeStr string) unit.Dp {
	var totalWidth float32
	for _, ch := range timeStr {
		if ch == ':' {
			totalWidth += s.ColonWidth
		} else {
			totalWidth += s.DigitWidth + s.DigitSpacing
		}
	}
	totalWidth -= s.DigitSpacing
	return unit.Dp(totalWidth)
}

// GetHeight returns the height of the display.
func (s *SevenSegmentDisplay) GetHeight() unit.Dp {
	return unit.Dp(s.DigitHeight)
}
