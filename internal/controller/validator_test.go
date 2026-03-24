// Package controller contains tests for the input validator.
package controller

import (
	"strconv"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: countdown-timer, Property 1: 入力バリデーションの正確性
// **Validates: Requirements 1.4, 1.5**
//
// *任意の*分（0-99）と秒（0-59）の組み合わせに対して、合計が1秒以上99分59秒以下であれば
// 有効と判定され、それ以外は無効と判定される。
func TestInputValidator_Property1_ValidationAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	validator := NewInputValidator()

	// Property: Any minutes (0-99) and seconds (0-59) combination with total >= 1 and <= 5999 is valid
	// This tests the complete validation flow: ValidateMinutes -> ValidateSeconds -> ValidateTotal
	properties.Property("valid input range is accepted", prop.ForAll(
		func(minutes, seconds int) bool {
			total := minutes*60 + seconds

			// Validate minutes string input
			minutesStr := strconv.Itoa(minutes)
			parsedMinutes, errMinutes := validator.ValidateMinutes(minutesStr)

			// Validate seconds string input
			secondsStr := strconv.Itoa(seconds)
			parsedSeconds, errSeconds := validator.ValidateSeconds(secondsStr)

			// Check if individual validations pass
			isValidMinutes := errMinutes == nil && parsedMinutes == minutes
			isValidSeconds := errSeconds == nil && parsedSeconds == seconds

			// If individual validations pass, check total
			if isValidMinutes && isValidSeconds {
				errTotal := validator.ValidateTotal(parsedMinutes, parsedSeconds)
				isValidTotal := total >= 1 && total <= 5999
				return (errTotal == nil) == isValidTotal
			}

			// If individual validations fail, that's expected for out-of-range values
			return true
		},
		gen.IntRange(0, 99), // minutes: 0-99 (valid range)
		gen.IntRange(0, 59), // seconds: 0-59 (valid range)
	))

	// Property: Minutes outside valid range (0-99) are rejected
	properties.Property("invalid minutes range is rejected", prop.ForAll(
		func(minutes int) bool {
			minutesStr := strconv.Itoa(minutes)
			_, err := validator.ValidateMinutes(minutesStr)

			isInvalidRange := minutes < 0 || minutes > 99
			if isInvalidRange {
				return err == ErrInvalidMinutes
			}
			return err == nil
		},
		gen.IntRange(-10, 110), // include invalid range
	))

	// Property: Seconds outside valid range (0-59) are rejected
	properties.Property("invalid seconds range is rejected", prop.ForAll(
		func(seconds int) bool {
			secondsStr := strconv.Itoa(seconds)
			_, err := validator.ValidateSeconds(secondsStr)

			isInvalidRange := seconds < 0 || seconds > 59
			if isInvalidRange {
				return err == ErrInvalidSeconds
			}
			return err == nil
		},
		gen.IntRange(-10, 70), // include invalid range
	))

	// Property: Zero total time (0 minutes, 0 seconds) is rejected
	properties.Property("zero total time is rejected", prop.ForAll(
		func(_ int) bool {
			errTotal := validator.ValidateTotal(0, 0)
			return errTotal == ErrZeroTime
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Maximum valid time (99:59 = 5999 seconds) is accepted
	properties.Property("maximum valid time is accepted", prop.ForAll(
		func(_ int) bool {
			// First validate individual inputs
			_, errMinutes := validator.ValidateMinutes("99")
			_, errSeconds := validator.ValidateSeconds("59")
			errTotal := validator.ValidateTotal(99, 59)

			return errMinutes == nil && errSeconds == nil && errTotal == nil
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Minimum valid time (0:01 = 1 second) is accepted
	properties.Property("minimum valid time is accepted", prop.ForAll(
		func(_ int) bool {
			// First validate individual inputs
			_, errMinutes := validator.ValidateMinutes("0")
			_, errSeconds := validator.ValidateSeconds("1")
			errTotal := validator.ValidateTotal(0, 1)

			return errMinutes == nil && errSeconds == nil && errTotal == nil
		},
		gen.IntRange(0, 1), // dummy generator
	))

	// Property: Total time exceeding 5999 seconds is rejected
	properties.Property("total time exceeding maximum is rejected", prop.ForAll(
		func(extraSeconds int) bool {
			// 99 minutes 59 seconds = 5999 seconds is max
			// Test with values that would exceed this
			// Using 100 minutes would be invalid minutes, so we test ValidateTotal directly
			// with values that are individually valid but total exceeds limit
			// Note: With valid minutes (0-99) and seconds (0-59), max is 99*60+59=5999
			// So we can't actually exceed 5999 with valid individual inputs
			// This property verifies that ValidateTotal correctly rejects totals > 5999
			total := 6000 + extraSeconds
			minutes := total / 60
			seconds := total % 60

			// If minutes would be > 99, ValidateTotal should still reject based on total
			errTotal := validator.ValidateTotal(minutes, seconds)
			return errTotal != nil
		},
		gen.IntRange(0, 100), // extra seconds beyond 5999
	))

	properties.TestingRun(t)
}

// ============================================================================
// Unit Tests for Edge Cases
// **Validates: Requirements 1.5, 1.6**
// ============================================================================

// TestValidateMinutes_EmptyInput tests that empty string input returns ErrEmptyInput
func TestValidateMinutes_EmptyInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"whitespace only - space", " "},
		{"whitespace only - tab", "\t"},
		{"whitespace only - multiple spaces", "   "},
		{"whitespace only - mixed", " \t "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateMinutes(tc.input)
			if err != ErrEmptyInput {
				t.Errorf("ValidateMinutes(%q) = %v, want %v", tc.input, err, ErrEmptyInput)
			}
		})
	}
}

// TestValidateSeconds_EmptyInput tests that empty string input returns ErrEmptyInput
func TestValidateSeconds_EmptyInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"whitespace only - space", " "},
		{"whitespace only - tab", "\t"},
		{"whitespace only - multiple spaces", "   "},
		{"whitespace only - mixed", " \t "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateSeconds(tc.input)
			if err != ErrEmptyInput {
				t.Errorf("ValidateSeconds(%q) = %v, want %v", tc.input, err, ErrEmptyInput)
			}
		})
	}
}

