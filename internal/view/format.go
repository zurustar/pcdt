// Package view contains the UI components for the countdown timer.
package view

import "fmt"

// FormatTime formats a time value in milliseconds to MM:SS:mm format.
// For negative values, a minus sign is prepended.
// The mm part represents the upper 2 digits of milliseconds (00-99).
//
// Examples:
//   - 330420 ms -> "05:30:42"
//   - 0 ms -> "00:00:00"
//   - -15730 ms -> "-00:15:73"
func FormatTime(milliseconds int64) string {
	if milliseconds < 0 {
		return "-" + formatAbsoluteTime(-milliseconds)
	}
	return formatAbsoluteTime(milliseconds)
}

// formatAbsoluteTime formats a non-negative time value in milliseconds to MM:SS:mm format.
func formatAbsoluteTime(milliseconds int64) string {
	// Extract milliseconds part (upper 2 digits: 00-99)
	ms := (milliseconds % 1000) / 10

	// Convert to total seconds
	totalSeconds := milliseconds / 1000

	// Extract minutes and seconds
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", minutes, seconds, ms)
}

// FormatTimeFromSeconds formats a time value in seconds to MM:SS:mm format.
// The milliseconds part will be 00 since we only have second precision.
// For negative values, a minus sign is prepended.
//
// Examples:
//   - 330 seconds -> "05:30:00"
//   - 0 seconds -> "00:00:00"
//   - -15 seconds -> "-00:15:00"
func FormatTimeFromSeconds(seconds int) string {
	return FormatTime(int64(seconds) * 1000)
}

// FormatMilliseconds formats just the milliseconds part (00-99).
// This is useful for the animated milliseconds display.
func FormatMilliseconds(milliseconds int64) string {
	if milliseconds < 0 {
		milliseconds = -milliseconds
	}
	ms := (milliseconds % 1000) / 10
	return fmt.Sprintf("%02d", ms)
}
