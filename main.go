package main

import (
	"github.com/jlindsey/mcserv/rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	opts, err := parseOptions()
	if err != nil {
		log.Fatal(err)
	}

	registry := newServiceRegistry()

	server := rpc.NewServer(opts.SocketPath)
	err = server.Register(new(Service))
	if err != nil {
		log.Panic(err)
	}
	registry.add(server)

	registry.setupSignalHandler()
	registry.start()
	err = registry.wait()

	if err != nil {
		log.Error(err)
	}
}
