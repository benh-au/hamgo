package protocol

import (
	"encoding/binary"

	"github.com/donothingloop/hamgo/parameters"
)

// PayloadType defines the type of the payload
type PayloadType uint16

// Payload types.
const (
	PayloadCQ = 0
)

// Message is a message in the transport.
type Message struct {
	Version       uint16
	SeqCounter    uint16
	Source        Contact
	PayloadType   PayloadType
	PayloadLenght uint8
	Payload       []byte
}

// Bytes converts the message into a byte buffer.
func (m *Message) Bytes() []byte {
	buf := make([]byte, parameters.TransportMaxPackageSize)
	idx := 0

	binary.LittleEndian.PutUint16(buf[idx:], m.Version)
	idx += 2

	binary.LittleEndian.PutUint16(buf[idx:], m.SeqCounter)
	idx += 2

	cb := m.Source.Bytes()
	copy(buf[idx:], cb)
	idx += len(cb)

	buf[idx] = uint8(m.PayloadType)
	idx++

	buf[idx] = m.PayloadLenght
	idx++

	copy(buf[idx:], m.Payload)
	idx += len(m.Payload)

	return buf[:idx]
}

// ParseMessage parses a message from a buffer.
func ParseMessage(buf []byte) Message {
	msg := Message{}
	idx := 0

	msg.Version = binary.LittleEndian.Uint16(buf[idx : idx+2])
	idx += 2

	msg.SeqCounter = binary.LittleEndian.Uint16(buf[idx : idx+2])
	idx += 2

	ct, rbuf := ParseContact(buf[idx:])
	buf = rbuf
	idx = 0
	msg.Source = ct

	msg.PayloadType = PayloadType(buf[idx])
	idx++

	msg.PayloadLenght = buf[idx]
	idx++

	pbuf := make([]byte, msg.PayloadLenght)
	for i := uint8(0); i < msg.PayloadLenght; i++ {
		pbuf[i] = buf[idx]
		idx++
	}
	msg.Payload = pbuf

	return msg
}
