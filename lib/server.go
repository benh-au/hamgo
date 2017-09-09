package lib

import (
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
)

// TCPServer provides a tcp listening server.
type TCPServer struct {
	Port          uint
	listener      net.Listener
	close         bool
	NewConnection chan *Connection
}

func (t *TCPServer) worker() {
	for {
		logrus.Debug("TCPServer: accepting connections")

		conn, err := t.listener.Accept()
		if err != nil {
			logrus.WithError(err).Warn("TCPServer: failed to accept client")
			continue
		}

		logrus.Infof("TCPServer: new connection: %s", conn.RemoteAddr().String())

		// create a new connection
		c := Connection{
			Connection: conn,
			Send:       make(chan *Message),
			close:      make(chan interface{}),
			Received:   make(chan []byte),
		}

		// signal that a new connection is established
		t.NewConnection <- &c

		// start a goroutine for the worker
		go c.connectionWorker()
	}
}

// Start the tcp server and listen for incoming connections.
func (t *TCPServer) Start() error {
	logrus.Debugf("TCPServer: listening on port %d", t.Port)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", t.Port))

	if err != nil {
		logrus.WithError(err).Warn("TCPServer: Failed to listen on tcp port")
		return err
	}

	t.NewConnection = make(chan *Connection)
	t.listener = conn

	// start the listener accept worker
	go t.worker()

	return nil
}

// Stop the tcp server.
func (t *TCPServer) Stop() {
	logrus.Debug("TCPServer: stopping")

	t.close = true

	err := t.listener.Close()
	if err != nil {
		logrus.WithError(err).Warn("TCPServer: failed to stop listener")
	}
}
