package main

import (
	"github.com/jlindsey/mcserv/shared"
	log "github.com/sirupsen/logrus"
)

// Service is an instantiation of the shared.Service interface
type Service int

// Ping is a simple command that checks connectivity
func (s *Service) Ping(args *shared.PingArgs, ret *shared.Pong) error {
	log.WithFields(log.Fields{
		"func": "Ping",
		"args": args,
	}).Debug("Got RPC call")

	ret.OK = true

	log.WithFields(log.Fields{
		"func": "Ping",
		"args": args,
		"ret":  ret,
	}).Info("RPC call")

	return nil
}
