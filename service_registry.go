/*
mcserv
Copyright (C) 2017 Joshua Lindsey

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Service is an interface that defines the minimum methods required to be
// added to the ServiceRegistry.
type Service interface {
	// Start should start the underlying Service. It does not
	// necessarily have to block, though it is called in a goroutine
	// by the Registry.
	Start()

	// Stop should stop the underlying service. It should not block,
	// but should behave in such a way that Done() reports correctly.
	Stop()

	// Done returns an error channel that remains open as long as the
	// Service is running and is pushed either `nil` for a successful
	// completion, or an `error`.
	Done() chan error

	fmt.Stringer
}

// ServiceRegistry is a container for Service pointers, managing
// their lifecycles in a centralized interface.
type ServiceRegistry struct {
	services []Service
	finished []bool
	errors   []error
	sigChan  chan os.Signal
	mux      *sync.Mutex
}

func (r *ServiceRegistry) String() string {
	return fmt.Sprintf("serviceRegistry{registered: %d}", len(r.services))
}

// NewServiceRegistry returns a pointer to a new ServiceRegistry
func NewServiceRegistry() *ServiceRegistry {
	r := ServiceRegistry{
		services: make([]Service, 0),
		finished: make([]bool, 0),
		errors:   make([]error, 0),
		mux:      &sync.Mutex{},
	}

	return &r
}

// Add registers a new Service to the Registry.
func (r *ServiceRegistry) Add(f Service) {
	log.WithFields(log.Fields{
		"registry": r,
		"service":  f,
	}).Debug("Add service to registry")
	r.services = append(r.services, f)
}

// Start iterates through the registered Services and calls their Start()
// methods in a goroutine.
func (r *ServiceRegistry) Start() {
	logger := log.WithField("registry", r)
	logger.Debug("Starting services")

	for i := range r.services {
		service := r.services[i]
		logger.WithField("service", service).Debug("Starting service")
		go service.Start()
	}
}

// Stop iterates through the registered Services and calls their Stop() methods.
func (r *ServiceRegistry) Stop() {
	logger := log.WithField("registry", r)
	logger.Debug("Stopping services")

	for i := range r.services {
		service := r.services[i]
		logger.WithField("service", service).Debug("Stopping service")
		service.Stop()
	}
}

// SetupSignalHandler registers a handler for SIGINT and a listener
// on the resulting channel to cleanly stop registered Services.
func (r *ServiceRegistry) SetupSignalHandler() {
	logger := log.WithField("registry", r)
	logger.Debug("Setting up signal handler")

	r.sigChan = make(chan os.Signal)
	signal.Notify(r.sigChan, os.Interrupt)

	go func() {
		<-r.sigChan
		logger.Debug("Caught SIGINT")
		r.Stop()
	}()
}

// Wait will wait until the underlying Services' Done() methods
// indicate exit. If any of them return an error, this method
// will pass them through.
func (r *ServiceRegistry) Wait() error {
	for i := range r.services {
		go r.waitForFinisher(r.services[i])
	}

	for {
		if len(r.services) == len(r.finished) {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	if len(r.errors) > 0 {
		return fmt.Errorf("Finished with errors: %v", r.errors)
	}

	return nil
}

func (r *ServiceRegistry) waitForFinisher(f Service) {
	logger := log.WithFields(log.Fields{
		"registry": r,
		"service":  f,
	})
	logger.Debug("Waiting for service")

	err := <-f.Done()
	r.mux.Lock()
	if err != nil {
		logger.Error(err)
		r.errors = append(r.errors, err)

		logger.Error("Error encountered. Stopping all services")
		r.Stop()
	}
	r.finished = append(r.finished, true)
	r.mux.Unlock()

	logger.Debug("Service finished")
}
