package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/sirupsen/logrus"
)

// Server is our RPC server
type Server struct {
	socket            *net.UnixListener
	rpcServer         *rpc.Server
	rpcServerDoneChan chan error
	shouldQuit        bool

	// SocketPath is the path of the RPC control socket
	SocketPath string
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{SocketPath: %v}", s.SocketPath)
}

// NewServer instantiates a new RPC server
func NewServer(socketPath string) *Server {
	s := Server{
		rpcServer:         rpc.NewServer(),
		rpcServerDoneChan: make(chan error),
		shouldQuit:        false,
		SocketPath:        socketPath,
	}

	return &s
}

// Register adds a new RPC service to the server
func (s *Server) Register(i interface{}) error {
	log.WithFields(log.Fields{
		"service": i,
		"server":  s,
	}).Debug("Registering RPC")
	return s.rpcServer.Register(i)
}

// Done returns a channel that will block until the server is exited.
// Any errors from this channel should be handled by the calling routine.
func (s *Server) Done() chan error {
	return s.rpcServerDoneChan
}

// Stop the RPC server
func (s *Server) Stop() {
	s.shouldQuit = true
	if err := s.socket.SetDeadline(time.Now().Add(-10 * time.Second)); err != nil {
		log.Panic(err)
	}
}

func (s *Server) finish() {
	s.rpcServerDoneChan <- nil
}

// Start the RPC server. Should be called in a goroutine.
func (s *Server) Start() {
	logger := log.WithField("server", s)
	logger.Info("Starting RPC Server")

	err := s.openSocket()
	if err != nil {
		s.rpcServerDoneChan <- err
		return
	}
	defer func() {
		if err = s.socket.Close(); err != nil {
			logger.Panic(err)
		}
	}()

	for {
		if s.shouldQuit {
			break
		}

		conn, err := s.socket.Accept()
		if err != nil {
			if err, ok := err.(*net.OpError); ok && err.Timeout() {
				logger.Debug("Timeout")
				continue
			}

			s.rpcServerDoneChan <- err
			return
		}

		logger.WithField("local", conn.LocalAddr()).Debug("New Connection")

		go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

	logger.Info("Stopped RPC Server")
	s.finish()
}

func (s *Server) openSocket() error {
	l, err := net.Listen("unix", s.SocketPath)
	if err != nil {
		return err
	}

	cast, ok := l.(*net.UnixListener)
	if !ok {
		return fmt.Errorf("Unable to cast socket to correct type")
	}

	s.socket = cast
	return nil
}
