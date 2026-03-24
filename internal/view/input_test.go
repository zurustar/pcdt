// Package view contains tests for the input screen.
package view

import (
	"strconv"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: countdown-timer, Property 2: 無効入力時のスタートボタン無効化
// **Validates: Requirements 2.2**
//
// *任意の*無効な入力状態（空、範囲外、非数値）に対して、スタートボタンは無効化されている。
// 有効な入力に対しては、IsValid()がtrueを返す。
//
// Note: The InputScreen editor has a filter that only allows digits (0-9) and
// MaxLen of 2, so negative numbers and values > 99 cannot be entered through
// the UI. These tests verify the validation logic for inputs that can actually
// occur through the UI.
func TestInputScreen_Property2_StartButtonDisabledForInvalidInput(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Empty input results in IsValid() returning false
	properties.Property("empty input results in invalid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			// Both fields empty
			screen.SetMinutes("")
			screen.SetSeconds("")
			screen.validate()

			return !screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Only minutes empty with zero seconds results in invalid state
	properties.Property("only minutes empty with zero seconds results in invalid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			// Minutes empty, seconds is 0 (total would be 0)
			screen.SetMinutes("")
			screen.SetSeconds("0")
			screen.validate()

			return !screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Only seconds empty with zero minutes results in invalid state
	properties.Property("only seconds empty with zero minutes results in invalid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			// Seconds empty, minutes is 0 (total would be 0)
			screen.SetMinutes("0")
			screen.SetSeconds("")
			screen.validate()

			return !screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Zero total time (0 minutes, 0 seconds) results in invalid state
	properties.Property("zero total time results in invalid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			screen.SetMinutes("0")
			screen.SetSeconds("0")
			screen.validate()

			return !screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Valid input results in IsValid() returning true
	// Note: The editor has MaxLen=2, so we test with valid 2-digit values
	properties.Property("valid input results in valid state", prop.ForAll(
		func(minutes, seconds int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			total := minutes*60 + seconds
			// Skip if total is 0 (invalid case)
			if total < 1 {
				return true
			}

			screen.SetMinutes(strconv.Itoa(minutes))
			screen.SetSeconds(strconv.Itoa(seconds))
			screen.validate()

			return screen.IsValid()
		},
		gen.IntRange(0, 99), // valid minutes range (0-99, max 2 digits)
		gen.IntRange(0, 59), // valid seconds range (0-59)
	))

	// Property: Valid input with only minutes (seconds defaults to 0)
	properties.Property("valid minutes with empty seconds results in valid state", prop.ForAll(
		func(minutes int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			// Minutes > 0, seconds empty (defaults to 0, total > 0)
			screen.SetMinutes(strconv.Itoa(minutes))
			screen.SetSeconds("")
			screen.validate()

			return screen.IsValid()
		},
		gen.IntRange(1, 99), // valid minutes range (> 0 to ensure total > 0)
	))

	// Property: Valid input with only seconds (minutes defaults to 0)
	properties.Property("valid seconds with empty minutes results in valid state", prop.ForAll(
		func(seconds int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			// Seconds > 0, minutes empty (defaults to 0, total > 0)
			screen.SetMinutes("")
			screen.SetSeconds(strconv.Itoa(seconds))
			screen.validate()

			return screen.IsValid()
		},
		gen.IntRange(1, 59), // valid seconds range (> 0 to ensure total > 0)
	))

	// Property: Maximum valid time (99:59) results in valid state
	properties.Property("maximum valid time results in valid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			screen.SetMinutes("99")
			screen.SetSeconds("59")
			screen.validate()

			return screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Minimum valid time (0:01) results in valid state
	properties.Property("minimum valid time results in valid state", prop.ForAll(
		func(_ int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			screen.SetMinutes("0")
			screen.SetSeconds("1")
			screen.validate()

			return screen.IsValid()
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Seconds out of range (60-99, which can be typed as 2 digits) results in invalid state
	// Note: The editor allows 2 digits, so values 60-99 can be entered for seconds
	properties.Property("seconds out of range (60-99) results in invalid state", prop.ForAll(
		func(seconds int) bool {
			theme := NewAlertTheme()
			screen := NewInputScreen(theme)

			screen.SetMinutes("5")
			screen.SetSeconds(strconv.Itoa(seconds))
			screen.validate()

			return !screen.IsValid()
		},
		gen.IntRange(60, 99), // invalid seconds range (can be typed as 2 digits)
	))

	properties.TestingRun(t)
}
