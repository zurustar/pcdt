// Package controller contains the business logic controllers for the countdown timer.
package controller

import (
	"errors"
	"sync"
	"time"

	"evangelion-timer/internal/model"
)

// ErrInvalidStartSeconds is returned when trying to start with invalid seconds
var ErrInvalidStartSeconds = errors.New("seconds must be greater than 0")

// TimerController はタイマーの制御を行う
type TimerController struct {
	model    *model.TimerModel
	ticker   *time.Ticker
	stopChan chan struct{}
	onTick   func(state *model.TimerModel)
	mu       sync.Mutex
}

// NewTimerController creates a new TimerController
func NewTimerController() *TimerController {
	return &TimerController{
		model: model.NewTimerModel(),
	}
}

// GetModel returns the underlying timer model
func (tc *TimerController) GetModel() *model.TimerModel {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.model
}

// Start は指定秒数でカウントダウンを開始する
func (tc *TimerController) Start(seconds int) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if seconds <= 0 {
		return ErrInvalidStartSeconds
	}

	// Stop any existing ticker
	tc.stopTickerLocked()

	// Initialize the model
	tc.model.SetInitialSeconds(seconds)
	tc.model.Start()

	// Start the ticker
	tc.startTickerLocked()

	return nil
}


// Stop はカウントダウンを一時停止する
func (tc *TimerController) Stop() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.model.IsRunning() {
		tc.stopTickerLocked()
		tc.model.Pause()
	}
}

// Resume はカウントダウンを再開する
func (tc *TimerController) Resume() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.model.IsPaused() {
		tc.model.Resume()
		tc.startTickerLocked()
	}
}

// Reset はタイマーをリセットする（待機状態に戻す）
func (tc *TimerController) Reset() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.stopTickerLocked()
	tc.model.Reset()
}

// Toggle は一時停止/再開を切り替える
func (tc *TimerController) Toggle() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.model.IsRunning() {
		tc.stopTickerLocked()
		tc.model.Pause()
	} else if tc.model.IsPaused() {
		tc.model.Resume()
		tc.startTickerLocked()
	}
}

// OnTick はティックごとのコールバックを設定する
func (tc *TimerController) OnTick(callback func(state *model.TimerModel)) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.onTick = callback
}

// startTickerLocked starts the internal ticker (must be called with lock held)
func (tc *TimerController) startTickerLocked() {
	tc.stopChan = make(chan struct{})
	tc.ticker = time.NewTicker(1 * time.Second)

	// Capture the ticker channel to avoid race condition with stopTickerLocked
	tickerChan := tc.ticker.C
	stopChan := tc.stopChan

	go func() {
		for {
			select {
			case <-stopChan:
				return
			case <-tickerChan:
				tc.mu.Lock()
				if tc.model.IsRunning() {
					// Decrement remaining seconds (continues into negative)
					tc.model.Tick()

					// Call the onTick callback if set
					if tc.onTick != nil {
						callback := tc.onTick
						state := tc.model
						tc.mu.Unlock()
						callback(state)
						continue
					}
				}
				tc.mu.Unlock()
			}
		}
	}()
}

// stopTickerLocked stops the internal ticker (must be called with lock held)
func (tc *TimerController) stopTickerLocked() {
	if tc.ticker != nil {
		tc.ticker.Stop()
		tc.ticker = nil
	}
	if tc.stopChan != nil {
		close(tc.stopChan)
		tc.stopChan = nil
	}
}
