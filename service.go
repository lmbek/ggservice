package ggservice

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// IService defines the interface for managing a service with start, stop, and force shutdown capabilities.
type IService interface {
	Start(startFunc func() error, runFunc func() error, stopFunc func() error, forceShutdownFunc func() error) error
	Restart() error
	Stop() error
	ForceShutdown() error
	GetIsRunning() bool
	GetGracefulShutdownTime() time.Duration
	SetGracefulShutdownTime(gracefulShutdownTime time.Duration)
	GetRunSleepDuration() time.Duration
	SetRunSleepDuration(runSleepDuration time.Duration)
	GetLogLevel() int
	SetLogLevel(logLevel int)
}

// Service represents a service that can be started, stopped, and forcefully shutdown with graceful handling.
type Service struct {
	Name                            string        // Name of the service
	gracefulShutdownTime            time.Duration // Timeout duration for graceful shutdown
	isRunning                       bool          // Flag indicating whether the service is running
	canRestart                      bool
	runSleepDuration                time.Duration
	logLevel                        int
	isInitialized                   bool
	isListenForInterruptInitialized bool
	isInterrupted                   bool
	customFunctions                 [4]func() error
}

// Log levels
const (
	LOG_LEVEL_NONE  = iota // 0: No logging
	LOG_LEVEL_ERROR        // 1: Log errors
	LOG_LEVEL_WARN         // 2: Log warnings and errors
	LOG_LEVEL_INFO         // 3: Log info, warnings, and errors
	LOG_LEVEL_ALL          // 4: Log all (info, warnings, and errors)
)

// New creates a new instance of Service with the given name and graceful shutdown timeout.
func New(service *Service) IService {
	return &Service{
		Name:                 service.Name,
		gracefulShutdownTime: 5 * time.Second,
		logLevel:             LOG_LEVEL_ALL,
		isRunning:            true,
		canRestart:           true, // can only be changed by the program
	}
}

// NewService creates a new instance of Service with the given name and graceful shutdown timeout.
func NewService(name string) IService {
	return New(&Service{Name: name})
}

func (s *Service) GetGracefulShutdownTime() time.Duration {
	return s.gracefulShutdownTime
}

func (s *Service) SetGracefulShutdownTime(gracefulShutdownTime time.Duration) {
	s.gracefulShutdownTime = gracefulShutdownTime
}

func (s *Service) SetRunSleepDuration(runSleepDuration time.Duration) {
	s.runSleepDuration = runSleepDuration
}

func (s *Service) GetRunSleepDuration() time.Duration {
	return s.runSleepDuration
}

func (s *Service) GetLogLevel() int {
	return s.logLevel
}

func (s *Service) SetLogLevel(logLevel int) {
	s.logLevel = logLevel
}

func (s *Service) GetIsRunning() bool {
	return s.isInitialized && s.isRunning && s.canRestart
}

// Start starts the service with custom start, run, and stop functions.
func (s *Service) Start(startFunc func() error, runFunc func() error, stopFunc func() error, forceShutdownFunc func() error) error {
	if s.isInitialized {
		if s.logLevel >= LOG_LEVEL_WARN {
			time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
			log.Println("Already started")
		}
		return nil
	}

	s.isInitialized = true
	s.isRunning = true
	s.canRestart = false
	s.customFunctions[0] = startFunc
	s.customFunctions[1] = runFunc
	s.customFunctions[2] = stopFunc
	s.customFunctions[3] = forceShutdownFunc

	if s.logLevel >= LOG_LEVEL_INFO {
		time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
		log.Printf("Starting service: %s\n", s.Name)
	}

	// Custom start function if provided
	if startFunc != nil {
		err := startFunc()
		if err != nil {
			return err
		}
	} else {
		// do nothing
	}

	// Custom run func if provided (in a loop as long as the service is running)
	if runFunc != nil {
		// listen for interrupts for running service
		if !s.isListenForInterruptInitialized {
			s.isListenForInterruptInitialized = true
			go s.listenForInterrupt(forceShutdownFunc) // Listen for interrupt signals
		}

		for s.isRunning {
			err := runFunc()
			if err != nil {
				return err
			}

			// if we define the runSleepDuration to be above every millisecond, then we are allowed to sleep
			if s.GetRunSleepDuration() > 1*time.Millisecond {
				time.Sleep(s.GetRunSleepDuration())
			}
		}
	} else {
		// do nothing
	}

	// Custom stop func if provided
	if stopFunc != nil {
		err := stopFunc()
		if err != nil {
			return err
		}
	} else {
		// do nothing
	}

	if s.logLevel >= LOG_LEVEL_INFO {
		time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
		log.Printf("%s stopped gracefully\n", s.Name)
	}

	s.canRestart = true
	s.isInitialized = false
	return nil
}

// Restart restarts the service
func (s *Service) Restart() error {
	if !s.isInterrupted {
		if s.logLevel >= LOG_LEVEL_INFO {
			time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
			log.Println("Calling for restart of service: " + s.Name)
			err := s.Stop() // ignore stop err
			if err != nil {
				log.Println(err)
			}

		}

		for {
			if s.logLevel >= LOG_LEVEL_INFO {
				time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
			}
			if s.isInterrupted {
				break
			}

			if s.canRestart {
				s.canRestart = false
				s.isRunning = true
				err := s.Start(s.customFunctions[0], s.customFunctions[1], s.customFunctions[2], s.customFunctions[3])
				return err
			}
		}
	}
	return errors.New("restart failed (start was interrupted)")
}

// Stop stops the service by setting isRunning to false.
func (s *Service) Stop() error {
	if s.isRunning {
		if s.logLevel >= LOG_LEVEL_INFO {
			time.Sleep(20 * time.Millisecond) // to prevent log package from race condition logging most of the time
			log.Println("Stopping service: " + s.Name)
		}
		s.isRunning = false
		return nil
	}

	if s.logLevel >= LOG_LEVEL_WARN {
		return errors.New("Service was not running: " + s.Name)
	}
	return nil
}

// ForceShutdown forcefully stops both the service and the whole program and logs an error. (note: forcing shutdown is not graceful)
func (s *Service) ForceShutdown() error {
	err := s.Stop()
	if err != nil {
		return err
	}
	if s.logLevel >= LOG_LEVEL_ERROR {
		log.Println("(Timeout) forced shutdown of program with all its running services")
	}
	os.Exit(1)
	return nil
}

// listenForInterrupt listens for interrupt signals and triggers shutdown.
func (s *Service) listenForInterrupt(forceShutdown func() error) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal // Block until a signal is received
	s.isInterrupted = true
	// printing interrupt signal warning regardless of s.PrintLog
	if s.logLevel >= LOG_LEVEL_WARN {
		log.Printf("%s received interrupt signal, initiating graceful shutdown (timeout: %v)\n", s.Name, s.gracefulShutdownTime)
	}

	signal.Stop(osSignal)
	close(osSignal)

	err := s.Stop() // Stop the service
	if err != nil {
		log.Println(err)
	}

	// Schedule a forced shutdown if the graceful shutdown time elapses
	go func() {
		<-time.After(s.gracefulShutdownTime)

		// Custom forceShutdown func if provided
		if forceShutdown != nil {
			_ = forceShutdown() // ignore err
		} else {
			// if forceShutdownFunc is not implemented by the user, then run ForceShutdown (exits program with log)
			_ = s.ForceShutdown() // ignore err
		}
	}()
}
