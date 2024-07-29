package ggservice_test

import (
	"errors"
	"fmt"
	"github.com/lmbek/ggservice"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func ExampleNewService() {
	service := ggservice.NewService("My Service")

	ExampleStart := func() error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampeForceShutdown := func() error {
		fmt.Println("this runs when service forceExits")
		return nil
	}

	ExampeStop := func() error {
		log.Println("running what should happen when we stop program")
		return nil
	}

	err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleNew() {
	service := ggservice.New(&ggservice.Service{Name: "My Service"})

	ExampleStart := func() error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		fmt.Println("this runs if run is set (if loop is not wished then use nil instead)")
		return nil
	}

	ExampeStop := func() error {
		log.Println("running what should happen when we stop program")
		return nil
	}

	ExampeForceShutdown := func() error {
		fmt.Println("this runs when service forceExits")
		return nil
	}

	err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
	if err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	service := ggservice.New(&ggservice.Service{Name: "My Service"})
	if service == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestNewService(t *testing.T) {
	service := ggservice.NewService("My Service")
	if service == nil {
		t.Error("Could not create new service")
	}
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

func TestService_ForceShutdown(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_listenForInterrupt(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func ExampleService_Restart() {
	ExampleStart := func() error {
		fmt.Println("this runs when service starts")
		return nil
	}

	ExampleRun := func() error {
		log.Println("start of work")
		time.Sleep(485 * time.Millisecond) // note: if the graceful timer duration is below amount of work needed to be done, it will forceExit
		log.Println("end of work")
		return nil
	}

	ExampeStop := func() error {
		log.Println("running what should happen when we stop program")
		return nil
	}

	ExampeForceShutdown := func() error {
		log.Println(errors.New("custom forced stop: timeout"))
		os.Exit(-1)
		return nil
	}

	// creating new service with name and graceful shutdown time duration
	service := ggservice.NewService("SSG Service")
	service.SetGracefulShutdownTime(5 * time.Second)
	service.SetLogLevel(ggservice.LOG_LEVEL_ALL)

	// starting the service (please note you can choose to not implement any of these by using nil instead)
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(12)
	go func() {
		fmt.Println("start")
		err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(6 * time.Second)
		fmt.Println("restart")
		err := service.Restart()
		if err != nil {
			log.Println(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(13 * time.Second)
		fmt.Println("stop")
		err := service.Stop()
		if err != nil {
			log.Println("Could not stop service: " + err.Error())
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(17 * time.Second)
		fmt.Println("start")
		err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(20 * time.Second)
		fmt.Println("stop")
		err := service.Stop()
		if err != nil {
			log.Println("Could not stop service: " + err.Error())
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(25 * time.Second)
		fmt.Println("restart")
		err := service.Restart()
		if err != nil {
			log.Println(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(33 * time.Second)
		fmt.Println("stop")
		err := service.Stop()
		if err != nil {
			log.Println("Could not stop service: " + err.Error())
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(38 * time.Second)
		fmt.Println("stop")
		err := service.Stop()
		if err != nil {
			log.Println("Could not stop service: " + err.Error())
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(45 * time.Second)
		fmt.Println("start")
		err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(55 * time.Second)
		fmt.Println("start")
		err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(62 * time.Second)
		fmt.Println("start")
		err := service.Start(ExampleStart, ExampleRun, ExampeStop, ExampeForceShutdown) // this is a blocking call
		if err != nil {
			log.Fatalln(err)
		}
		waitgroup.Done()
	}()

	go func() {
		time.Sleep(70 * time.Second)
		fmt.Println("stop")
		err := service.Stop()
		if err != nil {
			log.Println("Could not stop service: " + err.Error())
		}
		waitgroup.Done()
	}()

	// if we wish to stop the service, we can do so before the wait function in a waitgroup goroutine
	//service.Stop()
	//service.ForceShutdown()

	waitgroup.Wait() // this is a blocking call
}
