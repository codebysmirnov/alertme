package main

import (
	"flag"
	"log"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	hello := widget.NewLabel("You need to rest!")
	closeButton := widget.NewButton("Close", func() {
		w.Hide()
		stopTicker <- struct{}{}      // Отправляем сигнал для остановки текущего таймера
		startNextTicker <- struct{}{} // Отправляем сигнал для запуска следующего таймера
		log.Println("Close button pressed")
	})

	w.SetContent(container.NewVBox(
		hello,
		closeButton,
	))

	go func() {
		for {
			select {
			case <-ticker.C:
				startTime := time.Now() // Запоминаем время начала отдыха
				w.Show()
				log.Println("Notification appeared")
				<-stopTicker // Ждем сигнала для остановки текущего таймера
				w.Hide()
				endTime := time.Now()                 // Запоминаем время нажатия кнопки Close
				elapsedTime := endTime.Sub(startTime) // Вычисляем прошедшее время
				log.Printf("Rest duration: %v\n", elapsedTime)
				ticker.Stop()                                                         // Останавливаем текущий таймер
				<-startNextTicker                                                     // Ждем сигнала для запуска следующего таймера
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Запускаем новый таймер
			case <-startNextTicker: // Ждем сигнала для запуска следующего таймера
				ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute) // Запускаем новый таймер
			}
		}
	}()

	log.Println("Rest timer program successfully started!")

	defer func() {
		log.Println("Rest timer program successfully terminated!")
	}()

	a.Run()
}
