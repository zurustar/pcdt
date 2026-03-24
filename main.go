// Package main is the entry point for the countdown timer application.
package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"

	"countdown-timer/internal/controller"
	"countdown-timer/internal/model"
	"countdown-timer/internal/view"
)

func main() {
	go func() {
		if err := run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run() error {
	// Initialize theme
	theme := view.NewAlertTheme()

	// Initialize timer controller
	timerController := controller.NewTimerController()

	// Initialize application with theme and controller
	application := view.NewApp(
		view.WithTheme(theme),
		view.WithController(timerController),
	)

	// Create the window
	application.CreateWindow()

	// Create material theme for Gio widgets
	materialTheme := material.NewTheme()

	// Create screens
	inputScreen := view.NewInputScreen(theme)
	countdownScreen := view.NewCountdownScreen(theme, timerController)

	// Create help dialog
	helpDialog := view.NewHelpDialog(theme)

	// Set up screen transition callbacks
	inputScreen.SetOnStart(func(minutes, seconds int) {
		totalSeconds := minutes*60 + seconds
		if err := timerController.Start(totalSeconds); err != nil {
			log.Printf("Failed to start timer: %v", err)
			return
		}
		application.TransitionToCountdown()
	})

	countdownScreen.SetOnReset(func() {
		application.TransitionToInput()
		inputScreen.Reset()
		countdownScreen.Reset()
	})

	// Set up timer tick callback to invalidate window for redraw
	timerController.OnTick(func(state *model.TimerModel) {
		application.InvalidateWindow()
	})

	// Main event loop
	var ops op.Ops
	for {
		switch e := application.Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle global keyboard events (help dialog toggle)
			handleGlobalKeys(gtx, helpDialog)

			// Layout based on current screen
			switch application.CurrentScreen() {
			case view.ScreenInput:
				inputScreen.Layout(gtx, materialTheme)
			case view.ScreenCountdown:
				countdownScreen.Layout(gtx, materialTheme)
			}

			// Overlay help dialog if visible
			if helpDialog.IsVisible() {
				helpDialog.Layout(gtx, materialTheme)
			}

			e.Frame(gtx.Ops)

		case app.ConfigEvent:
			// Handle window configuration changes (resize, etc.)
			// The window size is automatically handled by Gio

		}
	}
}

// handleGlobalKeys handles global keyboard shortcuts (H/? for help)
func handleGlobalKeys(gtx layout.Context, helpDialog *view.HelpDialog) {
	for {
		event, ok := gtx.Event(
			key.Filter{Name: "H"},
			key.Filter{Name: "?"},
		)
		if !ok {
			break
		}
		if e, ok := event.(key.Event); ok && e.State == key.Press {
			helpDialog.Toggle()
		}
	}
}
