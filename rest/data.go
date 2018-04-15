package rest

import (
	"encoding/json"
	"net"

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

func messageToCQ(msg *protocol.Message) string {
	var ips []string

	for _, v := range msg.Source.IPs {
		switch v.Type {
		case protocol.ContactIPv4:
			nip := net.IP(v.Data)
			ips = append(ips, nip.String())

			// TODO: IPv6
		}
	}

	cm := CQMessage{
		Sequence: msg.SeqCounter,
		Message:  string(msg.Payload),
		Contact: Contact{
			Callsign: string(msg.Source.Callsign),
			Type:     msg.Source.Type,
			IPs:      ips,
		},
		ACK: ((msg.Flags & protocol.FlagACK) != 0),
	}

	dat, err := json.Marshal(cm)
	if err != nil {
		logrus.WithError(err).Warn("REST: failed to convert message to json")
		return ""
	}

	return string(dat)
}

func messageToJSON(msg *protocol.Message) string {
	switch msg.PayloadType {
	case protocol.PayloadCQ:
		return messageToCQ(msg)
	}

	return ""
}

/*
	Example CQ message:

	{
		"sequence": 1,
		"contact": {
			"type": 0,
			"ips": [
				"44.143.25.1"
			],
			"callsign": "OE1VQS"
		},
		"message": "Test CQ"
	}

	Example CURL:
	curl -H "Content-Type: application/json" -X POST -d '{"sequence": 1, "contact": { "type": 0, "ips": ["44.143.25.1"], "callsign": "OE1VQS"}, "message": "Test CQ"}' http://server:9125/api/spread/cq
*/
