package protocol

import (
	"encoding/binary"

	"github.com/donothingloop/hamgo/parameters"
)

// PayloadType defines the type of the payload
type PayloadType uint16

// Payload types.
const (
	PayloadCQ    = 0
	PayloadDebug = 1
	PayloadHamgo = 2
)

// Flags for the protcol.
const (
	FlagNoCache = (1 << 1)
)

// Message is a message in the transport.
type Message struct {
	Version       uint16
	SeqCounter    uint32
	TTL           uint8
	Flags         uint8
	Source        Contact
	PathLength    uint16
	Path          string
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

	binary.LittleEndian.PutUint32(buf[idx:], m.SeqCounter)
	idx += 4

	buf[idx] = m.TTL
	idx++

	buf[idx] = m.Flags
	idx++

	cb := m.Source.Bytes()
	copy(buf[idx:], cb)
	idx += len(cb)

	binary.LittleEndian.PutUint16(buf[idx:], m.PathLength)
	idx += 2

	copy(buf[idx:], []byte(m.Path))
	idx += len(m.Path)

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

	msg.SeqCounter = binary.LittleEndian.Uint32(buf[idx : idx+4])
	idx += 4

	msg.TTL = buf[idx]
	idx++

	msg.Flags = buf[idx]
	idx++

	ct, rbuf := ParseContact(buf[idx:])
	buf = rbuf
	idx = 0
	msg.Source = ct

	msg.PathLength = binary.LittleEndian.Uint16(buf[idx:])
	idx += 2

	msg.Path = string(buf[idx : idx+int(msg.PathLength)])
	idx += int(msg.PathLength)

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
