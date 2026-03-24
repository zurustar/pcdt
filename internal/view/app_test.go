// Package view contains the UI components for the countdown timer.
package view

import (
	"testing"

	"gioui.org/unit"

	"evangelion-timer/internal/controller"
	"evangelion-timer/internal/model"
)

func TestScreenState_String(t *testing.T) {
	tests := []struct {
		state    ScreenState
		expected string
	}{
		{ScreenInput, "Input"},
		{ScreenCountdown, "Countdown"},
		{ScreenState(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("ScreenState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewApp_DefaultValues(t *testing.T) {
	app := NewApp()

	// Check default screen is Input
	if app.CurrentScreen() != ScreenInput {
		t.Errorf("NewApp().CurrentScreen() = %v, want %v", app.CurrentScreen(), ScreenInput)
	}

	// Check theme is initialized
	if app.Theme == nil {
		t.Error("NewApp().Theme should not be nil")
	}

	// Check controller is initialized
	if app.Controller == nil {
		t.Error("NewApp().Controller should not be nil")
	}

	// Check config is initialized
	if app.Config == nil {
		t.Error("NewApp().Config should not be nil")
	}
}

func TestNewApp_WithOptions(t *testing.T) {
	customConfig := &model.AppConfig{
		AlwaysOnTop:  true,
		WindowWidth:  800,
		WindowHeight: 600,
		LastMinutes:  10,
		LastSeconds:  30,
	}
	customTheme := NewAlertTheme()
	customController := controller.NewTimerController()

	app := NewApp(
		WithConfig(customConfig),
		WithTheme(customTheme),
		WithController(customController),
	)

	if app.Config != customConfig {
		t.Error("WithConfig option not applied correctly")
	}

	if app.Theme != customTheme {
		t.Error("WithTheme option not applied correctly")
	}

	if app.Controller != customController {
		t.Error("WithController option not applied correctly")
	}

	// Check window dimensions from config
	width, height := app.GetWindowSize()
	if width != unit.Dp(800) || height != unit.Dp(600) {
		t.Errorf("Window size = (%v, %v), want (800, 600)", width, height)
	}
}

func TestApp_SetScreen(t *testing.T) {
	app := NewApp()

	// Initial screen should be Input
	if app.CurrentScreen() != ScreenInput {
		t.Errorf("Initial screen = %v, want %v", app.CurrentScreen(), ScreenInput)
	}

	// Set to Countdown
	app.SetScreen(ScreenCountdown)
	if app.CurrentScreen() != ScreenCountdown {
		t.Errorf("After SetScreen(ScreenCountdown), CurrentScreen() = %v, want %v", app.CurrentScreen(), ScreenCountdown)
	}

	// Set back to Input
	app.SetScreen(ScreenInput)
	if app.CurrentScreen() != ScreenInput {
		t.Errorf("After SetScreen(ScreenInput), CurrentScreen() = %v, want %v", app.CurrentScreen(), ScreenInput)
	}
}

func TestApp_WindowSize(t *testing.T) {
	app := NewApp()

	// Set new window size
	app.SetWindowSize(unit.Dp(500), unit.Dp(400))

	width, height := app.GetWindowSize()
	if width != unit.Dp(500) || height != unit.Dp(400) {
		t.Errorf("GetWindowSize() = (%v, %v), want (500, 400)", width, height)
	}

	// Check config is updated
	if app.Config.WindowWidth != 500 || app.Config.WindowHeight != 400 {
		t.Errorf("Config not updated: WindowWidth=%v, WindowHeight=%v", app.Config.WindowWidth, app.Config.WindowHeight)
	}
}

func TestApp_HandleResize(t *testing.T) {
	app := NewApp()

	app.HandleResize(600, 450)

	width, height := app.GetWindowSize()
	if width != unit.Dp(600) || height != unit.Dp(450) {
		t.Errorf("After HandleResize(600, 450), GetWindowSize() = (%v, %v), want (600, 450)", width, height)
	}
}

func TestApp_AlwaysOnTop(t *testing.T) {
	app := NewApp()

	// Default should be false
	if app.IsAlwaysOnTop() {
		t.Error("Default AlwaysOnTop should be false")
	}

	// Set to true
	app.SetAlwaysOnTop(true)
	if !app.IsAlwaysOnTop() {
		t.Error("After SetAlwaysOnTop(true), IsAlwaysOnTop() should be true")
	}

	// Toggle
	app.ToggleAlwaysOnTop()
	if app.IsAlwaysOnTop() {
		t.Error("After ToggleAlwaysOnTop(), IsAlwaysOnTop() should be false")
	}

	// Toggle again
	app.ToggleAlwaysOnTop()
	if !app.IsAlwaysOnTop() {
		t.Error("After second ToggleAlwaysOnTop(), IsAlwaysOnTop() should be true")
	}
}

func TestApp_TransitionToCountdown(t *testing.T) {
	app := NewApp()

	app.TransitionToCountdown()

	if app.CurrentScreen() != ScreenCountdown {
		t.Errorf("After TransitionToCountdown(), CurrentScreen() = %v, want %v", app.CurrentScreen(), ScreenCountdown)
	}
}

func TestApp_TransitionToInput(t *testing.T) {
	app := NewApp()

	// Start timer and transition to countdown
	app.Controller.Start(60)
	app.SetScreen(ScreenCountdown)

	// Transition back to input
	app.TransitionToInput()

	if app.CurrentScreen() != ScreenInput {
		t.Errorf("After TransitionToInput(), CurrentScreen() = %v, want %v", app.CurrentScreen(), ScreenInput)
	}

	// Timer should be reset (idle)
	if !app.Controller.GetModel().IsIdle() {
		t.Error("After TransitionToInput(), timer should be idle")
	}
}

func TestApp_InitialScreen(t *testing.T) {
	// This test validates the design document example:
	// "例: アプリ起動時に入力画面が表示される"
	app := NewApp()

	if app.CurrentScreen() != ScreenInput {
		t.Error("expected input screen on startup")
	}
}
