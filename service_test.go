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

// TESTS:

var startFunc = func() error {
	fmt.Println("started service...")
	return nil
}

var runFunc = func() error {
	fmt.Println("running service...")
	time.Sleep(1 * time.Second)
	return nil
}

var stopFunc = func() error {
	fmt.Println("stopped service...")
	time.Sleep(1 * time.Second)
	return nil
}

func TestNew(test *testing.T) {
	service := ggservice.New(&ggservice.Service{Name: "My Service"})
	if service == nil {
		test.Error("Could not create new service")
	}
}

func TestNewService(t *testing.T) {
	service := ggservice.NewService("My Service")
	if service == nil {
		t.Error("Could not create new service")
	}
}

func TestService_StartStop(t *testing.T) {
	t.Run("Without custom functions", func(t *testing.T) {
		service := ggservice.NewService("My Service")
		err := service.Start(nil, nil, nil, nil)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("With start function", func(t *testing.T) {
		service := ggservice.NewService("My Service")

		err := service.Start(startFunc, nil, nil, nil)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("With stop function", func(t *testing.T) {
		service := ggservice.NewService("My Service")

		err := service.Start(nil, nil, stopFunc, nil)
		if err != nil {
			t.Error(err)
		}
	})
	// waitgroup dependent below:
	t.Run("With custom functions", func(t *testing.T) {
		waitgroup := sync.WaitGroup{}
		waitgroup.Add(2)
		var service ggservice.IService

		go func() {
			defer waitgroup.Done()
			service = ggservice.NewService("My Service")
			err := service.Start(nil, runFunc, nil, nil)
			if err != nil {
				t.Error(err)
			}
		}()
		go func() {
			defer waitgroup.Done()
			time.Sleep(3 * time.Second)
			err := service.Stop()
			if err != nil {
				t.Error(err)
			}
		}()
		waitgroup.Wait()
	})
	t.Run("With full custom functions", func(t *testing.T) {
		waitgroup := sync.WaitGroup{}
		waitgroup.Add(2)
		var service2 ggservice.IService

		forceShutdownFunc := func() error {
			fmt.Println("stopped service...")
			time.Sleep(1 * time.Second)
			return nil
		}

		go func() {
			defer waitgroup.Done()
			service2 = ggservice.NewService("My Service")
			err := service2.Start(startFunc, runFunc, stopFunc, forceShutdownFunc)
			if err != nil {
				t.Error(err)
			}
		}()

		go func() {
			defer waitgroup.Done()
			time.Sleep(3 * time.Second)
			err := service2.Stop()
			if err != nil {
				t.Error(err)
			}
		}()
		waitgroup.Wait()
	})
}

func TestService_Restart(t *testing.T) {
	service := ggservice.NewService("My Service")

	testFunc := func() {
		waitgroup := sync.WaitGroup{}

		go func() {
			waitgroup.Add(1)
			defer waitgroup.Done()
			time.Sleep(3 * time.Second)
			err := service.Restart() // this is a blocking call
			if err != nil {
				t.Error(err)
			}
		}()
		go func() {
			waitgroup.Add(1)
			defer waitgroup.Done()
			time.Sleep(6 * time.Second)
			err := service.Stop()
			if err != nil {
				t.Error(err)
			}
		}()
		err := service.Start(nil, runFunc, nil, nil) // this is a blocking call
		waitgroup.Wait()
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("with_logLevel_default", func(t *testing.T) {
		testFunc()
	})

	t.Run("with_logLevel_all", func(t *testing.T) {
		service.SetLogLevel(ggservice.LOG_LEVEL_ALL)
		testFunc()
	})

	t.Run("with_logLevel_warn", func(t *testing.T) {
		service.SetLogLevel(ggservice.LOG_LEVEL_WARN)
		testFunc()
	})

	t.Run("with_logLevel_error", func(t *testing.T) {
		service.SetLogLevel(ggservice.LOG_LEVEL_ERROR)
		testFunc()
	})

	t.Run("with_logLevel_none", func(t *testing.T) {
		service.SetLogLevel(ggservice.LOG_LEVEL_NONE)
		testFunc()
	})

}

func TestService_ForceShutdown(t *testing.T) {
	t.Errorf("test not implemented yet")
}

func TestService_listenForInterrupt(t *testing.T) {
	t.Errorf("test not implemented yet")
}

// EXAMPLES:

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
