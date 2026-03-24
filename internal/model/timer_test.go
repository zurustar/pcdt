// Package model contains tests for the timer model.
package model

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: countdown-timer, Property 4: タイマー状態に基づく更新動作
// **Validates: Requirements 3.1, 3.5**
//
// *任意の*タイマー状態に対して、実行中であれば1秒後に残り秒数が1減少し、
// 一時停止中であれば残り秒数は変化しない。
func TestTimerModel_Property4_TickBehaviorBasedOnState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: When timer is running, Tick() decrements RemainingSeconds by 1
	properties.Property("running timer decrements by 1 on tick", prop.ForAll(
		func(initialSeconds int, remainingSeconds int) bool {
			timer := NewTimerModel()
			timer.SetInitialSeconds(initialSeconds)
			timer.RemainingSeconds = remainingSeconds
			timer.Start()

			beforeTick := timer.RemainingSeconds
			timer.Tick()
			afterTick := timer.RemainingSeconds

			return afterTick == beforeTick-1
		},
		gen.IntRange(1, 5999),  // initialSeconds: 1秒〜99分59秒
		gen.IntRange(-100, 5999), // remainingSeconds: マイナス値も含む
	))

	// Property: When timer is paused, Tick() does not change RemainingSeconds
	properties.Property("paused timer does not change on tick", prop.ForAll(
		func(initialSeconds int, remainingSeconds int) bool {
			timer := NewTimerModel()
			timer.SetInitialSeconds(initialSeconds)
			timer.RemainingSeconds = remainingSeconds
			timer.Start()
			timer.Pause()

			beforeTick := timer.RemainingSeconds
			timer.Tick()
			afterTick := timer.RemainingSeconds

			return afterTick == beforeTick
		},
		gen.IntRange(1, 5999),  // initialSeconds: 1秒〜99分59秒
		gen.IntRange(-100, 5999), // remainingSeconds: マイナス値も含む
	))

	// Property: When timer is idle, Tick() does not change RemainingSeconds
	properties.Property("idle timer does not change on tick", prop.ForAll(
		func(initialSeconds int, remainingSeconds int) bool {
			timer := NewTimerModel()
			timer.SetInitialSeconds(initialSeconds)
			timer.RemainingSeconds = remainingSeconds
			// Timer stays in Idle state (not started)

			beforeTick := timer.RemainingSeconds
			timer.Tick()
			afterTick := timer.RemainingSeconds

			return afterTick == beforeTick
		},
		gen.IntRange(1, 5999),  // initialSeconds: 1秒〜99分59秒
		gen.IntRange(-100, 5999), // remainingSeconds: マイナス値も含む
	))

	properties.TestingRun(t)
}
