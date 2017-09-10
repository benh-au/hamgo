package updproto

import (
	"github.com/Sirupsen/logrus"
	"github.com/donothingloop/hamgo/node"
	"github.com/donothingloop/hamgo/protocol"
)

// Handler handles the upd proto.
type Handler struct {
	node *node.Node
}

// NewHandler creates a new update protocol handler.
func NewHandler(n *node.Node) *Handler {
	return &Handler{
		node: n,
	}
}

// PeerConnectedHandler handles peer connected messages from the node logic.
func (h *Handler) PeerConnectedHandler(peer *node.Peer) {
	qry := protocol.UpdPayloadCacheRequest{}

	// build the query message
	for _, e := range h.node.Cache {
		qe := protocol.UpdRequestCacheEntry{
			SeqCounter: e.SeqCounter,
			Source:     e.Source,
		}

		qry.Entries = append(qry.Entries, qe)
	}

	qry.NumEntries = uint32(len(qry.Entries))

	qbuf := qry.Bytes()
	umsg := protocol.UpdPayload{
		DataLength: uint16(len(qbuf)),
		Data:       qbuf,
		Operation:  protocol.UpdOperationCacheRequest,
	}

	ubuf := umsg.Bytes()

	// build the protcol message
	msg := protocol.Message{
		Version:       0,
		SeqCounter:    0,
		TTL:           0,
		Flags:         protocol.FlagNoCache,
		Source:        h.node.Local,
		PathLength:    0,
		Path:          "",
		PayloadType:   protocol.PayloadUpd,
		PayloadLenght: uint32(len(ubuf)),
		Payload:       ubuf,
	}

	msgbuf := msg.Bytes()
	logrus.WithField("payload", qry).Debug("UpProto: sending query message")
	peer.QueueMessage(msgbuf)
}

// msgInRequest checks if a message is already cached on the querying node.
func (h *Handler) msgInRequest(msg *protocol.Message, req *protocol.UpdPayloadCacheRequest) bool {
	for _, m := range req.Entries {
		if (m.SeqCounter == msg.SeqCounter) && m.Source.Compare(&msg.Source) {
			return true
		}
	}

	return false
}

// handleRequest handles a request for the update protocol.
func (h *Handler) handleRequest(upd *protocol.UpdPayload, src *node.Peer) {
	req := protocol.ParsePayloadCacheRequest(upd.Data)
	res := protocol.UpdPayloadCacheResponse{}

	logrus.WithField("payload", req).Debug("UpProto: received query")

	for _, c := range h.node.Cache {
		if !h.msgInRequest(c, &req) {
			res.Entries = append(res.Entries, *c)
		}
	}

	res.NumEntries = uint32(len(res.Entries))

	pbuf := res.Bytes()
	updres := protocol.UpdPayload{
		Data:       pbuf,
		DataLength: uint16(len(pbuf)),
		Operation:  protocol.UpdOperationCacheResponse,
	}

	ubuf := updres.Bytes()

	// TODO: send message via direct message
	msg := protocol.Message{
		Version:    0,
		SeqCounter: 0,

		// TTL set to zero so that the message is not spread
		TTL:    0,
		Flags:  protocol.FlagNoCache,
		Source: h.node.Local,

		PathLength:    0,
		Path:          "",
		PayloadType:   protocol.PayloadUpd,
		PayloadLenght: uint32(len(ubuf)),
		Payload:       ubuf,
	}

	msgBuf := msg.Bytes()

	logrus.WithField("payload", res).Debug("UpProto: sending response message")
	src.QueueMessage(msgBuf)
}

func (h *Handler) handleResponse(upd *protocol.UpdPayload, src *node.Peer) {
	res := protocol.ParsePayloadCacheResponse(upd.Data)
	logrus.WithField("payload", res).Debug("UpProto: received response")

	for _, e := range res.Entries {
		logrus.WithField("entry", e).Debug("UpProto: received entry")
		logrus.Debug("UpProto: caching message")

		d := e
		h.node.AddToCache(&d)
	}
}

// UpdHandler handles messages for the update protocol.
func (h *Handler) UpdHandler(msg *protocol.Message, src *node.Peer) {
	// ignore non-hamgo messages
	if msg.PayloadType != protocol.PayloadUpd {
		return
	}

	logrus.Debug("UpProto: received message")

	upd, err := protocol.ParseUpdPayload(msg.Payload)
	if err != nil {
		logrus.WithError(err).Warn("UPDProto: failed to handle message")
		return
	}

	switch upd.Operation {
	case protocol.UpdOperationCacheRequest:
		h.handleRequest(upd, src)
		break

	case protocol.UpdOperationCacheResponse:
		h.handleResponse(upd, src)
		break
	}
}
