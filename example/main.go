package main

import (
	"github.com/lmbek/ggservice"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	loadService()
}

var ServiceName1 = "My Service 1"
var ServiceName2 = "My Service 2"
var ServiceName3 = "My Service 3"
var ServiceName4 = "My Service 4"

// loadService - creates a new service and starts it
func loadService() {
	// starting the service (please note you can choose to not implement any of these by using nil instead)
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(4)
	go func() {
		// creating new service with name and graceful shutdown time duration
		service := ggservice.NewService(ServiceName1)
		service.SetGracefulShutdownTime(5 * time.Second)
		service.SetLogLevel(ggservice.LOG_LEVEL_INFO)
		err := service.Start(start, run, stop, forceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		// creating new service with name and graceful shutdown time duration
		service := ggservice.NewService(ServiceName2)
		service.SetGracefulShutdownTime(5 * time.Second)
		service.SetLogLevel(ggservice.LOG_LEVEL_INFO)
		err := service.Start(nil, run2, nil, nil)
		if err != nil {
			log.Println(err)
		}
		waitgroup.Done()
	}()

	go func() {
		// creating new service with name and graceful shutdown time duration
		service := ggservice.NewService(ServiceName3)
		service.SetGracefulShutdownTime(5 * time.Second)
		service.SetRunSleepDuration(12 * time.Second)
		service.SetLogLevel(ggservice.LOG_LEVEL_INFO)
		err := service.Start(nil, run3, nil, nil)
		if err != nil {
			log.Println(err)
		}
		waitgroup.Done()
	}()

	go func() {
		// creating new service with name and graceful shutdown time duration
		service := ggservice.NewService(ServiceName4)
		service.SetGracefulShutdownTime(5 * time.Second)
		service.SetRunSleepDuration(12 * time.Second)
		service.SetLogLevel(ggservice.LOG_LEVEL_INFO)
		err := service.Start(nil, nil, nil, nil)
		if err != nil {
			log.Println(err)
		}
		waitgroup.Done()
	}()

	// if we wish to stop the service, we can do so before the wait function with a waitgroup goroutine
	//service.Stop()
	//service.ForceShutdown()

	waitgroup.Wait() // this is a blocking call
}

// start runs when the service starts
func start() error {
	return nil
}

// run loops from the application is started until it is stopped, terminated or ForceShutdown (please use with time.Sleep in between frames)
func run() error {
	log.Println("start of work")
	time.Sleep(485 * time.Millisecond) // note: if the graceful timer duration is below amount of work needed to be done, it will forceExit
	log.Println("end of work")
	return nil
}

// stop is executing what should happen when we stop the program
func stop() error {
	log.Println("running what should happen when we stop program")
	return nil
}

// forceShutdown is being run when the application is trying to force a shutdown (non-gracefully)
func forceShutdown() error {
	log.Println("(Timeout) forced shutdown of program with all its running services")
	os.Exit(-1)
	return nil
}

// run loops from the application is started until it is stopped, terminated or ForceShutdown (please use with time.Sleep in between frames)
func run2() error {
	log.Println("start of work2")
	time.Sleep(894 * time.Millisecond) // note: if the graceful timer duration is below amount of work needed to be done, it will forceExit
	log.Println("end of work2")
	return nil
}

// run loops from the application is started until it is stopped, terminated or ForceShutdown (please use with time.Sleep in between frames)
func run3() error {
	log.Println("work3 (with custom loop time.Sleep)")
	return nil
}
