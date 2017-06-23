package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type options struct {
	Verbose []bool `short:"v" long:"verbose" description:"increases logging verbosity"`

	Args struct {
		CMD string `description:"path to server command to run"`
	} `positional-args:"yes" required:"yes"`
}

var opts options

func main() {
	parseOptions()
}

func parseOptions() {
	if _, err := flags.Parse(&opts); err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}

		panic(err)
	}
}
