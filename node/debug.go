package node

import (
	"github.com/Sirupsen/logrus"
	"github.com/donothingloop/hamgo/protocol"
)

func parseDebugPayload(buf []byte) *protocol.Debug {
	dbg := &protocol.Debug{
		Operation: buf[0],
	}

	return dbg
}

// debugHandler handles incoming debug messages
func (n *Node) debugHandler(msg *protocol.Message) {
	// ignore non-debug messages
	if msg.PayloadType != protocol.PayloadDebug {
		return
	}

	logrus.Info("Received debug message")
	dbg := parseDebugPayload(msg.Payload)

	switch dbg.Operation {
	case protocol.DebugOperationBroadcast:
		// TODO: implement
		break
	}
}
