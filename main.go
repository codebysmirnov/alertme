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
	stopTicker := make(chan struct{})      // Канал для остановки текущего таймера
	startNextTicker := make(chan struct{}) // Канал для запуска следующего таймера
	totalRestTime := time.Duration(0)      // Переменная для подсчета общего времени отдыха

	restNotificationText := widget.NewLabel("You need to rest! Press 'Continue button' after the rest")
	continueButton := widget.NewButton("Continue", func() {
		w.Hide()
		stopTicker <- struct{}{}      // Отправляем сигнал для остановки текущего таймера
		startNextTicker <- struct{}{} // Отправляем сигнал для запуска следующего таймера
		log.Println("Continue button pressed")
	})

	totalRestLabel := widget.NewLabel("Total Rest Time: 0s")
	restDurationLabel := widget.NewLabel("Rest Duration: 0s")

	exitButton := widget.NewButton("Exit", func() {
		w.Hide()
		stopTicker <- struct{}{} // Отправляем сигнал для остановки текущего таймера
		log.Println("Exit button pressed")
		dialog.ShowInformation("Rest Timer", "Program interrupted", w)
		os.Exit(0) // Завершаем программу
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
				startTime := time.Now() // Запоминаем время начала отдыха
				w.Show()
				log.Println("Notification appeared")

				go func() {
					for {
						select {
						case <-stopTicker: // Ждем сигнала для остановки текущего таймера
							w.Hide()
							endTime := time.Now()                 // Запоминаем время нажатия кнопки Continue
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
