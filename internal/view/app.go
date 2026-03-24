// Package view contains the UI components for the countdown timer.
package view

import (
	"sync"

	"gioui.org/app"
	"gioui.org/unit"

	"evangelion-timer/internal/controller"
	"evangelion-timer/internal/model"
)

// ScreenState represents the current screen of the application
type ScreenState int

const (
	// ScreenInput is the initial time input screen
	ScreenInput ScreenState = iota
	// ScreenCountdown is the countdown display screen
	ScreenCountdown
)

// String returns the string representation of ScreenState
func (s ScreenState) String() string {
	switch s {
	case ScreenInput:
		return "Input"
	case ScreenCountdown:
		return "Countdown"
	default:
		return "Unknown"
	}
}

// App represents the main application window and state management
type App struct {
	// Window is the Gio application window
	Window *app.Window

	// Theme is the alert-style visual theme
	Theme *AlertTheme

	// Controller is the timer controller
	Controller *controller.TimerController

	// Config is the application configuration
	Config *model.AppConfig

	// currentScreen is the current screen state
	currentScreen ScreenState

	// mu protects concurrent access to app state
	mu sync.RWMutex

	// Window dimensions
	windowWidth  unit.Dp
	windowHeight unit.Dp
}

// AppOption is a functional option for configuring the App
type AppOption func(*App)

// WithConfig sets the application configuration
func WithConfig(config *model.AppConfig) AppOption {
	return func(a *App) {
		a.Config = config
	}
}

// WithController sets the timer controller
func WithController(ctrl *controller.TimerController) AppOption {
	return func(a *App) {
		a.Controller = ctrl
	}
}

// WithTheme sets the visual theme
func WithTheme(theme *AlertTheme) AppOption {
	return func(a *App) {
		a.Theme = theme
	}
}

// NewApp creates a new application instance with the given options
func NewApp(opts ...AppOption) *App {
	a := &App{
		currentScreen: ScreenInput,
		Theme:         NewAlertTheme(),
		Controller:    controller.NewTimerController(),
		Config:        model.DefaultConfig(),
	}

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	// Set window dimensions from config
	a.windowWidth = unit.Dp(a.Config.WindowWidth)
	a.windowHeight = unit.Dp(a.Config.WindowHeight)

	return a
}

// CreateWindow creates and configures the Gio application window
func (a *App) CreateWindow() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create window options
	options := []app.Option{
		app.Title("Countdown Timer"),
		app.Size(a.windowWidth, a.windowHeight),
		app.MinSize(unit.Dp(300), unit.Dp(200)),
	}

	// Create the window
	a.Window = new(app.Window)
	a.Window.Option(options...)

	// Apply always-on-top if configured
	if a.Config.AlwaysOnTop {
		a.SetAlwaysOnTop(true)
	}
}

// CurrentScreen returns the current screen state
func (a *App) CurrentScreen() ScreenState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.currentScreen
}

// SetScreen transitions to the specified screen state
func (a *App) SetScreen(screen ScreenState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.currentScreen = screen
}

// GetWindowSize returns the current window dimensions
func (a *App) GetWindowSize() (width, height unit.Dp) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.windowWidth, a.windowHeight
}

// SetWindowSize updates the window dimensions
func (a *App) SetWindowSize(width, height unit.Dp) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.windowWidth = width
	a.windowHeight = height

	// Update config
	a.Config.WindowWidth = float32(width)
	a.Config.WindowHeight = float32(height)
}

// HandleResize handles window resize events
func (a *App) HandleResize(width, height int) {
	a.SetWindowSize(unit.Dp(width), unit.Dp(height))
}

// SetAlwaysOnTop sets the always-on-top window option
func (a *App) SetAlwaysOnTop(enabled bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Config.AlwaysOnTop = enabled

	// Note: Gio's always-on-top support varies by platform
	// On some platforms, this may require platform-specific code
	// For now, we store the preference and apply it where supported
	if a.Window != nil {
		// Gio doesn't have a direct always-on-top API in the standard options
		// This would need platform-specific implementation
		// For macOS, this could be done via CGO or external libraries
		// For now, we just store the preference
	}
}

// IsAlwaysOnTop returns whether always-on-top is enabled
func (a *App) IsAlwaysOnTop() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Config.AlwaysOnTop
}

// ToggleAlwaysOnTop toggles the always-on-top setting
func (a *App) ToggleAlwaysOnTop() {
	a.SetAlwaysOnTop(!a.IsAlwaysOnTop())
}

// SaveConfig saves the current configuration
func (a *App) SaveConfig() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Config.Save()
}

// TransitionToCountdown transitions from input screen to countdown screen
func (a *App) TransitionToCountdown() {
	a.SetScreen(ScreenCountdown)
}

// TransitionToInput transitions from countdown screen to input screen
// This also resets the timer
func (a *App) TransitionToInput() {
	a.Controller.Reset()
	a.SetScreen(ScreenInput)
}

// InvalidateWindow requests a window redraw
func (a *App) InvalidateWindow() {
	if a.Window != nil {
		a.Window.Invalidate()
	}
}
