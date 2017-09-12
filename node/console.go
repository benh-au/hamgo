package node

import (
	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
)

// consoleHandler prints received messages in a readable format
func (n *Node) consoleHandler(msg *protocol.Message) {
	logrus.WithFields(logrus.Fields{
		"Sequence":        msg.SeqCounter,
		"Version":         msg.Version,
		"Source Callsign": string(msg.Source.Callsign),
		"Source Type":     msg.Source.Type,
		"Flags":           msg.Flags,
		"Payload Lenght":  msg.PayloadLenght,
		"Payload Type":    msg.PayloadType,
		"Payload":         string(msg.Payload),
		"Path":            string(msg.Path),
	}).Info("message received")
}
