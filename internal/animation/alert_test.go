package animation

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestGetAnimationState_Thresholds(t *testing.T) {
	testCases := []struct {
		seconds  int
		expected AnimationState
	}{
		{15, Normal},   // > 10 seconds
		{11, Normal},   // > 10 seconds
		{10, Blink},    // 1-10 seconds
		{5, Blink},     // 1-10 seconds
		{1, Blink},     // 1-10 seconds
		{0, Critical},  // <= 0 seconds
		{-5, Critical}, // <= 0 seconds (negative)
		{-100, Critical}, // deep negative
	}

	for _, tc := range testCases {
		result := GetAnimationState(tc.seconds)
		if result != tc.expected {
			t.Errorf("GetAnimationState(%d): expected %s, got %s", tc.seconds, tc.expected, result)
		}
	}
}

func TestAnimationController_BlinkStateNormal(t *testing.T) {
	ac := NewAnimationController()
	
	// Normal state should always return true (always visible)
	for i := 0; i < 10; i++ {
		if !ac.GetBlinkState(Normal) {
			t.Error("GetBlinkState(Normal) should always return true")
		}
	}
}

func TestAnimationController_Reset(t *testing.T) {
	ac := NewAnimationController()
	ac.blinkOn = false
	
	ac.Reset()
	
	if !ac.blinkOn {
		t.Error("Reset should set blinkOn to true")
	}
}

func TestAnimationState_String(t *testing.T) {
	testCases := []struct {
		state    AnimationState
		expected string
	}{
		{Normal, "Normal"},
		{Blink, "Blink"},
		{Critical, "Critical"},
		{AnimationState(99), "Unknown"},
	}

	for _, tc := range testCases {
		result := tc.state.String()
		if result != tc.expected {
			t.Errorf("AnimationState(%d).String(): expected %s, got %s", tc.state, tc.expected, result)
		}
	}
}


// Feature: countdown-timer, Property 9: アニメーション状態の閾値遷移
// **Validates: Requirements 5.4, 5.5**
func TestAnimationState_Property9_ThresholdTransitions(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Any remaining seconds > 10 → Normal state
	properties.Property("remaining seconds > 10 results in Normal state", prop.ForAll(
		func(seconds int) bool {
			state := GetAnimationState(seconds)
			return state == Normal
		},
		gen.IntRange(11, 10000), // seconds > 10
	))

	// Property: Any remaining seconds 1-10 → Blink state
	properties.Property("remaining seconds 1-10 results in Blink state", prop.ForAll(
		func(seconds int) bool {
			state := GetAnimationState(seconds)
			return state == Blink
		},
		gen.IntRange(1, 10), // seconds 1-10 inclusive
	))

	// Property: Any remaining seconds <= 0 → Critical state
	properties.Property("remaining seconds <= 0 results in Critical state", prop.ForAll(
		func(seconds int) bool {
			state := GetAnimationState(seconds)
			return state == Critical
		},
		gen.IntRange(-10000, 0), // seconds <= 0
	))

	// Combined property: Animation state is determined correctly for any integer
	properties.Property("animation state is correctly determined for any remaining seconds", prop.ForAll(
		func(seconds int) bool {
			state := GetAnimationState(seconds)

			if seconds > 10 {
				return state == Normal
			}
			if seconds >= 1 && seconds <= 10 {
				return state == Blink
			}
			// seconds <= 0
			return state == Critical
		},
		gen.IntRange(-10000, 10000),
	))

	properties.TestingRun(t)
}
