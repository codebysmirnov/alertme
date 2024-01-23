package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	stopTicker := make(chan struct{})      // Канал для остановки текущего таймера
	startNextTicker := make(chan struct{}) // Канал для запуска следующего таймера
	totalRestTime := time.Duration(0)      // Переменная для подсчета общего времени отдыха

	hello := widget.NewLabel("You need to rest!")
	closeButton := widget.NewButton("Close", func() {
		w.Hide()
		stopTicker <- struct{}{}      // Отправляем сигнал для остановки текущего таймера
		startNextTicker <- struct{}{} // Отправляем сигнал для запуска следующего таймера
		log.Println("Close button pressed")
	})

	totalRestLabel := widget.NewLabel("Total Rest Time: 0s")
	restDurationLabel := widget.NewLabel("Rest Duration: 0s")

	w.SetContent(container.NewVBox(
		hello,
		closeButton,
		totalRestLabel,
		restDurationLabel,
	))

	go func() {
		for {
			select {
			case <-ticker.C:
				startTime := time.Now() // Запоминаем время начала отдыха
				w.Show()
				log.Println("Notification appeared")

				go func() {
					for {
						select {
						case <-stopTicker: // Ждем сигнала для остановки текущего таймера
							w.Hide()
							endTime := time.Now()                 // Запоминаем время нажатия кнопки Close
							elapsedTime := endTime.Sub(startTime) // Вычисляем прошедшее время
							totalRestTime += elapsedTime          // Добавляем прошедшее время к общему времени отдыха
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

				<-startNextTicker                                                     // Ждем сигнала для запуска следующего таймера
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Запускаем новый таймер
			case <-startNextTicker: // Ждем сигнала для запуска следующего таймера
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Запускаем новый таймер
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
