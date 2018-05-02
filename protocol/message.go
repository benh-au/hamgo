package protocol

import (
	"encoding/binary"

	"github.com/Sirupsen/logrus"

	"github.com/donothingloop/hamgo/parameters"
)

// PayloadType defines the type of the payload
type PayloadType uint16

// Payload types.
const (
	PayloadCQ                 = 0
	PayloadDebug              = 1
	PayloadUpd                = 2
	PayloadAck                = 3
	PayloadMessengerCQ        = 4
	PayloadMessengerGroup     = 5
	PayloadMessengerBroadcast = 6
	PayloadMessengerEmergency = 7
)

// Flags for the protocol.
const (
	FlagNoCache = (1 << 0)
	FlagACK     = (1 << 1)
)

// Message is a message in the transport.
type Message struct {
	Version       uint16      `json:"version"`
	SeqCounter    uint64      `json:"sequence"`
	TTL           uint8       `json:"ttl"`
	Flags         uint8       `json:"flags"`
	Source        Contact     `json:"source"`
	PathLength    uint16      `json:"pathLength"`
	Path          string      `json:"path"`
	PayloadType   PayloadType `json:"payloadType"`
	PayloadLenght uint32      `json:"payloadLength"`
	Payload       []byte      `json:"payload"`
}

// Bytes converts the message into a byte buffer.
func (m *Message) Bytes() []byte {
	buf := make([]byte, parameters.TransportMaxPackageSize)
	idx := 0

	binary.LittleEndian.PutUint16(buf[idx:], m.Version)
	idx += 2

	binary.LittleEndian.PutUint64(buf[idx:], m.SeqCounter)
	idx += 8

	buf[idx] = m.TTL
	idx++

	buf[idx] = m.Flags
	idx++

	cb := m.Source.Bytes()
	copy(buf[idx:], cb)
	idx += len(cb)

	binary.LittleEndian.PutUint16(buf[idx:], m.PathLength)
	idx += 2

	if m.PathLength != 0 {
		pth := []byte(m.Path)
		copy(buf[idx:], pth[0:int(m.PathLength)])
		idx += len(m.Path)
	}

	buf[idx] = uint8(m.PayloadType)
	idx++

	binary.LittleEndian.PutUint32(buf[idx:idx+4], m.PayloadLenght)
	idx += 4

	copy(buf[idx:], m.Payload)
	idx += len(m.Payload)

	return buf[:idx]
}

// ParseMessage parses a message from a buffer.
func ParseMessage(buf []byte) (*Message, []byte) {
	msg := Message{}
	idx := 0

	if len(buf) < 2+8+1+1 {
		logrus.Warn("Message: failed to parse message")
		return nil, nil
	}

	msg.Version = binary.LittleEndian.Uint16(buf[idx : idx+2])
	idx += 2

	msg.SeqCounter = binary.LittleEndian.Uint64(buf[idx : idx+8])
	idx += 8

	msg.TTL = buf[idx]
	idx++

	msg.Flags = buf[idx]
	idx++

	ct, rbuf := ParseContact(buf[idx:])
	if ct == nil {
		logrus.Warn("Message: failed to parse contact")
		return nil, nil
	}

	buf = rbuf
	idx = 0
	msg.Source = *ct

	if len(buf) < 2 {
		logrus.Warn("Message: failed to parse message")
		return nil, nil
	}

	msg.PathLength = binary.LittleEndian.Uint16(buf[idx:])
	idx += 2

	if msg.PathLength != 0 {

		if len(buf) < 2+int(msg.PathLength) {
			logrus.Warn("Message: failed to parse path")
			return nil, nil
		}

		msg.Path = string(buf[idx : idx+int(msg.PathLength)])
		idx += int(msg.PathLength)
	}

	if len(buf) < idx+1+4 {
		logrus.Warn("Message: failed to parse path length")
		return nil, nil
	}

	msg.PayloadType = PayloadType(buf[idx])
	idx++

	msg.PayloadLenght = binary.LittleEndian.Uint32(buf[idx : idx+4])
	idx += 4

	if len(buf) < idx+int(msg.PayloadLenght) {
		logrus.Warnf("Message: failed to parse payload, buffer smaller than payload length: %d < %d", len(buf), idx+int(msg.PayloadLenght))
		return nil, nil
	}

	pbuf := make([]byte, msg.PayloadLenght)
	for i := uint32(0); i < msg.PayloadLenght; i++ {
		pbuf[i] = buf[idx]
		idx++
	}
	msg.Payload = pbuf

	return &msg, buf[idx:]
}
