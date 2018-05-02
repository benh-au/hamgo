package rest

import (
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/donothingloop/hamgo/node"

	"github.com/donothingloop/hamgo/protocol"
	"github.com/gorilla/websocket"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const protocolVersion = 1

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

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

	flags := uint8(0)

	if msg.ACK {
		flags |= protocol.FlagACK
	}

	// build the network message
	nmsg := protocol.Message{
		Version:    protocolVersion,
		SeqCounter: msg.Sequence,
		Source:     ctg,
		TTL:        255,
		Flags:      flags,

		PayloadType:   protocol.PayloadCQ,
		PayloadLenght: uint32(len(msg.Message)),
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
		str := messageToJSON(m)
		if str == "" {
			continue
		}

		if first {
			first = false
		} else {
			response += ", "
		}

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

func (h *Handler) ws(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()
	closech := make(chan interface{})
	closed := false
	lck := sync.Mutex{}

	cb := func(msg *protocol.Message, src *node.Peer) {
		if closed {
			return
		}

		logrus.Info("REST: sending message to websocket")

		str := messageToJSON(msg)
		err := ws.WriteMessage(websocket.TextMessage, []byte(str))

		if err != nil {
			lck.Lock()
			if !closed {
				closed = true
				close(closech)
			}
			lck.Unlock()
		}
	}

	go func() {
		for !closed {
			msg := &protocol.Message{}
			err := ws.ReadJSON(msg)
			if err != nil {
				logrus.WithError(err).Warn("REST: failed to read incoming message")
				return
			}

			logrus.Debugf("REST: spreading msg:\n %+v", msg)
			err = h.node.SpreadMessage(msg)
			if err != nil {
				logrus.WithError(err).Warn("REST: failed to send msg from ws")
				continue
			}

			logrus.Info("REST: message from ws spread")
		}
	}()

	cd := &node.MessageCallback{
		Cb: cb,
	}

	h.node.AddCallback(cd)
	defer h.node.RemoveCallback(cd)

	<-closech

	return nil
}

func (h *Handler) registerAPI(e *echo.Group) {
	spread := e.Group("/spread")
	spread.POST("/cq", h.cqmessage)

	e.GET("/cache", h.cache)
	e.GET("/ws", h.ws)
}
