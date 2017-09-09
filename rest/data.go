package rest

import "github.com/donothingloop/hamgo/protocol"

// Contact for the rest api.
type Contact struct {
	Type     protocol.ContactType `json:"type"`
	IPs      []string             `json:"ips"`
	Callsign string               `json:"callsign"`
}

// CQMessage indicates the users location.
type CQMessage struct {
	Sequence uint16  `json:"sequence"`
	Contact  Contact `json:"contact"`
	Message  string  `json:"message"`
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
