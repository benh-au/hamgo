package node

import (
	"errors"
	"sync"

	"github.com/donothingloop/hamgo/lib"
	"github.com/donothingloop/hamgo/parameters"
	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
)

// Node is a node in the gossip protocol.
type Node struct {
	server      lib.TCPServer
	settings    parameters.Settings
	station     parameters.Station
	peers       []*Peer
	logic       *Logic
	close       chan interface{}
	cbs         []*MessageCallback
	cbsPeerConn []*PeerConnCallback
	Cache       []*protocol.Message
	cacheLock   sync.Mutex
	Local       protocol.Contact
}

// MessageCallback is a callback that is called when a message was received.
type MessageCallback struct {
	Cb func(*protocol.Message, *Peer)
}

// PeerConnCallback is called when a peer reconnects.
type PeerConnCallback struct {
	PeerConnected func(*Peer)
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

// triggerPeerConnected is used to trigger all peer connected callbacks.
func (n *Node) triggerPeerConnected(peer *Peer) {
	for _, cb := range n.cbsPeerConn {
		cb.PeerConnected(peer)
	}
}

// AddPeerConnCallback adds a peer connected callback.
func (n *Node) AddPeerConnCallback(cb *PeerConnCallback) {
	n.cbsPeerConn = append(n.cbsPeerConn, cb)
}

// RemovePeerConnCallback removes a previously added peer connected callback.
func (n *Node) RemovePeerConnCallback(cb *PeerConnCallback) {
	var cbs []*PeerConnCallback

	for _, v := range n.cbsPeerConn {
		if v != cb {
			cbs = append(cbs, v)
		}
	}

	n.cbsPeerConn = cbs
}

// existsInCache checks if a given message exists in the cache.
func (n *Node) existsInCache(msg *protocol.Message) bool {
	logrus.WithField("msg", msg).Debug("Node: check if message is in cache")

	for _, v := range n.Cache {
		if (v.SeqCounter == msg.SeqCounter) && v.Source.Compare(&msg.Source) {
			logrus.WithField("msg", msg).Debug("Node: found message in cache")
			return true
		}
	}

	logrus.Debug("Node: did not find message in cache")
	return false
}

// AddToCache adds a remote message to the cache.
func (n *Node) AddToCache(msg *protocol.Message) {
	// append local node to path
	msg.Path += ";" + n.station.Callsign
	msg.PathLength = uint16(len(msg.Path))

	// decrease TTL
	if msg.TTL != 0 {
		msg.TTL--
	}

	n.pushToCache(msg)
}

// pushToCache pushes a message to cache if it does not exist in it
func (n *Node) pushToCache(msg *protocol.Message) bool {
	n.cacheLock.Lock()
	defer n.cacheLock.Unlock()

	if (msg.Flags & protocol.FlagNoCache) != 0 {
		logrus.Debug("node: not caching message with no-cache flag")
		return true
	}

	if n.existsInCache(msg) {
		logrus.Debug("Node: message already cached, ignoring")
		return false
	}

	// remove first cache entry, if cache is full
	if uint(len(n.Cache)) == n.settings.LogicSettings.CacheSize && len(n.Cache) > 1 {
		n.Cache = n.Cache[1:]
	}

	n.Cache = append(n.Cache, msg)
	return true
}

// handleCallbacks calls all registered callbacks that hook the received messages.
func (n *Node) handleCallbacks(msg *protocol.Message, src *Peer) {
	// call some fixed handlers
	n.consoleHandler(msg)

	for _, v := range n.cbs {
		v.Cb(msg, src)
	}
}

// SpreadMessage spreads a message by gossip.
func (n *Node) SpreadMessage(msg *protocol.Message) error {
	if n.settings.LogicSettings.ReadOnly {
		logrus.Warn("Node: node is read-only, ignoring spread message")
		return errors.New("read-only node")
	}

	// message already cached, ignoring
	if !n.pushToCache(msg) {
		return nil
	}

	go n.handleCallbacks(msg, nil)

	return n.logic.SpreadMessage(msg)
}

// handleMessage handles a message from a peer.
func (n *Node) handleMessage(msg []byte, src *Peer) {
	pmsg, _ := protocol.ParseMessage(msg)

	// message already cached, ignoring
	if !n.pushToCache(&pmsg) {
		return
	}

	go n.handleCallbacks(&pmsg, src)

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

		case <-p.reconnected:
			logrus.Debug("Node: peer connected")
			n.triggerPeerConnected(p)
			break

		case msg := <-p.Received:
			logrus.Debug("Node: message received")
			n.handleMessage(msg, p)
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

		// start the peer worker
		go n.peerWorker(p)

		// start the peer
		p.Start()
		p.Reconnect()
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

		// start the peer worker
		go n.peerWorker(np)

		np.Start()

		// set the connection and start the read
		np.SetConnection(conn)

		n.logic.peers = append(n.logic.peers, np)
		p = np
	} else {
		logrus.Debug("Node: setting connections")
		p.SetConnection(conn)
	}

	go n.triggerPeerConnected(p)
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
		station:  station,
		server: lib.TCPServer{
			Port: settings.Port,
		},
		logic: &Logic{
			settings:        settings.LogicSettings,
			settingsStation: station,
		},
		Local: protocol.Contact{
			Type:           protocol.ContactTypeFixed,
			CallsignLength: uint8(len(station.Callsign)),
			Callsign:       []byte(station.Callsign),
			NumberIPs:      0,
			IPs:            []protocol.ContactIP{},
		},
	}

	n.logic.Local = n.Local

	return n, nil
}

// Init the node.
func (n *Node) Init() error {
	// create the peer instances
	n.createPeers()

	// set the peers for the logic
	n.logic.peers = n.peers

	logrus.Debug("Node: starting server")
	err := n.server.Start()
	if err != nil {
		logrus.WithError(err).Warn("Node: failed to start TCP server")
		return err
	}

	// start the connection worker for the server
	go n.connectionWorker()

	return nil
}
