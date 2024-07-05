package main

import (
	"errors"
	"fmt"
	"github.com/lmbek/ggservice"
	"log"
	"time"
)

func main() {
	loadService()
}

func loadService() {
	//alternative //s := new_ssg.NewService("SSG Service", 5*time.Second)
	service := ggservice.New(&ggservice.Service{Name: "My Service", GracefulShutdownTime: 5 * time.Second})
	err := service.Start(start, run, forcedTimeoutStop)
	if err != nil {
		log.Fatal(err)
	}
}

// start runs when the service starts
func start() error {
	fmt.Println("starting function")
	return nil
}

// run loops from the application is started until it is stopped, terminated or ForceShutdown (please use with time.Sleep in between frames)
func run() error {
	fmt.Println("start of work")

	// service will run the rest of the task if graceful shutdown timer is > 1
	time.Sleep(7 * time.Second)
	// service will force shutdown if the rest of the task length
	// is > 5 (or whatever graceful shutdown timer is set to)
	// time.Sleep(8 * time.Second)

	fmt.Println("end of work")
	return nil
}

// forcedStop is being run when the application is trying to force a shutdown (non-gracefully)
func forcedTimeoutStop() {
	log.Fatal(errors.New("forced stop: timeout"))
}
