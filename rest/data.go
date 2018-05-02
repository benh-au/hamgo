package rest

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/donothingloop/hamgo/protocol"
)

// Contact for the rest api.
type Contact struct {
	Type     protocol.ContactType `json:"type"`
	IPs      []string             `json:"ips"`
	Callsign string               `json:"callsign"`
}

// CQMessage indicates the users location.
type CQMessage struct {
	Sequence uint64  `json:"sequence"`
	Contact  Contact `json:"contact"`
	Message  string  `json:"message"`
	ACK      bool    `json:"ack,omitempty"`
}

func messageToJSON(msg *protocol.Message) string {
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Warn("REST: failed to convert message to json")
		return ""
	}

	return string(data)
}
