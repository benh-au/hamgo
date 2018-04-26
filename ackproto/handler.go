package ackproto

import (
	"github.com/Sirupsen/logrus"
	"github.com/donothingloop/hamgo/node"
	"github.com/donothingloop/hamgo/protocol"
)

// Handler is a handler for ACK payloads.
type Handler struct {
	node *node.Node
}

// NewHandler creates a new handler for the ACK protocol.
func NewHandler(n *node.Node) *Handler {
	return &Handler{
		node: n,
	}
}

// ACKHandler provides a handler for ack payloads.
func (h *Handler) ACKHandler(msg *protocol.Message, src *node.Peer) {
	// ignore non-ack messages
	if msg.PayloadType != protocol.PayloadAck {
		return
	}

	// catch empty messages
	if src == nil || msg == nil {
		return
	}

	logrus.Debug("ACKHandler: received message")

	ack := protocol.ParseACKPayload(msg.Payload)
	if ack == nil {
		logrus.Warn("ACKHandler: failed to handle messages")
		return
	}

	logrus.WithFields(logrus.Fields{
		"Reply from":   string(msg.Source.Callsign),
		"Source":       string(ack.Source.Callsign),
		"Contact Type": ack.Source.Type,
		"Sequence":     ack.SeqCounter,
	}).Info("ACKHandler: ACK msg received")
}