// TestValidateMinutes_NonNumericInput tests that non-numeric input returns ErrInvalidFormat
func TestValidateMinutes_NonNumericInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"letters only", "abc"},
		{"single letter", "a"},
		{"uppercase letters", "ABC"},
		{"mixed case", "AbC"},
		{"special characters", "!@#"},
		{"symbols", "$%^&*"},
		{"mixed letters and numbers", "12abc"},
		{"mixed numbers and letters", "abc12"},
		{"letters with spaces", "a b c"},
		{"unicode characters", "日本語"},
		{"emoji", "🕐"},
		{"equals sign", "=10"},
		{"slash", "5/2"},
		{"asterisk", "5*2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateMinutes(tc.input)
			if err != ErrInvalidFormat {
				t.Errorf("ValidateMinutes(%q) = %v, want %v", tc.input, err, ErrInvalidFormat)
			}
		})
	}
}

// TestValidateMinutes_PlusSign tests that plus sign prefix is accepted by strconv.Atoi
// Note: Go's strconv.Atoi accepts leading + sign, so "+5" parses as 5
func TestValidateMinutes_PlusSign(t *testing.T) {
	validator := NewInputValidator()

	result, err := validator.ValidateMinutes("+5")
	if err != nil {
		t.Errorf("ValidateMinutes(\"+5\") returned unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("ValidateMinutes(\"+5\") = %d, want 5", result)
	}
}

// TestValidateSeconds_NonNumericInput tests that non-numeric input returns ErrInvalidFormat
func TestValidateSeconds_NonNumericInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"letters only", "abc"},
		{"single letter", "a"},
		{"uppercase letters", "ABC"},
		{"mixed case", "AbC"},
		{"special characters", "!@#"},
		{"symbols", "$%^&*"},
		{"mixed letters and numbers", "12abc"},
		{"mixed numbers and letters", "abc12"},
		{"letters with spaces", "a b c"},
		{"unicode characters", "日本語"},
		{"emoji", "🕐"},
		{"equals sign", "=10"},
		{"slash", "5/2"},
		{"asterisk", "5*2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateSeconds(tc.input)
			if err != ErrInvalidFormat {
				t.Errorf("ValidateSeconds(%q) = %v, want %v", tc.input, err, ErrInvalidFormat)
			}
		})
	}
}

