# Graceful Go Service (GGService)
GGService is a Go package designed for building robust and gracefully shutdown services. It provides a framework to easily manage service lifecycle, handle interrupts, and ensure smooth operations even during shutdowns. Ideal for applications requiring reliable service management. Contains graceful shutdowns and custom functions.

## Features

- **Graceful Shutdowns:** GgService allows services to handle interrupts and shutdowns gracefully, ensuring minimal disruption.
- **Customizable:** Easily integrate with your Go applications by providing custom start and run functions.
- **Simple API:** Straightforward API for starting, stopping, and managing service lifecycles.

## Installation

To use GgService in your Go project, simply run:

```bash
go get github.com/your-username/ggservice
```

## Example code (How to use)

```bash
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/your-username/ggservice"
)

func main() {
	fmt.Println("Starting GgService example...")

	// Create a new GgService instance
	service := ggservice.New(&ggservice.Service{Name: "My Service", GracefulShutdownTime: 5 * time.Second})

  // Start the service
	err := service.Start(start, run, forcedStop)
	if err != nil {
		log.Fatal("Failed to start service:", err)
	}

  // you can also start this in a Goroutine and use service.Stop() and service.ForceShutdown() to stop the service internally
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
func forcedStop(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```

## Contributors
Lars M Bek (https://github.com/lmbek)
