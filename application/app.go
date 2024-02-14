package application

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

const (
	appName = "AlertMe"
)

// AlertMe represents the main application structure.
type AlertMe struct {
	App          fyne.App      // App holds the Fyne application instance.
	RestNotifier *RestNotifier // RestNotifier is an instance of RestNotifier responsible for rest notifications.
}

// newApplication creates and initializes a new AlertMe instance.
func newApplication() AlertMe {
	return AlertMe{
		App: app.NewWithID(appName),
	}
}

// Run starts the AlertMe application and runs the Fyne app.
func (a *AlertMe) Run() {
	go a.RestNotifier.run() // Start the rest notifier in a separate goroutine.
	a.App.Run()             // Run the Fyne app.
}

// Init initializes the AlertMe application with the specified rest interval in minutes.
func Init(intervalMinutes int) AlertMe {
	a := newApplication()
	a.setupSystemTrayMenu()

	a.RestNotifier = restNotifierInit(
		intervalMinutes,
		a.App.NewWindow("Rest time"),
		a.App.NewWindow("Rest statistic"),
	)

	return a
}

// setTrayMenuIcon sets the system tray menu icon for the desktop app.
func setTrayMenuIcon(da desktop.App) {
	ico, err := fyne.LoadResourceFromPath("./assets/main-tray-ico.png")
	if err != nil {
		log.Printf("load tray icon failed %s", err)
		return
	}
	da.SetSystemTrayIcon(ico)
}

// setupSystemTrayMenu configures the system tray menu for the AlertMe application.
func (a *AlertMe) setupSystemTrayMenu() {
	if desk, ok := a.App.(desktop.App); ok {
		setTrayMenuIcon(desk)

		m := fyne.NewMenu("Main",
			fyne.NewMenuItem("Start rest", func() {
				log.Println("Tray: the start notification button pushed")
				a.RestNotifier.showNotification()
			}),
			fyne.NewMenuItem("Statistic", func() {
				log.Println("Tray: the statistic button pushed")
				a.RestNotifier.showStatistic()
			}),
		)

		desk.SetSystemTrayMenu(m)
	}
}
