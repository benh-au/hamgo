package protocol

import (
	"bytes"
	"hamgo/parameters"
)

// ContactType defines the type of the contact.
type ContactType uint8

// Contact types.
const (
	ContactTypeUser  = 0
	ContactTypeFixed = 1
)

// ContactIPType defines a ip address type.
type ContactIPType uint8

// IP address types.
const (
	ContactIPv4 = 0
	ContactIPv6 = 1
)

// ContactIP defines the ip address of the contact.
type ContactIP struct {
	Type   ContactIPType
	Length uint8
	Data   []byte
}

// Bytes converts the ip address to bytes.
func (c *ContactIP) Bytes() []byte {
	buf := make([]byte, 2+len(c.Data))
	idx := 0

	buf[idx] = uint8(c.Type)
	idx++

	buf[idx] = c.Length
	idx++

	for i := 0; i < len(c.Data); i++ {
		buf[idx] = c.Data[i]
		idx++
	}

	return buf
}

// ParseContactIP parses a contact ip and returns the read length.
func ParseContactIP(buf []byte) (ContactIP, int) {
	idx := 0
	ci := ContactIP{}

	ci.Type = ContactIPType(buf[idx])
	idx++

	ci.Length = buf[idx]
	idx++

	dbuf := make([]byte, ci.Length)
	for i := uint8(0); i < ci.Length; i++ {
		dbuf[i] = buf[idx]
		idx++
	}

	ci.Data = dbuf
	return ci, idx
}

// Contact represents a source in the gossip protocol.
type Contact struct {
	Type           ContactType
	CallsignLength uint8
	Callsign       []byte
	NumberIPs      uint8
	IPs            []ContactIP
}

func (c *Contact) equalIPs(other *Contact) bool {
	for i, ip := range c.IPs {
		oip := other.IPs[i]

		if (!bytes.Equal(ip.Data, oip.Data)) ||
			(ip.Length != oip.Length) ||
			(ip.Type != oip.Type) {
			return false
		}
	}

	return true
}

// Compare the one contact to the other.
func (c *Contact) Compare(other *Contact) bool {
	return (c.Type == other.Type) &&
		(c.CallsignLength == other.CallsignLength) &&
		(bytes.Equal(c.Callsign, other.Callsign)) &&
		(c.NumberIPs == other.NumberIPs) &&
		(c.equalIPs(other))
}

// Bytes converts the contact to bytes.
func (c *Contact) Bytes() []byte {
	buf := make([]byte, parameters.TransportMaxPackageSize)
	idx := 0

	buf[idx] = uint8(c.Type)
	idx++

	buf[idx] = c.CallsignLength
	idx++

	// copy the callsign
	for i := 0; i < len(c.Callsign); i++ {
		buf[idx] = c.Callsign[i]
		idx++
	}

	buf[idx] = c.NumberIPs
	idx++

	// iterate IPs
	for _, v := range c.IPs {
		ip := v.Bytes()
		copy(buf[idx:], ip)

		idx += len(ip)
	}

	return buf[:idx]
}

// ParseContact parses a contact from a buffer and returns the remainder.
func ParseContact(msg []byte) (Contact, []byte) {
	idx := 0
	c := Contact{}

	c.Type = ContactType(msg[idx])
	idx++

	c.CallsignLength = msg[idx]
	idx++

	// copy the callsign
	for i := uint8(0); i < c.CallsignLength; i++ {
		c.Callsign = append(c.Callsign, msg[idx])
		idx++
	}

	c.NumberIPs = msg[idx]
	idx++

	// parse the ip addresses
	for i := uint8(0); i < c.NumberIPs; i++ {
		ibuf, n := ParseContactIP(msg[idx:])
		idx += n

		c.IPs = append(c.IPs, ibuf)
	}

	return c, msg[idx:]
}
