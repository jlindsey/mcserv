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

type service interface {
	fmt.Stringer
	Start()
	Stop()
	Done() chan error
}

type serviceRegistry struct {
	services []service
	finished []bool
	errors   []error
	sigChan  chan os.Signal
	mux      *sync.Mutex
}

func (r *serviceRegistry) String() string {
	return fmt.Sprintf("serviceRegistry{registered: %d}", len(r.services))
}

func newServiceRegistry() *serviceRegistry {
	r := serviceRegistry{
		services: make([]service, 0),
		finished: make([]bool, 0),
		errors:   make([]error, 0),
		mux:      &sync.Mutex{},
	}

	return &r
}

func (r *serviceRegistry) add(f service) {
	log.WithFields(log.Fields{
		"registry": r,
		"service":  f,
	}).Debug("Add service to registry")
	r.services = append(r.services, f)
}

func (r *serviceRegistry) start() {
	logger := log.WithField("registry", r)
	logger.Debug("Starting services")

	for i := range r.services {
		service := r.services[i]
		logger.WithField("service", service).Debug("Starting service")
		go service.Start()
	}
}

func (r *serviceRegistry) stop() {
	logger := log.WithField("registry", r)
	logger.Debug("Stopping services")

	for i := range r.services {
		service := r.services[i]
		logger.WithField("service", service).Debug("Stopping service")
		service.Stop()
	}
}

func (r *serviceRegistry) setupSignalHandler() {
	logger := log.WithField("registry", r)
	logger.Debug("Setting up signal handler")

	r.sigChan = make(chan os.Signal)
	signal.Notify(r.sigChan, os.Interrupt)

	go func() {
		<-r.sigChan
		logger.Debug("Caught SIGINT")
		r.stop()
	}()
}

func (r *serviceRegistry) wait() error {
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

func (r *serviceRegistry) waitForFinisher(f service) {
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
	}
	r.finished = append(r.finished, true)
	r.mux.Unlock()

	logger.Debug("Service finished")
}
