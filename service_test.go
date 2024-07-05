package ggservice_test

import (
	"fmt"
	"ggservice"
	"log"
	"testing"
	"time"
)

func ExampleNewService() {
	service := ggservice.NewService("My Service", 5*time.Second)

	ExampleStart := func() error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampleForcedTimeoutStop := func() {
		fmt.Println("this runs when service starts")
	}

	err := service.Start(ExampleStart, ExampleRun, ExampleForcedTimeoutStop)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleNew() {
	service := ggservice.New(&ggservice.Service{Name: "My Service", GracefulShutdownTime: 5 * time.Second})

	ExampleStart := func() error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampleForcedTimeoutStop := func() {
		fmt.Println("this runs when service starts")
	}

	err := service.Start(ExampleStart, ExampleRun, ExampleForcedTimeoutStop)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestNewService(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_ForceShutdown(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_Start(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_Stop(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_listenForInterrupt(t *testing.T) {
	t.Errorf("test not implemented yet")
}
