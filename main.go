package main

import (
	"flag"
	"log"

	"alertme/application"
)

func main() {
	var intervalMinutes int
	flag.IntVar(&intervalMinutes, "interval", 25, "Rest timer interval in minutes")
	flag.Parse()

	log.Printf("Rest timer interval set to %d minutes\n", intervalMinutes)

	alertMe := application.Init(intervalMinutes)

	log.Println("The program successfully started!")

	defer func() {
		log.Printf("Total rest time during the program: %v\n", alertMe.RestNotifier.GetTotalRestTime())
		log.Println("The program successfully terminated!")
	}()

	alertMe.Run()
}
