// Package controller contains tests for the timer controller.
package controller

import (
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: countdown-timer, Property 3: カウントダウン開始時の初期値設定
// **Validates: Requirements 2.4**
//
// *任意の*有効な入力時間（分、秒）に対して、スタート後のタイマーの残り秒数は
// `分 * 60 + 秒` と等しい。
func TestTimerController_Property3_InitialValueAfterStart(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: After Start(seconds), remaining seconds equals the input seconds
	properties.Property("start sets remaining seconds to input value", prop.ForAll(
		func(seconds int) bool {
			controller := NewTimerController()
			err := controller.Start(seconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()
			result := model.GetRemainingSeconds() == seconds &&
				model.GetInitialSeconds() == seconds &&
				model.IsRunning()

			// Clean up: stop the ticker goroutine
			controller.Reset()
			time.Sleep(10 * time.Millisecond) // Allow goroutine to exit

			return result
		},
		gen.IntRange(1, 5999), // valid seconds: 1秒〜99分59秒
	))

	// Property: Start with minutes and seconds combination sets correct total
	properties.Property("start with minutes and seconds sets correct total", prop.ForAll(
		func(minutes, seconds int) bool {
			totalSeconds := minutes*60 + seconds
			if totalSeconds <= 0 || totalSeconds > 5999 {
				return true // skip invalid combinations
			}

			controller := NewTimerController()
			err := controller.Start(totalSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()
			result := model.GetRemainingSeconds() == totalSeconds

			// Clean up: stop the ticker goroutine
			controller.Reset()
			time.Sleep(10 * time.Millisecond) // Allow goroutine to exit

			return result
		},
		gen.IntRange(0, 99), // minutes: 0-99
		gen.IntRange(0, 59), // seconds: 0-59
	))

	// Property: Start with invalid seconds (0 or negative) returns error
	properties.Property("start with invalid seconds returns error", prop.ForAll(
		func(seconds int) bool {
			controller := NewTimerController()
			err := controller.Start(seconds)
			// No cleanup needed - Start fails before creating goroutine
			return err == ErrInvalidStartSeconds
		},
		gen.IntRange(-100, 0), // invalid seconds: negative and zero
	))

	properties.TestingRun(t)
}

// Feature: countdown-timer, Property 5: 停止・再開のラウンドトリップ
// **Validates: Requirements 3.2, 3.3**
//
// *任意の*実行中のタイマーに対して、停止してから再開すると、
// 停止直前の残り秒数からカウントダウンが継続される。
func TestTimerController_Property5_StopResumeRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Stop then Resume preserves the remaining seconds value
	// This test verifies the state management without relying on the ticker goroutine
	properties.Property("stop and resume preserves remaining seconds", prop.ForAll(
		func(initialSeconds int, tickCount int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			// Simulate some ticks by directly manipulating the model
			model := controller.GetModel()
			for i := 0; i < tickCount; i++ {
				model.Tick()
			}

			// Record the remaining seconds before stop
			remainingBeforeStop := model.GetRemainingSeconds()

			// Stop the timer
			controller.Stop()
			if !model.IsPaused() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Verify remaining seconds unchanged after stop
			remainingAfterStop := model.GetRemainingSeconds()
			if remainingAfterStop != remainingBeforeStop {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Resume the timer
			controller.Resume()
			if !model.IsRunning() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Verify remaining seconds unchanged after resume
			remainingAfterResume := model.GetRemainingSeconds()
			result := remainingAfterResume == remainingBeforeStop

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // initialSeconds: 1秒〜99分59秒
		gen.IntRange(0, 100),  // tickCount: 0〜100回のティック
	))

	// Property: Stop preserves remaining seconds (simpler test without Resume)
	properties.Property("stop preserves remaining seconds value", prop.ForAll(
		func(initialSeconds int, tickCount int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()
			// Simulate some ticks
			for i := 0; i < tickCount; i++ {
				model.Tick()
			}

			remainingBefore := model.GetRemainingSeconds()

			// Stop the timer
			controller.Stop()

			// Verify state
			result := model.IsPaused() && model.GetRemainingSeconds() == remainingBefore

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // initialSeconds: 1秒〜99分59秒
		gen.IntRange(0, 100),  // tickCount: 0〜100回のティック
	))

	properties.TestingRun(t)
}

// Feature: countdown-timer, Property 6: ゼロ到達後のマイナス継続
// **Validates: Requirements 4.1**
//
// *任意の*カウントダウンに対して、ゼロに到達した後も停止せず、
// 1秒後に残り秒数が-1になる。
func TestTimerController_Property6_ContinuesIntoNegative(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Timer continues into negative after reaching zero
	properties.Property("timer continues into negative after zero", prop.ForAll(
		func(initialSeconds int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// Tick until we reach zero
			for i := 0; i < initialSeconds; i++ {
				model.Tick()
			}

			// Verify we're at zero
			if model.GetRemainingSeconds() != 0 {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Timer should still be running at zero
			if !model.IsRunning() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// One more tick should go negative
			model.Tick()
			if model.GetRemainingSeconds() != -1 {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Timer should still be running in negative
			result := model.IsRunning()

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 100), // initialSeconds: 1〜100秒（テスト効率のため小さめ）
	))

	// Property: Timer continues counting down in negative territory
	properties.Property("timer continues counting in negative territory", prop.ForAll(
		func(initialSeconds int, extraTicks int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// Tick past zero
			totalTicks := initialSeconds + extraTicks
			for i := 0; i < totalTicks; i++ {
				model.Tick()
			}

			expectedRemaining := initialSeconds - totalTicks
			result := model.GetRemainingSeconds() == expectedRemaining &&
				model.IsRunning() &&
				model.IsNegative()

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 50), // initialSeconds: 1〜50秒
		gen.IntRange(1, 50), // extraTicks: 1〜50回の追加ティック
	))

	// Property: IsNegative returns true only when remaining seconds < 0
	properties.Property("IsNegative correctly identifies negative state", prop.ForAll(
		func(initialSeconds int, tickCount int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			for i := 0; i < tickCount; i++ {
				model.Tick()
			}

			remaining := model.GetRemainingSeconds()
			result := model.IsNegative() == (remaining < 0)

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 100), // initialSeconds: 1〜100秒
		gen.IntRange(0, 200), // tickCount: 0〜200回のティック
	))

	properties.TestingRun(t)
}

// Feature: countdown-timer, Property 10: スペースキーによる状態トグル
// **Validates: Requirements 7.2**
//
// *任意の*カウントダウン画面の状態（実行中または一時停止中）に対して、
// スペースキーを押すと状態が反転する（実行中→一時停止、一時停止→実行中）。
func TestTimerController_Property10_SpaceKeyToggle(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Toggle on running timer changes state to paused
	properties.Property("toggle on running timer changes to paused", prop.ForAll(
		func(initialSeconds int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// Verify timer is running
			if !model.IsRunning() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Toggle (simulates space key press)
			controller.Toggle()

			// Verify timer is now paused
			result := model.IsPaused()

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // valid seconds: 1秒〜99分59秒
	))

	// Property: Toggle on paused timer changes state to running
	properties.Property("toggle on paused timer changes to running", prop.ForAll(
		func(initialSeconds int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// First stop the timer to get it into paused state
			controller.Stop()
			if !model.IsPaused() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Toggle (simulates space key press)
			controller.Toggle()

			// Verify timer is now running
			result := model.IsRunning()

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // valid seconds: 1秒〜99分59秒
	))

	// Property: Double toggle returns to original state (running -> paused -> running)
	properties.Property("double toggle returns to original running state", prop.ForAll(
		func(initialSeconds int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// Verify initial running state
			if !model.IsRunning() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// First toggle: running -> paused
			controller.Toggle()
			if !model.IsPaused() {
				controller.Reset()
				time.Sleep(10 * time.Millisecond)
				return false
			}

			// Second toggle: paused -> running
			controller.Toggle()
			result := model.IsRunning()

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // valid seconds: 1秒〜99分59秒
	))

	// Property: Toggle preserves remaining seconds value
	properties.Property("toggle preserves remaining seconds", prop.ForAll(
		func(initialSeconds int, tickCount int) bool {
			controller := NewTimerController()
			err := controller.Start(initialSeconds)
			if err != nil {
				return false
			}

			model := controller.GetModel()

			// Simulate some ticks
			for i := 0; i < tickCount; i++ {
				model.Tick()
			}

			remainingBefore := model.GetRemainingSeconds()

			// Toggle to pause
			controller.Toggle()
			remainingAfterPause := model.GetRemainingSeconds()

			// Toggle to resume
			controller.Toggle()
			remainingAfterResume := model.GetRemainingSeconds()

			// Remaining seconds should be preserved through both toggles
			result := remainingAfterPause == remainingBefore && remainingAfterResume == remainingBefore

			// Clean up
			controller.Reset()
			time.Sleep(10 * time.Millisecond)

			return result
		},
		gen.IntRange(1, 5999), // initialSeconds: 1秒〜99分59秒
		gen.IntRange(0, 100),  // tickCount: 0〜100回のティック
	))

	properties.TestingRun(t)
}
