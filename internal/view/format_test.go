// Package view contains the UI components for the countdown timer.
package view

import (
	"fmt"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name         string
		milliseconds int64
		expected     string
	}{
		// Positive values
		{"5 minutes 30 seconds 420ms", 330420, "05:30:42"},
		{"1 minute", 60000, "01:00:00"},
		{"59 seconds 990ms", 59990, "00:59:99"},
		{"10 seconds", 10000, "00:10:00"},
		{"1 second", 1000, "00:01:00"},
		{"500ms", 500, "00:00:50"},
		{"99ms", 99, "00:00:09"},
		{"10ms", 10, "00:00:01"},
		{"9ms", 9, "00:00:00"},

		// Zero
		{"zero", 0, "00:00:00"},

		// Negative values
		{"negative 15 seconds 730ms", -15730, "-00:15:73"},
		{"negative 1 minute", -60000, "-01:00:00"},
		{"negative 5 minutes 30 seconds", -330000, "-05:30:00"},
		{"negative 500ms", -500, "-00:00:50"},

		// Large values
		{"99 minutes 59 seconds 990ms", 5999990, "99:59:99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.milliseconds)
			if result != tt.expected {
				t.Errorf("FormatTime(%d) = %q, want %q", tt.milliseconds, result, tt.expected)
			}
		})
	}
}

func TestFormatTimeFromSeconds(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"5 minutes 30 seconds", 330, "05:30:00"},
		{"zero", 0, "00:00:00"},
		{"negative 15 seconds", -15, "-00:15:00"},
		{"99 minutes 59 seconds", 5999, "99:59:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeFromSeconds(tt.seconds)
			if result != tt.expected {
				t.Errorf("FormatTimeFromSeconds(%d) = %q, want %q", tt.seconds, result, tt.expected)
			}
		})
	}
}

func TestFormatMilliseconds(t *testing.T) {
	tests := []struct {
		name         string
		milliseconds int64
		expected     string
	}{
		{"420ms", 420, "42"},
		{"990ms", 990, "99"},
		{"500ms", 500, "50"},
		{"99ms", 99, "09"},
		{"10ms", 10, "01"},
		{"9ms", 9, "00"},
		{"0ms", 0, "00"},
		{"negative 730ms", -730, "73"},
		{"1420ms (wraps)", 1420, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMilliseconds(tt.milliseconds)
			if result != tt.expected {
				t.Errorf("FormatMilliseconds(%d) = %q, want %q", tt.milliseconds, result, tt.expected)
			}
		})
	}
}


// Feature: countdown-timer, Property 7: 負の値の表示フォーマット
// **Validates: Requirements 4.2**
//
// *任意の*負の秒数に対して、フォーマット結果は `-MM:SS:mm` 形式であり、
// `-` プレフィックスが付き、分と秒は絶対値から計算される。
func TestFormatTime_Property7_NegativeValueFormat(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Any negative milliseconds value formats with a `-` prefix
	properties.Property("negative milliseconds format with minus prefix", prop.ForAll(
		func(ms int64) bool {
			negativeMs := -ms // Ensure negative
			formatted := FormatTime(negativeMs)

			// Must start with `-`
			return strings.HasPrefix(formatted, "-")
		},
		gen.Int64Range(1, 5999990), // 1ms to 99:59:99 in milliseconds
	))

	// Property: Minutes and seconds are calculated from absolute value
	properties.Property("minutes and seconds calculated from absolute value", prop.ForAll(
		func(ms int64) bool {
			negativeMs := -ms // Ensure negative
			formatted := FormatTime(negativeMs)

			// Calculate expected values from absolute value
			absMs := ms
			expectedMsDigits := (absMs % 1000) / 10
			totalSeconds := absMs / 1000
			expectedMinutes := totalSeconds / 60
			expectedSeconds := totalSeconds % 60

			expected := fmt.Sprintf("-%02d:%02d:%02d", expectedMinutes, expectedSeconds, expectedMsDigits)

			return formatted == expected
		},
		gen.Int64Range(1, 5999990), // 1ms to 99:59:99 in milliseconds
	))

	// Property: Negative format is exactly positive format with `-` prefix
	properties.Property("negative format equals minus prefix plus positive format", prop.ForAll(
		func(ms int64) bool {
			positiveFormatted := FormatTime(ms)
			negativeFormatted := FormatTime(-ms)

			// Negative format should be "-" + positive format
			return negativeFormatted == "-"+positiveFormatted
		},
		gen.Int64Range(1, 5999990), // 1ms to 99:59:99 in milliseconds
	))

	// Property: Zero is not negative (edge case)
	properties.Property("zero does not have minus prefix", prop.ForAll(
		func(_ bool) bool {
			formatted := FormatTime(0)
			return !strings.HasPrefix(formatted, "-") && formatted == "00:00:00"
		},
		gen.Bool(), // dummy generator to run the test
	))

	properties.TestingRun(t)
}