// TestValidateSeconds_PlusSign tests that plus sign prefix is accepted by strconv.Atoi
// Note: Go's strconv.Atoi accepts leading + sign, so "+5" parses as 5
func TestValidateSeconds_PlusSign(t *testing.T) {
	validator := NewInputValidator()

	result, err := validator.ValidateSeconds("+5")
	if err != nil {
		t.Errorf("ValidateSeconds(\"+5\") returned unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("ValidateSeconds(\"+5\") = %d, want 5", result)
	}
}

// TestValidateMinutes_DecimalNumbers tests that decimal numbers return ErrInvalidFormat
func TestValidateMinutes_DecimalNumbers(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"simple decimal", "1.5"},
		{"decimal with zero", "0.5"},
		{"large decimal", "99.9"},
		{"multiple decimal points", "1.2.3"},
		{"decimal at start", ".5"},
		{"decimal at end", "5."},
		{"comma as decimal", "1,5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateMinutes(tc.input)
			if err != ErrInvalidFormat {
				t.Errorf("ValidateMinutes(%q) = %v, want %v", tc.input, err, ErrInvalidFormat)
			}
		})
	}
}

// TestValidateSeconds_DecimalNumbers tests that decimal numbers return ErrInvalidFormat
func TestValidateSeconds_DecimalNumbers(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name  string
		input string
	}{
		{"simple decimal", "1.5"},
		{"decimal with zero", "0.5"},
		{"large decimal", "59.9"},
		{"multiple decimal points", "1.2.3"},
		{"decimal at start", ".5"},
		{"decimal at end", "5."},
		{"comma as decimal", "1,5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateSeconds(tc.input)
			if err != ErrInvalidFormat {
				t.Errorf("ValidateSeconds(%q) = %v, want %v", tc.input, err, ErrInvalidFormat)
			}
		})
	}
}

// TestValidateMinutes_NegativeNumbers tests that negative number strings return ErrInvalidMinutes
// Note: "-0" is parsed as 0 by strconv.Atoi, which is valid
func TestValidateMinutes_NegativeNumbers(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"negative one", "-1", ErrInvalidMinutes},
		{"negative ten", "-10", ErrInvalidMinutes},
		{"large negative", "-100", ErrInvalidMinutes},
		{"negative with space", "- 5", ErrInvalidFormat},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateMinutes(tc.input)
			if err != tc.expectedErr {
				t.Errorf("ValidateMinutes(%q) = %v, want %v", tc.input, err, tc.expectedErr)
			}
		})
	}
}

// TestValidateMinutes_NegativeZero tests that "-0" is parsed as 0 (valid)
// Note: Go's strconv.Atoi parses "-0" as 0
func TestValidateMinutes_NegativeZero(t *testing.T) {
	validator := NewInputValidator()

	result, err := validator.ValidateMinutes("-0")
	if err != nil {
		t.Errorf("ValidateMinutes(\"-0\") returned unexpected error: %v", err)
	}
	if result != 0 {
		t.Errorf("ValidateMinutes(\"-0\") = %d, want 0", result)
	}
}

// TestValidateSeconds_NegativeNumbers tests that negative number strings return ErrInvalidSeconds
// Note: "-0" is parsed as 0 by strconv.Atoi, which is valid
func TestValidateSeconds_NegativeNumbers(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"negative one", "-1", ErrInvalidSeconds},
		{"negative ten", "-10", ErrInvalidSeconds},
		{"large negative", "-100", ErrInvalidSeconds},
		{"negative with space", "- 5", ErrInvalidFormat},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validator.ValidateSeconds(tc.input)
			if err != tc.expectedErr {
				t.Errorf("ValidateSeconds(%q) = %v, want %v", tc.input, err, tc.expectedErr)
			}
		})
	}
}

