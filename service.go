package ggservice

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// IService defines the interface for managing a service with start, stop, and force shutdown capabilities.
type IService interface {
	Start(startFunc func() error, runFunc func() error, forceExitFunc func() error) error
	Stop()
	ForceShutdown()
}

// Service represents a service that can be started, stopped, and forcefully shutdown with graceful handling.
type Service struct {
	Name                 string        // Name of the service
	GracefulShutdownTime time.Duration // Timeout duration for graceful shutdown
	isRunning            bool          // Flag indicating whether the service is running
}

// New creates a new instance of Service with the given name and graceful shutdown timeout.
func New(service *Service) IService {
	return &Service{
		Name:                 service.Name,
		GracefulShutdownTime: service.GracefulShutdownTime,
		isRunning:            true,
	}
}

// NewService creates a new instance of Service with the given name and graceful shutdown timeout.
func NewService(name string, gracefulShutdownTime time.Duration) IService {
	return New(&Service{
		Name:                 name,
		GracefulShutdownTime: gracefulShutdownTime,
	})
}

// Start starts the service with custom start, run, and stop functions.
func (s *Service) Start(startFunc func() error, runFunc func() error, forceExitFunc func() error) error {
	go s.listenForInterrupt(forceExitFunc) // Listen for interrupt signals

	// Execute custom start function if provided
	fmt.Printf("Starting service: %s\n", s.Name)
	if startFunc != nil {
		err := startFunc()
		if err != nil {
			return err
		}
	}

	// Execute custom run function in a loop as long as the service is running
	if runFunc != nil {
		for s.isRunning {
			err := runFunc()
			if err != nil {
				return err
			}
		}
		fmt.Printf("%s stopped gracefully\n", s.Name)
	}

	return nil
}

// Stop stops the service by setting isRunning to false.
func (s *Service) Stop() {
	s.isRunning = false
}

// ForceShutdown forcefully stops the service and logs a fatal error. (note: forcing shutdown is not graceful)
func (s *Service) ForceShutdown() {
	s.Stop()
	log.Fatal("Forced shutdown")
}

// listenForInterrupt listens for interrupt signals and triggers shutdown.
func (s *Service) listenForInterrupt(forceExit func() error) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal // Block until a signal is received
	fmt.Printf("%s received interrupt signal, initiating graceful shutdown (timeout: %v)\n", s.Name, s.GracefulShutdownTime)
	signal.Stop(osSignal)
	close(osSignal)

	s.Stop() // Stop the service

	// Schedule a forced shutdown if the graceful shutdown time elapses
	go func() {
		<-time.After(s.GracefulShutdownTime)

		// Execute custom ForceExit func
		if forceExit != nil {
			_ = forceExit() // ignoring error in interrupt functionality (the error is meant for manual use only)
		} else {
			// if forceExitFunc is not implemented by the user, then run ForceShutdown (exits program with log)
			s.ForceShutdown()
		}
	}()
}
