package lib

import (
	"bufio"
	"errors"
	"net"

	"github.com/donothingloop/hamgo/parameters"

	"github.com/Sirupsen/logrus"
)

// Connection stores a connection handle with a peer.
type Connection struct {
	Connection  net.Conn
	EscapeNext  bool
	FrameActive bool
	Received    chan []byte
	Send        chan *Message
	close       chan interface{}
	Closed      bool
}

// Message is a message that is sent to the connection.
type Message struct {
	Data     []byte
	Callback func(*Connection, error)
}

// Close the connection.
func (c *Connection) Close() {
	c.Closed = true
	close(c.close)
	c.Connection.Close()
}

func (c *Connection) sendMessage(msg *Message) {
	logrus.WithField("msg", msg).Debug("Connection: sending message")

	buf := make([]byte, parameters.TransportMaxPackageSize)
	idx := 0

	buf[idx] = TransportFrameStart
	idx++

	if idx >= parameters.TransportMaxPackageSize {
		logrus.Warnf("Connection: package size exceeded, max: %d", parameters.TransportMaxPackageSize)
		msg.Callback(c, errors.New("package size exceeded"))
		return
	}

	for i := 0; i < len(msg.Data); i++ {
		if (msg.Data[i] == TransportFrameEnd) ||
			(msg.Data[i] == TransportFrameStart) ||
			(msg.Data[i] == TransportEscape) {
			buf[idx] = TransportEscape
			idx++

			if idx >= parameters.TransportMaxPackageSize {
				logrus.Warnf("Connection: package size exceeded, max: %d", parameters.TransportMaxPackageSize)
				msg.Callback(c, errors.New("package size exceeded"))
				return
			}
		}

		buf[idx] = msg.Data[i]
		idx++

		if idx >= parameters.TransportMaxPackageSize {
			logrus.Warnf("Connection: package size exceeded, max: %d", parameters.TransportMaxPackageSize)
			msg.Callback(c, errors.New("package size exceeded"))
			return
		}
	}

	buf[idx] = TransportFrameEnd
	idx++

	if idx >= parameters.TransportMaxPackageSize {
		logrus.Warnf("Connection: package size exceeded, max: %d", parameters.TransportMaxPackageSize)
		msg.Callback(c, errors.New("package size exceeded"))
		return
	}

	sbuf := buf[:idx]
	logrus.WithField("buf", sbuf).Debug("Connection: wire message built")

	logrus.WithField("len", len(sbuf)).Debug("Connection: sending wire message")

	// send the message
	n, err := c.Connection.Write(sbuf)
	if n != idx {
		logrus.Warnf("Connection: message not completely sent, sent: %d, should: %d", n, idx)
	}

	logrus.Debug("Connection: result callback")

	// call the result callback
	msg.Callback(c, err)
}

func (c *Connection) sendWorker() {
	logrus.Debug("Connection: sendWorker: active")

	for {
		select {
		case msg := <-c.Send:
			c.sendMessage(msg)
			break

		// close on request
		case <-c.close:
			logrus.Debug("Connection: sendWorker: closing")
			return
		}
	}
}

func (c *Connection) connectionWorker() {
	logrus.Debug("Connection: starting worker")

	// start the send worker
	go c.sendWorker()

	rd := bufio.NewReader(c.Connection)

	// make a buffer that can hold the maximum package size
	buf := make([]byte, parameters.TransportMaxPackageSize)
	idx := 0

	for {
		// check if the package exceeds the buffer bounds and drop it
		if idx >= len(buf) {
			logrus.Warnf("Connection: dropping message as it exceeds the maximum package size of %d", parameters.TransportMaxPackageSize)
			idx = 0
			continue
		}

		b, err := rd.ReadByte()
		logrus.WithField("byte", b).Debug("Connection: read byte")

		if err != nil {
			logrus.WithError(err).Warn("Connection: failed to read byte, closing connection")
			c.Close()
			return
		}

		// check if the byte should be escaped
		if c.EscapeNext {
			logrus.Debug("Connection: byte should be escaped")
			c.EscapeNext = false
			buf[idx] = b
			idx++
			continue
		}

		// check if the next byte should be escaped
		if b == TransportEscape {
			logrus.Debug("Connection: next byte should be escaped")
			c.EscapeNext = true
			continue
		}

		// check if this is a startframe marker
		if b == TransportFrameStart {
			logrus.Debug("Connection: start frame")

			if c.FrameActive {
				logrus.Warn("Connection: received start frame while frame was still pending")
			}

			c.FrameActive = true
			idx = 0
			continue
		}

		// check if this is an endframe marker
		if b == TransportFrameEnd {
			logrus.Debug("Connection: end frame")

			if !c.FrameActive {
				logrus.Warn("Connection: recevied endframe marker while no frame was active")
				continue
			}

			logrus.WithFields(logrus.Fields{
				"pkgsize": idx,
			}).Debugf("Connection: received frame end")

			// send the received data
			c.Received <- buf[:idx]

			// create a new buffer
			buf = make([]byte, parameters.TransportMaxPackageSize)
			idx = 0

			c.FrameActive = false
		}

		logrus.Debug("Connection: data byte")
		buf[idx] = b
		idx++
	}
}
