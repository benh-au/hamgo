package rest

import (
	"net"
	"strconv"

	"github.com/donothingloop/hamgo/protocol"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const protocolVersion = 1

// spread a cqmessage
func (h *Handler) cqmessage(c echo.Context) error {
	msg := CQMessage{}

	if err := c.Bind(&msg); err != nil {
		return err
	}

	ips := []protocol.ContactIP{}

	// build ip addresses
	for _, v := range msg.Contact.IPs {
		ip := net.ParseIP(v)

		// TODO: ipv6
		cip := protocol.ContactIP{
			Type:   protocol.ContactIPv4,
			Length: uint8(len(ip)),
			Data:   []byte(ip),
		}

		ips = append(ips, cip)
	}

	// build the network contact
	ctg := protocol.Contact{
		Type:           msg.Contact.Type,
		CallsignLength: uint8(len(msg.Contact.Callsign)),
		Callsign:       []byte(msg.Contact.Callsign),
		NumberIPs:      uint8(len(msg.Contact.IPs)),
		IPs:            ips,
	}

	// build the network message
	nmsg := protocol.Message{
		Version:       protocolVersion,
		SeqCounter:    msg.Sequence,
		Source:        ctg,
		PayloadType:   protocol.PayloadCQ,
		PayloadLenght: uint8(len(msg.Message)),
		Payload:       []byte(msg.Message),
	}

	logrus.WithField("msg", nmsg).Debug("spreading CQ message")

	// spread the message
	h.node.SpreadMessage(&nmsg)

	return c.NoContent(200)
}

// cache returns the current cache
func (h *Handler) cache(c echo.Context) error {
	max := c.QueryParam("max")
	maxI := -1

	if max != "" {
		mi, err := strconv.Atoi(max)
		maxI = mi

		if err != nil {
			return err
		}
	}

	response := "["
	first := true
	cnt := 0

	for _, m := range h.node.Cache {
		if first {
			first = false
		} else {
			response += ", "
		}

		str := messageToJSON(m)
		response += str

		cnt++

		// limit results if max is set
		if cnt > maxI && maxI != -1 {
			break
		}
	}

	response += "]"

	return c.String(200, response)
}

func (h *Handler) registerAPI(e *echo.Group) {
	spread := e.Group("/spread")
	spread.POST("/cq", h.cqmessage)

	e.GET("/cache", h.cache)
}
