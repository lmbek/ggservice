package ggservice_test

import (
	"fmt"
	"github.com/lmbek/ggservice"
	"log"
	"testing"
	"time"
)

func ExampleNewService() {
	service := ggservice.NewService("My Service", 5*time.Second, true, nil)

	ExampleStart := func(args ...any) error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampleForceExit := func() error {
		fmt.Println("this runs when service forceExits")
		return nil
	}

	err := service.Start(ExampleStart, ExampleRun, ExampleForceExit)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleNew() {
	service := ggservice.New(&ggservice.Service{Name: "My Service", GracefulShutdownTime: 5 * time.Second, Args: nil})

	ExampleStart := func(args ...any) error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampeForceExit := func() error {
		fmt.Println("this runs when service forceExits")
		return nil
	}

	err := service.Start(ExampleStart, ExampleRun, ExampeForceExit)
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
