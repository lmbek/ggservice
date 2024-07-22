package main

import (
	"errors"
	"fmt"
	"github.com/lmbek/ggservice"
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {
	loadService()
}

// loadService - creates a new service and starts it
func loadService() {
	// creating new service with name and graceful shutdown time duration
	service := ggservice.NewService("SSG Service", 5*time.Second, true, "starting", "function")

	// starting the service (please note you can choose to not implement any of these by using nil instead)
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)
	go func() {
		err := service.Start(start, run, forceExit) // this is a blocking call
		if err != nil {
			log.Fatal(err)
		}
		waitgroup.Done()
	}()

	// if we wish to stop the service, we can do so before the wait function
	//service.Stop()
	//service.ForceShutdown()

	waitgroup.Wait() // this is a blocking call
}

// start runs when the service starts (note: handling arguments is only needed if arguments are passed along)
func start(args ...any) error {
	argsNeeded := 2
	if len(args) < argsNeeded {
		return errors.New("Need at least " + strconv.Itoa(argsNeeded) + " args when creating new service")
	}
	return startWithArgs(args[0].(string), args[1].(string))
}

func startWithArgs(arg1 string, arg2 string) error {
	fmt.Println(arg1, arg2)
	return nil
}

// run loops from the application is started until it is stopped, terminated or ForceShutdown (please use with time.Sleep in between frames)
func run() error {
	fmt.Println("start of work")
	time.Sleep(10 * time.Second) // note: if the graceful timer duration is below amount of work needed to be done, it will forceExit
	fmt.Println("end of work")
	return nil
}

// forceExit is being run when the application is trying to force a shutdown (non-gracefully)
func forceExit() error {
	log.Fatalln(errors.New("forced stop: timeout"))
	return nil
}
