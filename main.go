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
