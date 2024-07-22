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
	// creating new service with name and graceful shutdown time duration
	service := ggservice.NewService("SSG Service", 5*time.Second)

	// we can use Start, Stop and ForceShutdown on the service

	// since Start is a blocking operation we would want to put Start, Stop or ForceShutdown into a goroutine
	// if we want to stop the service ourselves, otherwise we would also just wait for interrupt
	go func() {
		time.Sleep(1 * time.Second)
		// we can also use service.Stop or service.ForceShutdown, but notice service.Start is a blocking call
		// we can put service.Start in a go routine and use waitgroups etc... your choice.
		//service.Stop()          // waits till all operations are done (not waiting for graceful timer)
		//service.ForceShutdown() // force shutdown immediately
	}()

	// starting the service (please note you can choose to not implement any of these by using nil instead)
	err := service.Start(start, run, forceExit) // this is a blocking call
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
	time.Sleep(10 * time.Second) // note: if the graceful timer duration is below amount of work needed to be done, it will forceExit
	fmt.Println("end of work")
	return nil
}

// forceExit is being run when the application is trying to force a shutdown (non-gracefully)
func forceExit() {
	log.Fatalln(errors.New("forced stop: timeout"))
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