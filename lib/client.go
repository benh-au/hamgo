package lib

import (
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
)

// TCPClient is a tcp client for the communication between peers.
type TCPClient struct {
	Host string
	Port uint
}

// Start the connection.
func (c *TCPClient) Start() (*Connection, error) {
	logrus.Infof("TCPClient: connecting to %s:%d", c.Host, c.Port)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		logrus.WithError(err).Warn("TCPClient: failed to connect to host")
		return nil, err
	}

	logrus.Info("TCPClient: connection established")

	co := Connection{
		Connection: conn,
		Send:       make(chan *Message),
		close:      make(chan interface{}),
		Received:   make(chan []byte),
	}

	// start a goroutine for the worker
	go co.connectionWorker()

	return &co, nil
}
