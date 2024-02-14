package application

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// RestNotifier represents a notifier for rest intervals.
type RestNotifier struct {
	restWindow        fyne.Window   // restWindow holds the Fyne UI restWindow for rest notifications.
	statisticWindow   fyne.Window   // statisticWindow shows statistic about rest
	ticker            *time.Ticker  // ticker is a time.Ticker for triggering rest notifications at regular intervals.
	stopTicker        chan struct{} // stopTicker is a channel for signaling to stop the current rest timer.
	startNextTicker   chan struct{} // startNextTicker is a channel for signaling to start the next rest timer.
	skipTicker        chan struct{} // skipTicker is a channel for signaling to skip the rest notification.
	totalRestTime     time.Duration // totalRestTime holds the accumulated total rest time.
	totalRestLabel    *widget.Label // totalRestLabel is a Fyne widget for displaying the total accumulated rest time.
	restDurationLabel *widget.Label // restDurationLabel is a Fyne widget for displaying the current rest duration.
	intervalMinutes   int           // intervalMinutes is the interval duration in minutes between rest notifications.
}

// newRestNotifier creates a new RestNotifier with the given interval in minutes.
func newRestNotifier(intervalMinutes int) *RestNotifier {
	return &RestNotifier{
		ticker:            time.NewTicker(time.Duration(intervalMinutes) * time.Minute),
		stopTicker:        make(chan struct{}),
		startNextTicker:   make(chan struct{}),
		skipTicker:        make(chan struct{}),
		totalRestTime:     time.Duration(0),
		totalRestLabel:    widget.NewLabel("Total Rest Time: 0s"),
		restDurationLabel: widget.NewLabel("Rest Duration: 0s"),
		intervalMinutes:   intervalMinutes,
	}
}

// restNotifierInit initializes a RestNotifier with the provided parameters.
func restNotifierInit(
	intervalMinutes int,
	restWindow fyne.Window,
	statisticWindow fyne.Window,
) *RestNotifier {
	notifier := newRestNotifier(intervalMinutes)
	notifier.initializeRestWindow(restWindow)
	notifier.setupRestWindowUI()
	notifier.initializeStatisticWindow(statisticWindow)
	notifier.setupRestStatisticUI()

	return notifier
}

// initializeRestWindow sets up the initial properties for the rest restWindow.
func (rn *RestNotifier) initializeRestWindow(window fyne.Window) {
	rn.restWindow = window
	rn.restWindow.Resize(fyne.NewSize(300, 200))
	rn.restWindow.SetFixedSize(true)
	rn.restWindow.CenterOnScreen()
}

// initializeStatisticWindow sets up the initial properties for the statistic window.
func (rn *RestNotifier) initializeStatisticWindow(window fyne.Window) {
	rn.statisticWindow = window
	rn.statisticWindow.Resize(fyne.NewSize(300, 200))
	rn.statisticWindow.SetFixedSize(true)
	rn.statisticWindow.CenterOnScreen()
}

// setupRestStatisticUI configures the user interface components for the statistic window.
func (rn *RestNotifier) setupRestStatisticUI() {
	exitButton := widget.NewButton("Exit", func() {
		log.Println("Rest statistic exit button pressed")
		os.Exit(0)
	})
	rn.statisticWindow.SetContent(container.NewVBox(

		rn.totalRestLabel,
		exitButton,
	))
	rn.statisticWindow.SetCloseIntercept(func() {
		fmt.Println("closing window")
		rn.statisticWindow.Hide()
	})
}

// setupRestWindowUI configures the user interface components.
func (rn *RestNotifier) setupRestWindowUI() {
	continueButton := widget.NewButton("Continue", func() {
		rn.stopTicker <- struct{}{}      // Send a signal to stop the current timer
		rn.startNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Continue button pressed")
	})

	skipButton := widget.NewButton("Skip", func() {
		rn.skipTicker <- struct{}{}      // Send a signal to skip the rest
		rn.startNextTicker <- struct{}{} // Send a signal to start the next timer
		log.Println("Skip button pressed")
	})

	exitButton := widget.NewButton("Exit", func() {
		rn.stopTicker <- struct{}{} // Send a signal to stop the current timer
		log.Println("Exit button pressed")
		dialog.ShowInformation("Rest Timer", "Program interrupted", rn.restWindow)
		os.Exit(0) // Exit the program
	})

	restNotificationText := widget.NewLabel("You need to rest! Press 'Continue button' after the rest")
	rn.restWindow.SetContent(container.NewVBox(
		restNotificationText,
		continueButton,
		skipButton,
		exitButton,
		rn.totalRestLabel,
		rn.restDurationLabel,
	))
}

// showNotification displays the rest notification restWindow.
func (rn *RestNotifier) showNotification() {
	startTime := time.Now() // Record the start time of the rest
	rn.restWindow.Show()
	log.Println("Notification appeared")
	rn.ticker.Stop()

	go func() {
		for {
			select {
			case <-rn.stopTicker: // Wait for a signal to stop the current timer
				rn.restWindow.Hide()
				endTime := time.Now()                 // Record the time when the Continue or Skip button is pressed
				elapsedTime := endTime.Sub(startTime) // Calculate the elapsed time
				rn.totalRestTime += elapsedTime       // Add the elapsed time to the total rest time
				log.Printf("Rest duration: %v\n", elapsedTime)
				log.Printf("Total rest time: %v\n", rn.totalRestTime)
				// Update UI with the total rest time
				rn.totalRestLabel.SetText(fmt.Sprintf("Total Rest Time: %s", rn.totalRestTime.Round(time.Second).String()))
				return
			case <-rn.skipTicker: // Wait for a signal to skip the rest
				rn.restWindow.Hide()
				return
			default:
				elapsed := time.Since(startTime).Round(time.Second)
				// Update UI with the current rest duration
				rn.restDurationLabel.SetText(fmt.Sprintf("Rest Duration: %s", elapsed.String()))
				time.Sleep(time.Second)
			}
		}
	}()
}

// showStatistic displays the statistics window.
func (rn *RestNotifier) showStatistic() {
	if rn.statisticWindow == nil {
		log.Println("The statistics window is not initialized")
		return
	}
	rn.statisticWindow.Show()
}

// run initiates the main loop for rest notifications.
func (rn *RestNotifier) run() {
	for {
		select {
		case <-rn.ticker.C:
			rn.showNotification()
		case <-rn.startNextTicker:
			rn.ticker = time.NewTicker(time.Duration(rn.intervalMinutes) * time.Minute) // Start a new timer
		}
	}
}

// GetTotalRestTime returns the total accumulated rest time.
func (rn *RestNotifier) GetTotalRestTime() time.Duration {
	return rn.totalRestTime
}
