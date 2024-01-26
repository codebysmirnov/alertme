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

func main() {
	var intervalMinutes int
	flag.IntVar(&intervalMinutes, "interval", 25, "Rest timer interval in minutes")
	flag.Parse()

	log.Printf("Rest timer interval set to %d minutes\n", intervalMinutes)

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	stopTicker := make(chan struct{})      // Channel to stop the current timer
	startNextTicker := make(chan struct{}) // Channel to start the next timer
	skipTicker := make(chan struct{})      // Channel to skip the rest
	totalRestTime := time.Duration(0)      // Variable to count the total rest time

	App := app.NewWithID(appName)
	// enable the tray menu if this is a desktop application
	if desk, ok := App.(desktop.App); ok {
		m := fyne.NewMenu("Main",
			fyne.NewMenuItem("Start rest", func() {
				log.Println("rest started")
			}))
		desk.SetSystemTrayMenu(m)
	}

	w := App.NewWindow("Rest time")
	w.Resize(fyne.NewSize(300, 200))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetIcon(theme.FyneLogo())

	continueButton := widget.NewButton("Continue", func() {
		w.Hide()
		stopTicker <- struct{}{}      // Send a signal to stop the current timer
		startNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Continue button pressed")
	})

	skipButton := widget.NewButton("Skip", func() {
		w.Hide()
		skipTicker <- struct{}{}      // Send a signal to skip the rest
		startNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Skip button pressed")
	})

	totalRestLabel := widget.NewLabel("Total Rest Time: 0s")
	restDurationLabel := widget.NewLabel("Rest Duration: 0s")

	exitButton := widget.NewButton("Exit", func() {
		w.Hide()
		stopTicker <- struct{}{} // Send a signal to stop the current timer
		log.Println("Exit button pressed")
		dialog.ShowInformation("Rest Timer", "Program interrupted", w)
		os.Exit(0) // Exit the program
	})

	restNotificationText := widget.NewLabel("You need to rest! Press 'Continue button' after the rest")
	w.SetContent(container.NewVBox(
		restNotificationText,
		continueButton,
		skipButton,
		exitButton,
		totalRestLabel,
		restDurationLabel,
	))

	go func() {
		for {
			select {
			case <-ticker.C:
				startTime := time.Now() // Record the start time of the rest
				w.Show()
				log.Println("Notification appeared")

				go func() {
					for {
						select {
						case <-stopTicker: // Wait for a signal to stop the current timer
							w.Hide()
							endTime := time.Now()                 // Record the time when the Continue or Skip button is pressed
							elapsedTime := endTime.Sub(startTime) // Calculate the elapsed time
							totalRestTime += elapsedTime          // Add the elapsed time to the total rest time
							log.Printf("Rest duration: %v\n", elapsedTime)
							log.Printf("Total rest time: %v\n", totalRestTime)
							totalRestLabel.SetText(fmt.Sprintf("Total Rest Time: %s", totalRestTime.Round(time.Second).String()))
							return
						case <-skipTicker: // Wait for a signal to skip the rest
							w.Hide()
							return
						default:
							elapsed := time.Since(startTime).Round(time.Second)
							restDurationLabel.SetText(fmt.Sprintf("Rest Duration: %s", elapsed.String()))
							time.Sleep(time.Second)
						}
					}
				}()

				<-startNextTicker                                                     // Wait for a signal to start the next timer
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Start a new timer
			case <-startNextTicker: // Wait for a signal to start the next timer
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Start a new timer
			}
		}
	}()

	log.Println("Rest timer program successfully started!")

	defer func() {
		log.Printf("Total rest time during the program: %v\n", totalRestTime)
		log.Println("Rest timer program successfully terminated!")
	}()

	App.Run()
}
