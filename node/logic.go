package node

import (
	"errors"

	"github.com/donothingloop/hamgo/parameters"
	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
)

type cacheEntry struct {
	SeqCounter uint16
	Source     protocol.Contact
}

// Logic handles the forwarding of the nodes.
type Logic struct {
	settings        parameters.LogicSettings
	settingsStation parameters.Station
	cache           []*cacheEntry
	peers           []*Peer
}

func (n *Logic) isMessageCached(msg *protocol.Message) bool {
	logrus.Debug("Logic: check if message is cached")

	for _, c := range n.cache {
		if c.SeqCounter == msg.SeqCounter && c.Source.Compare(&msg.Source) {
			logrus.Debug("Logic: message is cached")
			return true
		}
	}

	logrus.Debug("Logic: message is not cached")
	return false
}

func (n *Logic) cacheMessage(msg *protocol.Message) {
	if uint(len(n.cache)) >= n.settings.CacheSize {
		n.cache = n.cache[1:]
		logrus.Debug("Logic: cache clean")
	}

	logrus.Debug("Logic: caching message")

	n.cache = append(n.cache, &cacheEntry{
		SeqCounter: msg.SeqCounter,
		Source:     msg.Source,
	})
}

// SpreadMessage caches a new message and spreads it afterwards.
func (n *Logic) SpreadMessage(msg *protocol.Message) error {
	if n.isMessageCached(msg) {
		return errors.New("message already cached")
	}

	logrus.Debug("Logic: spreading new message")

	// append local node to path
	msg.Path += ";" + n.settingsStation.Callsign
	msg.PathLength = uint16(len(msg.Path))

	if !n.isMessageCached(msg) {
		n.cacheMessage(msg)
		n.spreadCachedMessage(msg)
	} else {
		logrus.Info("Logic: message to be spread is already cached, ignoring")
	}

	return nil
}

// spreadCachedMessage spreads a message using the gossip protocol.
func (n *Logic) spreadCachedMessage(msg *protocol.Message) {
	buf := msg.Bytes()

	logrus.WithField("msg", msg).Debug("Logic: spreading cached message")

	for _, p := range n.peers {
		// enqueue the message for the peer to be sent
		p.QueueMessage(buf)
	}
}

// HandleMessage handles an incoming message from a peer.
func (n *Logic) HandleMessage(msg []byte) {
	logrus.Debug("Logic: parsing incoming message")

	// parse the incoming message
	m := protocol.ParseMessage(msg)

	logrus.Debug("Logic: handling incoming message")

	// check if the message is not cached and relay it, otherwise ignore it
	if !n.isMessageCached(&m) {

		// append local node to path
		m.Path += ";" + n.settingsStation.Callsign
		m.PathLength = uint16(len(m.Path))

		// cache message
		n.cacheMessage(&m)

		// spread the message to peers
		n.spreadCachedMessage(&m)
	} else {
		logrus.Debug("Logic: message already cached, ignoring")
	}
}
