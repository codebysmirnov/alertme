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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var intervalMinutes int
	flag.IntVar(&intervalMinutes, "interval", 25, "Rest timer interval in minutes")
	flag.Parse()

	log.Printf("Rest timer interval set to %d minutes\n", intervalMinutes)

	a := app.New()
	w := a.NewWindow("Rest time")

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	stopTicker := make(chan struct{})      // Channel to stop the current timer
	startNextTicker := make(chan struct{}) // Channel to start the next timer
	totalRestTime := time.Duration(0)      // Variable to count the total rest time

	restNotificationText := widget.NewLabel("You need to rest! Press 'Continue button' after the rest")
	continueButton := widget.NewButton("Continue", func() {
		w.Hide()
		stopTicker <- struct{}{}      // Send a signal to stop the current timer
		startNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Continue button pressed")
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

	w.SetContent(container.NewVBox(
		restNotificationText,
		continueButton,
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
							endTime := time.Now()                 // Record the time when the Continue button is pressed
							elapsedTime := endTime.Sub(startTime) // Calculate the elapsed time
							totalRestTime += elapsedTime          // Add the elapsed time to the total rest time
							log.Printf("Rest duration: %v\n", elapsedTime)
							log.Printf("Total rest time: %v\n", totalRestTime)
							totalRestLabel.SetText(fmt.Sprintf("Total Rest Time: %s", totalRestTime.Round(time.Second).String()))
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

	w.Resize(fyne.NewSize(300, 200))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetIcon(theme.FyneLogo())

	log.Println("Rest timer program successfully started!")

	defer func() {
		log.Printf("Total rest time during the program: %v\n", totalRestTime)
		log.Println("Rest timer program successfully terminated!")
	}()

	a.Run()
}
