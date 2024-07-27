# Graceful Go Service (GGService)
GGService is a Go package designed for building robust and gracefully shutdown services. It provides a framework to easily manage service lifecycle, handle interrupts, and ensure smooth operations even during shutdowns. Ideal for applications requiring reliable service management. Contains graceful shutdowns and custom functions.

go version 1.22+ <br>
[![Go report][go_report_img]][go_report_url]

## Features

- **Graceful Shutdowns:** GgService allows services to handle interrupts and shutdowns gracefully, ensuring minimal disruption.
- **Customizable:** Easily integrate with your Go applications by providing custom start and run functions.
- **Simple API:** Straightforward API for starting, stopping, and managing service lifecycles.

## Installation

To use GGService in your Go project, simply run:

```bash
go get github.com/lmbek/ggservice
```

## Example code (How to use)

```bash
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

// loadService - creates a new service and starts it
func loadService() {
	// starting the service (please note you can choose to not implement any of these by using nil instead)
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)
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

	// if we wish to stop the service, we can do so before the wait function with a waitgroup goroutine
	//service.Stop()
	//service.ForceShutdown()
	
	// we could also start multiple services with waitgroup

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


```

## Contributors
Lars M Bek (https://github.com/lmbek)


## License
[`ggservice`][repos_url] is free and open-source software licensed under the MIT License, created and supported by [Lars M Bek]. 


<!-- Go links -->
[repos_url]: https://github.com/lmbek/ggservice
[go_version_img]: go1.22+
[go_dev_url]: https://pkg.go.dev/github.com/lmbek/ggservice
[go_report_img]: https://goreportcard.com/badge/github.com/lmbek/ggservice
[go_report_url]: https://goreportcard.com/report/github.com/lmbek/ggservice