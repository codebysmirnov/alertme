package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "AlertMe"
)

type RestNotifier struct {
	App               fyne.App
	Window            fyne.Window
	Ticker            *time.Ticker
	StopTicker        chan struct{}
	StartNextTicker   chan struct{}
	SkipTicker        chan struct{}
	TotalRestTime     time.Duration
	TotalRestLabel    *widget.Label
	RestDurationLabel *widget.Label
	IntervalMinutes   int
}

func NewRestNotifier(intervalMinutes int) *RestNotifier {
	return &RestNotifier{
		App:               app.NewWithID(appName),
		Ticker:            time.NewTicker(time.Duration(intervalMinutes) * time.Minute),
		StopTicker:        make(chan struct{}),
		StartNextTicker:   make(chan struct{}),
		SkipTicker:        make(chan struct{}),
		TotalRestTime:     time.Duration(0),
		TotalRestLabel:    widget.NewLabel("Total Rest Time: 0s"),
		RestDurationLabel: widget.NewLabel("Rest Duration: 0s"),
		IntervalMinutes:   intervalMinutes,
	}
}

func (rn *RestNotifier) InitializeWindow() {
	rn.Window = rn.App.NewWindow("Rest time")
	rn.Window.Resize(fyne.NewSize(300, 200))
	rn.Window.SetFixedSize(true)
	rn.Window.CenterOnScreen()
	rn.Window.SetIcon(theme.FyneLogo())
}

func (rn *RestNotifier) SetupSystemTrayMenu() {
	if desk, ok := rn.App.(desktop.App); ok {
		m := fyne.NewMenu("Main",
			fyne.NewMenuItem("Start rest", func() {
				log.Println("Tray: the start notification button pushed")
				rn.ShowNotification()
			}))
		desk.SetSystemTrayMenu(m)
	}
}

func (rn *RestNotifier) SetupUI() {
	continueButton := widget.NewButton("Continue", func() {
		rn.StopTicker <- struct{}{}      // Send a signal to stop the current timer
		rn.StartNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Continue button pressed")
	})

	skipButton := widget.NewButton("Skip", func() {
		rn.SkipTicker <- struct{}{}      // Send a signal to skip the rest
		rn.StartNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Skip button pressed")
	})

	exitButton := widget.NewButton("Exit", func() {
		rn.StopTicker <- struct{}{} // Send a signal to stop the current timer
		log.Println("Exit button pressed")
		dialog.ShowInformation("Rest Timer", "Program interrupted", rn.Window)
		os.Exit(0) // Exit the program
	})

	restNotificationText := widget.NewLabel("You need to rest! Press 'Continue button' after the rest")
	rn.Window.SetContent(container.NewVBox(
		restNotificationText,
		continueButton,
		skipButton,
		exitButton,
		rn.TotalRestLabel,
		rn.RestDurationLabel,
	))
}

func (rn *RestNotifier) ShowNotification() {
	startTime := time.Now() // Record the start time of the rest
	rn.Window.Show()
	log.Println("Notification appeared")
	rn.Ticker.Stop()

	go func() {
		for {
			select {
			case <-rn.StopTicker: // Wait for a signal to stop the current timer
				rn.Window.Hide()
				endTime := time.Now()                 // Record the time when the Continue or Skip button is pressed
				elapsedTime := endTime.Sub(startTime) // Calculate the elapsed time
				rn.TotalRestTime += elapsedTime       // Add the elapsed time to the total rest time
				log.Printf("Rest duration: %v\n", elapsedTime)
				log.Printf("Total rest time: %v\n", rn.TotalRestTime)
				// Update UI with the total rest time
				rn.TotalRestLabel.SetText(fmt.Sprintf("Total Rest Time: %s", rn.TotalRestTime.Round(time.Second).String()))
				return
			case <-rn.SkipTicker: // Wait for a signal to skip the rest
				rn.Window.Hide()
				return
			default:
				elapsed := time.Since(startTime).Round(time.Second)
				// Update UI with the current rest duration
				rn.RestDurationLabel.SetText(fmt.Sprintf("Rest Duration: %s", elapsed.String()))
				time.Sleep(time.Second)
			}
		}
	}()
}

func (rn *RestNotifier) Run() {
	for {
		select {
		case <-rn.Ticker.C:
			rn.ShowNotification()
		case <-rn.StartNextTicker:
			rn.Ticker = time.NewTicker(time.Duration(rn.IntervalMinutes) * time.Minute) // Start a new timer
		}
	}
}

func main() {
	var intervalMinutes int
	flag.IntVar(&intervalMinutes, "interval", 25, "Rest timer interval in minutes")
	flag.Parse()

	log.Printf("Rest timer interval set to %d minutes\n", intervalMinutes)

	rt := NewRestNotifier(intervalMinutes)
	rt.InitializeWindow()
	rt.SetupSystemTrayMenu()
	rt.SetupUI()

	log.Println("Rest timer program successfully started!")

	defer func() {
		log.Printf("Total rest time during the program: %v\n", rt.TotalRestTime)
		log.Println("Rest timer program successfully terminated!")
	}()

	go rt.Run()
	rt.App.Run()
}
