package protocol

import (
	"encoding/binary"

	"github.com/donothingloop/hamgo/parameters"
)

// PayloadType defines the type of the payload
type PayloadType uint16

// Payload types.
const (
	PayloadCQ           = 0
	PayloadDebug        = 1
	PayloadUpd          = 2
	PayloadAck          = 3
	PayloadGroupMessage = 4
)

// Flags for the protcol.
const (
	FlagNoCache = (1 << 0)
	FlagACK     = (1 << 1)
)

// Message is a message in the transport.
type Message struct {
	Version       uint16
	SeqCounter    uint64
	TTL           uint8
	Flags         uint8
	Source        Contact
	PathLength    uint16
	Path          string
	PayloadType   PayloadType
	PayloadLenght uint32
	Payload       []byte
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
func ParseMessage(buf []byte) (Message, []byte) {
	msg := Message{}
	idx := 0

	msg.Version = binary.LittleEndian.Uint16(buf[idx : idx+2])
	idx += 2

	msg.SeqCounter = binary.LittleEndian.Uint64(buf[idx : idx+8])
	idx += 8

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

	if msg.PathLength != 0 {
		msg.Path = string(buf[idx : idx+int(msg.PathLength)])
		idx += int(msg.PathLength)
	}

	msg.PayloadType = PayloadType(buf[idx])
	idx++

	msg.PayloadLenght = binary.LittleEndian.Uint32(buf[idx : idx+8])
	idx += 8

	pbuf := make([]byte, msg.PayloadLenght)
	for i := uint32(0); i < msg.PayloadLenght; i++ {
		pbuf[i] = buf[idx]
		idx++
	}
	msg.Payload = pbuf

	return msg, buf[idx:]
}
