package node

import (
	"sync"
	"time"

	"github.com/donothingloop/hamgo/lib"
	"github.com/donothingloop/hamgo/parameters"

	"github.com/Sirupsen/logrus"
)

// Peer stores a peer of the gossip protocol.
type Peer struct {
	Settings         parameters.Settings
	connection       *lib.Connection
	queue            [][]byte
	checkMessages    chan interface{}
	close            chan interface{}
	connActiveClose  chan interface{}
	writeLock        sync.Mutex
	sendTries        uint
	connectionActive bool
	Received         chan []byte
	client           *lib.TCPClient
	fromServer       bool
	checkPending     bool
}

// NewPeer creates a new peer.
func NewPeer(host string, port uint, settings parameters.Settings) *Peer {
	return &Peer{
		Settings:        settings,
		checkMessages:   make(chan interface{}, 10),
		connActiveClose: make(chan interface{}),
		sendTries:       0,
		Received:        make(chan []byte),
		close:           make(chan interface{}),
		client: &lib.TCPClient{
			Host: host,
			Port: port,
		},
	}
}

// Close the peer.
func (p *Peer) Close() {
	close(p.close)
}

func (p *Peer) writeCallback(conn *lib.Connection, err error) {
	logrus.Debug("Peer: write callback")

	// remove the message from the queue if the send is successful
	if err == nil {
		logrus.WithField("queuelen", len(p.queue)).Debug("Peer: queuelen")

		if len(p.queue) == 0 || len(p.queue) == 1 {
			// clear queue
			p.queue = [][]byte{}
		} else {
			p.queue = p.queue[1:]
		}

		p.sendTries = 0

		logrus.WithField("queuelen", len(p.queue)).Debug("Peer: queuelen after")
		logrus.Debug("Peer: message sent successfully, removed from queue")
	} else {
		p.sendTries++
		logrus.Debug("Peer: message not sent successfully, retrying")

		if p.sendTries > p.Settings.Retries {
			logrus.Debug("Peer: maximum number of retrys reached, closing connection")

			// terminate the connection if it is faulty
			p.connection.Close()
			p.connectionActive = false

			// close the connection read worker
			close(p.connActiveClose)
		}

		p.sendTries = 0
	}

	p.writeLock.Unlock()
}

// Start the peer worker.
func (p *Peer) Start() {
	logrus.Debug("Peer: starting workers")

	go p.worker()

	if !p.fromServer {
		go p.reconnectWorker()
	}
}

// readWorker reads from the stream.
func (p *Peer) readWorker() {
	logrus.Debug("Peer: readWorker: active")

	for {
		select {
		case <-p.connActiveClose:
			logrus.Debug("Peer: connActiveClose signalled")
			return

		case msg := <-p.connection.Received:
			logrus.WithField("msg", msg).Debug("Peer: message received")
			p.Received <- msg
			break
		}
	}
}

// reconnect the connection.
func (p *Peer) reconnect() {
	logrus.Info("Peer: reconnecting")

	conn, err := p.client.Start()
	if err != nil {
		logrus.WithError(err).Warn("Peer: failed to reconnect")
		return
	}

	logrus.Debug("Peer: reconnected")

	p.connection = conn
	p.connectionActive = true

	p.connActiveClose = make(chan interface{})

	if !p.checkPending {
		p.checkPending = true

		// signal to check new messages
		p.checkMessages <- nil
	}

	// start the read worker
	go p.readWorker()
}

// SetConnection sets a new connection and initializes the workers.
func (p *Peer) SetConnection(conn *lib.Connection) {
	logrus.Debug("Peer: setting new connection")

	if p.connectionActive {
		p.connectionActive = false
		close(p.connActiveClose)
	}

	p.connection = conn
	p.connectionActive = true
	p.connActiveClose = make(chan interface{})

	if !p.checkPending {
		p.checkPending = true
		p.checkMessages <- nil
	}

	go p.readWorker()
}

// reconnectWorker handles the reconnecting of the connection.
func (p *Peer) reconnectWorker() {
	recon := time.Tick(time.Duration(p.Settings.ReconnectTimeout) * time.Second)

	for {
		select {
		case <-p.close:
			logrus.Debug("Peer: reconnect worker closed")
			return

		case <-recon:
			if !p.connectionActive {
				logrus.Debug("Peer: reconnect timer tick")
				p.reconnect()
			}
			break
		}
	}
}

// worker is the worker that sends messages to the peer.
func (p *Peer) worker() {
	for {
		for {
			p.writeLock.Lock()

			// if the connection is not active anymore, wait for a checkMessages signal
			if !p.connectionActive {
				p.writeLock.Unlock()
				logrus.Debug("Peer: worker: connection not active")
				break
			}

			// if the queue is empty, wait for the signal
			if len(p.queue) == 0 {
				p.writeLock.Unlock()
				logrus.Debug("Peer: worker: queue is empty")
				break
			}

			msg := &lib.Message{
				Data:     p.queue[0],
				Callback: p.writeCallback,
			}

			logrus.Debug("Peer: worker: acquiring lock")

			logrus.Debug("Peer: worker: sending message")
			p.connection.Send <- msg
		}

		logrus.Debug("Peer: worker: waiting for signal")

		// wait for a signal
		<-p.checkMessages

		p.checkPending = false
	}
}

// QueueMessage queues a message to be sent to the peer.
func (p *Peer) QueueMessage(msg []byte) {
	if uint(len(p.queue)) >= p.Settings.PeerQueueSize {
		logrus.Warn("Peer: peer queue full, dropping oldest message")
		p.queue = p.queue[1:]
	}

	p.queue = append(p.queue, msg)

	logrus.WithField("msg", msg).Debug("Peer: queueed peer message")

	if !p.checkPending {
		p.checkPending = true

		// send the check signal
		p.checkMessages <- nil
	}
}