// TestValidateSeconds_NegativeZero tests that "-0" is parsed as 0 (valid)
// Note: Go's strconv.Atoi parses "-0" as 0
func TestValidateSeconds_NegativeZero(t *testing.T) {
	validator := NewInputValidator()

	result, err := validator.ValidateSeconds("-0")
	if err != nil {
		t.Errorf("ValidateSeconds(\"-0\") returned unexpected error: %v", err)
	}
	if result != 0 {
		t.Errorf("ValidateSeconds(\"-0\") = %d, want 0", result)
	}
}

// TestValidateMinutes_WhitespaceHandling tests that whitespace is properly trimmed
func TestValidateMinutes_WhitespaceHandling(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"leading space", " 5", 5},
		{"trailing space", "5 ", 5},
		{"both spaces", " 5 ", 5},
		{"leading tab", "\t5", 5},
		{"trailing tab", "5\t", 5},
		{"multiple leading spaces", "   10", 10},
		{"multiple trailing spaces", "10   ", 10},
		{"mixed whitespace", " \t 15 \t ", 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validator.ValidateMinutes(tc.input)
			if err != nil {
				t.Errorf("ValidateMinutes(%q) returned error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("ValidateMinutes(%q) = %d, want %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestValidateSeconds_WhitespaceHandling tests that whitespace is properly trimmed
func TestValidateSeconds_WhitespaceHandling(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"leading space", " 5", 5},
		{"trailing space", "5 ", 5},
		{"both spaces", " 5 ", 5},
		{"leading tab", "\t5", 5},
		{"trailing tab", "5\t", 5},
		{"multiple leading spaces", "   10", 10},
		{"multiple trailing spaces", "10   ", 10},
		{"mixed whitespace", " \t 15 \t ", 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validator.ValidateSeconds(tc.input)
			if err != nil {
				t.Errorf("ValidateSeconds(%q) returned error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("ValidateSeconds(%q) = %d, want %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestValidateMinutes_RangeEdgeCases tests boundary values for minutes
func TestValidateMinutes_RangeEdgeCases(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name        string
		input       string
		expected    int
		expectedErr error
	}{
		{"minimum valid", "0", 0, nil},
		{"maximum valid", "99", 99, nil},
		{"just above maximum", "100", 0, ErrInvalidMinutes},
		{"well above maximum", "150", 0, ErrInvalidMinutes},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validator.ValidateMinutes(tc.input)
			if err != tc.expectedErr {
				t.Errorf("ValidateMinutes(%q) error = %v, want %v", tc.input, err, tc.expectedErr)
			}
			if tc.expectedErr == nil && result != tc.expected {
				t.Errorf("ValidateMinutes(%q) = %d, want %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestValidateSeconds_RangeEdgeCases tests boundary values for seconds
func TestValidateSeconds_RangeEdgeCases(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name        string
		input       string
		expected    int
		expectedErr error
	}{
		{"minimum valid", "0", 0, nil},
		{"maximum valid", "59", 59, nil},
		{"just above maximum", "60", 0, ErrInvalidSeconds},
		{"well above maximum", "100", 0, ErrInvalidSeconds},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validator.ValidateSeconds(tc.input)
			if err != tc.expectedErr {
				t.Errorf("ValidateSeconds(%q) error = %v, want %v", tc.input, err, tc.expectedErr)
			}
			if tc.expectedErr == nil && result != tc.expected {
				t.Errorf("ValidateSeconds(%q) = %d, want %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestValidateTotal_EdgeCases tests boundary values for total time validation
func TestValidateTotal_EdgeCases(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name        string
		minutes     int
		seconds     int
		expectedErr error
	}{
		{"zero total", 0, 0, ErrZeroTime},
		{"minimum valid (1 second)", 0, 1, nil},
		{"one minute", 1, 0, nil},
		{"maximum valid", 99, 59, nil},
		{"exceeds maximum", 100, 0, ErrInvalidMinutes},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateTotal(tc.minutes, tc.seconds)
			if err != tc.expectedErr {
				t.Errorf("ValidateTotal(%d, %d) = %v, want %v", tc.minutes, tc.seconds, err, tc.expectedErr)
			}
		})
	}
}
