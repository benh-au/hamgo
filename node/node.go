package node

import (
	"github.com/donothingloop/hamgo/lib"
	"github.com/donothingloop/hamgo/parameters"
	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
)

// Node is a node in the gossip protocol.
type Node struct {
	server   lib.TCPServer
	settings parameters.Settings
	peers    []*Peer
	logic    *Logic
	close    chan interface{}
	cbs      []MessageCallback
}

// MessageCallback is a callback that is called when a message was received.
type MessageCallback func(*protocol.Message)

// AddCallback adds a callback for received messages.
func (n *Node) AddCallback(cb MessageCallback) {
	n.cbs = append(n.cbs, cb)
}

// handleCallbacks calls all registered callbacks that hook the received messages.
func (n *Node) handleCallbacks(msg *protocol.Message) {
	for _, v := range n.cbs {
		v(msg)
	}
}

// SpreadMessage spreads a message by gossip.
func (n *Node) SpreadMessage(msg *protocol.Message) error {
	return n.logic.SpreadMessage(msg)
}

// handleMessage handles a message from a peer.
func (n *Node) handleMessage(msg []byte) {
	pmsg := protocol.ParseMessage(msg)
	go n.handleCallbacks(&pmsg)

	n.logic.HandleMessage(msg)
}

// Close the node.
func (n *Node) Close() {
	logrus.Debug("Node: closing")

	// close all peers
	for _, p := range n.peers {
		p.Close()
	}

	n.server.Stop()
}

func (n *Node) peerWorker(p *Peer) {
	logrus.Debug("Node: peerWorker: active")

	for {
		select {
		case <-n.close:
			logrus.Debug("Node: peerWorker: closing globally")
			return

		case <-p.close:
			logrus.Debug("Node: peerWorker: closing peer")
			return

		case msg := <-p.Received:
			logrus.Debug("Node: message received")
			n.handleMessage(msg)
			break
		}
	}
}

// createPeers creates the peers for the node.
func (n *Node) createPeers() {
	logrus.Debug("Node: creating peers")

	for _, v := range n.settings.Peers {
		p := NewPeer(v.Host, v.Port, n.settings)
		n.peers = append(n.peers, p)

		logrus.Debug("Node: starting peer")

		// start the peer
		p.Start()

		// start the peer worker
		go n.peerWorker(p)
	}
}

func (n *Node) findPeerByConn(conn *lib.Connection) *Peer {
	for _, p := range n.logic.peers {
		if p.connection == nil || conn.Connection == nil {
			continue
		}

		if p.connection.Connection.RemoteAddr() == conn.Connection.RemoteAddr() {
			return p
		}
	}

	return nil
}

func (n *Node) handleConnection(conn *lib.Connection) {
	logrus.Debug("Node: handling connection")
	logrus.Debug("Node: searching for existing connections")
	p := n.findPeerByConn(conn)

	if p == nil {
		logrus.Info("Node: creating new peer")
		np := NewPeer(conn.Connection.RemoteAddr().String(), n.settings.Port, n.settings)
		np.fromServer = true
		np.Start()

		// set the connection and start the read
		np.SetConnection(conn)

		n.logic.peers = append(n.logic.peers, np)

		// start the peer worker
		go n.peerWorker(np)
	} else {
		logrus.Debug("Node: setting connections")
		p.SetConnection(conn)
	}
}

// connectionWorker waits for new connections and adds them as peers.
func (n *Node) connectionWorker() {
	logrus.Debug("Node: connectionWorker: started")

	for {
		select {
		case <-n.close:
			logrus.Debug("Node: connectionWorker: closing")
			break

		case conn := <-n.server.NewConnection:
			n.handleConnection(conn)
			break
		}
	}
}

// NewNode creates a new node.
func NewNode(settings parameters.Settings, station parameters.Station) (*Node, error) {
	logrus.Debug("Node: creating new instance")

	n := &Node{
		settings: settings,
		server: lib.TCPServer{
			Port: settings.Port,
		},
		logic: &Logic{
			settings:        settings.LogicSettings,
			settingsStation: station,
		},
	}

	// create the peer instances
	n.createPeers()

	// set the peers for the logic
	n.logic.peers = n.peers

	logrus.Debug("Node: starting server")
	err := n.server.Start()
	if err != nil {
		logrus.WithError(err).Warn("Node: failed to start TCP server")
		return nil, err
	}

	// start the connection worker for the server
	go n.connectionWorker()

	// register handlers
	n.AddCallback(n.consoleHandler)

	return n, nil
}
