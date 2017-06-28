package shared

import "fmt"

type (
	// PingArgs is an empty struct
	PingArgs struct{}

	// Pong is the return of the Ping command
	Pong struct {
		// OK will be true
		OK bool
	}

	// RPC is the exported interface
	RPC interface {
		fmt.Stringer
		Ping(*PingArgs, *Pong) error
	}
)

func (a *PingArgs) String() string {
	return "PingArgs{}"
}

func (r *Pong) String() string {
	return fmt.Sprintf("Pong{OK: %v}", r.OK)
}
