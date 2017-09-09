package node

import (
	"sync"

	"github.com/donothingloop/hamgo/lib"
	"github.com/donothingloop/hamgo/parameters"
	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
)

// Node is a node in the gossip protocol.
type Node struct {
	server    lib.TCPServer
	settings  parameters.Settings
	peers     []*Peer
	logic     *Logic
	close     chan interface{}
	cbs       []*MessageCallback
	Cache     []*protocol.Message
	cacheLock sync.Mutex
}

// MessageCallback is a callback that is called when a message was received.
type MessageCallback struct {
	Cb func(*protocol.Message)
}

// AddCallback adds a callback for received messages.
func (n *Node) AddCallback(cb *MessageCallback) {
	n.cbs = append(n.cbs, cb)
}

// RemoveCallback removes the callback from the buffer.
func (n *Node) RemoveCallback(cb *MessageCallback) {
	var cbs []*MessageCallback

	for _, v := range n.cbs {
		if v != cb {
			cbs = append(cbs, v)
		}
	}

	n.cbs = cbs
}

// existsInCache checks if a given message exists in the cache.
func (n *Node) existsInCache(msg *protocol.Message) bool {
	for _, v := range n.Cache {
		if (v.SeqCounter == msg.SeqCounter) && v.Source.Compare(&msg.Source) {
			return true
		}
	}

	return false
}

// pushToCache pushes a message to cache if it does not exist in it
func (n *Node) pushToCache(msg *protocol.Message) {
	n.cacheLock.Lock()

	if n.existsInCache(msg) {
		logrus.Debug("Node: message already cached, ignoring")
		n.cacheLock.Unlock()
		return
	}

	// remove first cache entry, if cache is full
	if uint(len(n.Cache)) == n.settings.LogicSettings.CacheSize && len(n.Cache) > 1 {
		n.Cache = n.Cache[1:]
	}

	n.Cache = append(n.Cache, msg)
	n.cacheLock.Unlock()
}

// handleCallbacks calls all registered callbacks that hook the received messages.
func (n *Node) handleCallbacks(msg *protocol.Message) {
	n.consoleHandler(msg)

	for _, v := range n.cbs {
		v.Cb(msg)
	}
}

// SpreadMessage spreads a message by gossip.
func (n *Node) SpreadMessage(msg *protocol.Message) error {
	n.pushToCache(msg)
	go n.handleCallbacks(msg)

	return n.logic.SpreadMessage(msg)
}

// handleMessage handles a message from a peer.
func (n *Node) handleMessage(msg []byte) {
	pmsg := protocol.ParseMessage(msg)
	n.pushToCache(&pmsg)
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

	return n, nil
}
