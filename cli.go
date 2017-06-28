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

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type options struct {
	Verbose    bool   `short:"v" long:"verbose" description:"increases logging verbosity"`
	SocketPath string `short:"s" long:"socket" description:"path to the control socket" default:"/var/run/mcserv.sock"`

	Args struct {
		CMD string `description:"path to server command to run"`
	} `positional-args:"yes" required:"yes"`
}

func (o options) String() string {
	return fmt.Sprintf(
		"options{Verbose: %v, SocketPath: %s, CMD: %s}",
		o.Verbose,
		o.SocketPath,
		o.Args.CMD,
	)
}

func parseOptions() (options, error) {
	var opts options

	if _, err := flags.Parse(&opts); err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}

		return opts, err
	}

	log.SetLevel(log.InfoLevel)
	if opts.Verbose {
		log.SetLevel(log.DebugLevel)
		log.WithFields(log.Fields{
			"opts": opts,
		}).Debug("Debug")
	}

	return opts, nil
}
