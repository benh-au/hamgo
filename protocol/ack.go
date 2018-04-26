package protocol

import (
	"encoding/binary"

	"github.com/Sirupsen/logrus"
)

// ACKPayload represents an ack message.
type ACKPayload struct {
	Source     Contact
	SeqCounter uint64
}

// Bytes converts to payload to bytes.
func (a *ACKPayload) Bytes() []byte {
	ctbuf := a.Source.Bytes()
	buf := make([]byte, len(ctbuf)+8)
	idx := 0

	copy(buf[idx:], ctbuf)
	idx += len(ctbuf)

	binary.LittleEndian.PutUint64(buf[idx:idx+8], a.SeqCounter)
	idx += 8

	return buf
}

// ParseACKPayload parses an ACK payload.
func ParseACKPayload(buf []byte) *ACKPayload {
	ack := &ACKPayload{}

	ct, rbuf := ParseContact(buf)
	if ct == nil {
		logrus.Warn("ACK: failed to parse corrupt contact, skipping")
		return nil
	}

	if len(rbuf) < 8 {
		logrus.Warn("ACK: failed to parse, skipping msg")
		return nil
	}

	idx := 0
	ack.Source = *ct

	ack.SeqCounter = binary.LittleEndian.Uint64(rbuf[idx : idx+8])
	idx += 8

	return ack
}
