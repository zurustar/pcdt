// Package view contains the UI components for the countdown timer.
package view

import (
	"testing"
)

func TestNewHelpDialog(t *testing.T) {
	theme := NewAlertTheme()
	dialog := NewHelpDialog(theme)

	if dialog == nil {
		t.Fatal("NewHelpDialog returned nil")
	}

	if dialog.Theme != theme {
		t.Error("Theme not set correctly")
	}

	if dialog.IsVisible() {
		t.Error("Dialog should be hidden by default")
	}
}

func TestHelpDialog_ShowHide(t *testing.T) {
	theme := NewAlertTheme()
	dialog := NewHelpDialog(theme)

	// Initially hidden
	if dialog.IsVisible() {
		t.Error("Dialog should be hidden initially")
	}

	// Show
	dialog.Show()
	if !dialog.IsVisible() {
		t.Error("Dialog should be visible after Show()")
	}

	// Hide
	dialog.Hide()
	if dialog.IsVisible() {
		t.Error("Dialog should be hidden after Hide()")
	}
}

func TestHelpDialog_Toggle(t *testing.T) {
	theme := NewAlertTheme()
	dialog := NewHelpDialog(theme)

	// Initially hidden
	if dialog.IsVisible() {
		t.Error("Dialog should be hidden initially")
	}

	// Toggle to show
	dialog.Toggle()
	if !dialog.IsVisible() {
		t.Error("Dialog should be visible after first Toggle()")
	}

	// Toggle to hide
	dialog.Toggle()
	if dialog.IsVisible() {
		t.Error("Dialog should be hidden after second Toggle()")
	}

	// Toggle to show again
	dialog.Toggle()
	if !dialog.IsVisible() {
		t.Error("Dialog should be visible after third Toggle()")
	}
}

func TestGetShortcuts(t *testing.T) {
	shortcuts := getShortcuts()

	if len(shortcuts) == 0 {
		t.Fatal("getShortcuts returned empty list")
	}

	// Verify expected shortcuts are present
	expectedKeys := map[string]bool{
		"Enter":  false,
		"Space":  false,
		"Escape": false,
		"H / ?":  false,
	}

	for _, shortcut := range shortcuts {
		if shortcut.key == "" {
			t.Error("Shortcut key should not be empty")
		}
		if shortcut.description == "" {
			t.Error("Shortcut description should not be empty")
		}

		if _, exists := expectedKeys[shortcut.key]; exists {
			expectedKeys[shortcut.key] = true
		}
	}

	// Check all expected keys are present
	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Expected shortcut key %q not found", key)
		}
	}
}

func TestHelpDialog_MultipleShowCalls(t *testing.T) {
	theme := NewAlertTheme()
	dialog := NewHelpDialog(theme)

	// Multiple Show calls should keep it visible
	dialog.Show()
	dialog.Show()
	dialog.Show()

	if !dialog.IsVisible() {
		t.Error("Dialog should remain visible after multiple Show() calls")
	}
}

func TestHelpDialog_MultipleHideCalls(t *testing.T) {
	theme := NewAlertTheme()
	dialog := NewHelpDialog(theme)

	dialog.Show()

	// Multiple Hide calls should keep it hidden
	dialog.Hide()
	dialog.Hide()
	dialog.Hide()

	if dialog.IsVisible() {
		t.Error("Dialog should remain hidden after multiple Hide() calls")
	}
}
