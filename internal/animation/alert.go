// Package animation contains the animation controllers for the countdown timer.
package animation

import (
	"time"
)

// AnimationState represents the current animation state of the timer display.
type AnimationState int

const (
	// Normal is the default animation state for remaining time > 10 seconds.
	Normal AnimationState = iota
	// Blink is the animation state for remaining time 1-10 seconds.
	Blink
	// Critical is the animation state for remaining time <= 0 seconds (overtime).
	Critical
)

// String returns the string representation of the AnimationState.
func (s AnimationState) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Blink:
		return "Blink"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// AnimationController manages the animation state and timing for the countdown timer.
type AnimationController struct {
	// blinkInterval is the duration between blink state changes.
	blinkInterval time.Duration
	// criticalInterval is the duration between critical state changes.
	criticalInterval time.Duration
	// lastToggle tracks the last time the blink state was toggled.
	lastToggle time.Time
	// blinkOn indicates whether the blink is currently in the "on" state.
	blinkOn bool
}

// NewAnimationController creates a new AnimationController with default settings.
func NewAnimationController() *AnimationController {
	return &AnimationController{
		blinkInterval:    500 * time.Millisecond, // Blink every 500ms
		criticalInterval: 250 * time.Millisecond, // Critical blinks faster
		lastToggle:       time.Now(),
		blinkOn:          true,
	}
}

// GetAnimationState returns the appropriate animation state based on remaining seconds.
// - > 10 seconds: Normal animation
// - 1-10 seconds: Blink animation
// - <= 0 seconds: Critical animation (negative/overtime)
func GetAnimationState(remainingSeconds int) AnimationState {
	if remainingSeconds <= 0 {
		return Critical
	}
	if remainingSeconds <= 10 {
		return Blink
	}
	return Normal
}

// GetBlinkState returns whether the display should be visible based on the current
// animation state and timing. This is used for implementing the blink effect.
func (ac *AnimationController) GetBlinkState(state AnimationState) bool {
	if state == Normal {
		return true // Always visible in normal state
	}

	now := time.Now()
	interval := ac.blinkInterval
	if state == Critical {
		interval = ac.criticalInterval
	}

	if now.Sub(ac.lastToggle) >= interval {
		ac.blinkOn = !ac.blinkOn
		ac.lastToggle = now
	}

	return ac.blinkOn
}

// Reset resets the animation controller to its initial state.
func (ac *AnimationController) Reset() {
	ac.lastToggle = time.Now()
	ac.blinkOn = true
}

// SetBlinkInterval sets the interval for blink animation.
func (ac *AnimationController) SetBlinkInterval(interval time.Duration) {
	ac.blinkInterval = interval
}

// SetCriticalInterval sets the interval for critical animation.
func (ac *AnimationController) SetCriticalInterval(interval time.Duration) {
	ac.criticalInterval = interval
}

// GetBlinkProgress returns a value between 0.0 and 1.0 representing the current
// position in the blink cycle. This can be used for smooth animations.
func (ac *AnimationController) GetBlinkProgress(state AnimationState) float64 {
	if state == Normal {
		return 1.0 // Always fully visible
	}

	interval := ac.blinkInterval
	if state == Critical {
		interval = ac.criticalInterval
	}

	elapsed := time.Since(ac.lastToggle)
	progress := float64(elapsed) / float64(interval)
	if progress > 1.0 {
		progress = 1.0
	}

	// Return a value that oscillates between 0 and 1
	if ac.blinkOn {
		return 1.0 - progress
	}
	return progress
}
