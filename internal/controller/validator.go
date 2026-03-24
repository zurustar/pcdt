// Package controller contains the business logic controllers for the countdown timer.
package controller

import (
	"errors"
	"strconv"
	"strings"
)

// Error messages for input validation
var (
	ErrEmptyInput     = errors.New("時間を入力してください")
	ErrInvalidMinutes = errors.New("分は0〜99の範囲で入力してください")
	ErrInvalidSeconds = errors.New("秒は0〜59の範囲で入力してください")
	ErrZeroTime       = errors.New("1秒以上の時間を入力してください")
	ErrInvalidFormat  = errors.New("数値を入力してください")
)

// InputValidator validates time input for the countdown timer.
type InputValidator struct{}

// NewInputValidator creates a new InputValidator instance.
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// ValidateMinutes validates the minutes input string.
// Returns the parsed minutes value if valid, or an error if invalid.
// Valid range: 0-99
func (v *InputValidator) ValidateMinutes(input string) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, ErrEmptyInput
	}

	minutes, err := strconv.Atoi(input)
	if err != nil {
		return 0, ErrInvalidFormat
	}

	if minutes < 0 || minutes > 99 {
		return 0, ErrInvalidMinutes
	}

	return minutes, nil
}

// ValidateSeconds validates the seconds input string.
// Returns the parsed seconds value if valid, or an error if invalid.
// Valid range: 0-59
func (v *InputValidator) ValidateSeconds(input string) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, ErrEmptyInput
	}

	seconds, err := strconv.Atoi(input)
	if err != nil {
		return 0, ErrInvalidFormat
	}

	if seconds < 0 || seconds > 59 {
		return 0, ErrInvalidSeconds
	}

	return seconds, nil
}

// ValidateTotal validates that the total time is within the valid range.
// Total must be at least 1 second and at most 99 minutes 59 seconds (5999 seconds).
func (v *InputValidator) ValidateTotal(minutes, seconds int) error {
	total := minutes*60 + seconds

	if total < 1 {
		return ErrZeroTime
	}

	if total > 5999 {
		return ErrInvalidMinutes
	}

	return nil
}
