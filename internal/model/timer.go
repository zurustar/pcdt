// Package model contains the data models for the countdown timer.
package model

import "time"

// TimerStatus はタイマーの状態を表す列挙型
type TimerStatus int

const (
	StatusIdle    TimerStatus = iota // 待機中
	StatusRunning                    // 実行中
	StatusPaused                     // 一時停止中
)

// String returns the string representation of TimerStatus
func (s TimerStatus) String() string {
	switch s {
	case StatusIdle:
		return "Idle"
	case StatusRunning:
		return "Running"
	case StatusPaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

// TimerModel はタイマーの状態を保持する
type TimerModel struct {
	InitialSeconds   int         // 初期設定秒数
	RemainingSeconds int         // 残り秒数（負の値も可）
	Status           TimerStatus // タイマーの状態
	StartedAt        time.Time   // 開始時刻
	PausedAt         time.Time   // 一時停止時刻
}

// NewTimerModel creates a new TimerModel with default values
func NewTimerModel() *TimerModel {
	return &TimerModel{
		InitialSeconds:   0,
		RemainingSeconds: 0,
		Status:           StatusIdle,
	}
}

// GetRemainingSeconds は残り秒数を返す（負の値も可）
func (t *TimerModel) GetRemainingSeconds() int {
	return t.RemainingSeconds
}

// GetInitialSeconds は初期設定秒数を返す
func (t *TimerModel) GetInitialSeconds() int {
	return t.InitialSeconds
}

// IsRunning はタイマーが実行中かを返す
func (t *TimerModel) IsRunning() bool {
	return t.Status == StatusRunning
}

// IsPaused はタイマーが一時停止中かを返す
func (t *TimerModel) IsPaused() bool {
	return t.Status == StatusPaused
}

// IsIdle はタイマーが待機中かを返す
func (t *TimerModel) IsIdle() bool {
	return t.Status == StatusIdle
}

// IsNegative は超過状態（残り秒数が負）かを返す
func (t *TimerModel) IsNegative() bool {
	return t.RemainingSeconds < 0
}

// SetInitialSeconds sets the initial countdown seconds
func (t *TimerModel) SetInitialSeconds(seconds int) {
	t.InitialSeconds = seconds
	t.RemainingSeconds = seconds
}

// Start starts the timer
func (t *TimerModel) Start() {
	t.Status = StatusRunning
	t.StartedAt = time.Now()
}

// Pause pauses the timer
func (t *TimerModel) Pause() {
	t.Status = StatusPaused
	t.PausedAt = time.Now()
}

// Resume resumes the timer from paused state
func (t *TimerModel) Resume() {
	t.Status = StatusRunning
}

// Reset resets the timer to idle state
func (t *TimerModel) Reset() {
	t.Status = StatusIdle
	t.InitialSeconds = 0
	t.RemainingSeconds = 0
	t.StartedAt = time.Time{}
	t.PausedAt = time.Time{}
}

// Tick decrements the remaining seconds by 1
func (t *TimerModel) Tick() {
	if t.Status == StatusRunning {
		t.RemainingSeconds--
	}
}
