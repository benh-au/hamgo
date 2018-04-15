package protocol

import (
	"encoding/binary"
	"errors"
)

// Operations for hamgo protocol messages.
const (
	UpdOperationCacheRequest  = 0
	UpdOperationCacheResponse = 1
)

// UpdPayload defines the payload for hamgo signaling.
type UpdPayload struct {
	Operation  uint8
	DataLength uint16
	Data       []byte
}

// UpdRequestCacheEntry is sent in a cache request message to inform the remote peer
// of the messages that are already in the cache of the node so that only new messages
// have to be sent from the remote peer.
type UpdRequestCacheEntry struct {
	SeqCounter uint64
	Source     Contact
}

// UpdPayloadCacheRequest is the request sent by a node to get missing cache messages.
type UpdPayloadCacheRequest struct {
	NumEntries uint32
	Entries    []UpdRequestCacheEntry
}

// UpdPayloadCacheResponse is sent as an answer to a cache query in order to update the
// querying nodes cache.
type UpdPayloadCacheResponse struct {
	NumEntries uint32
	Entries    []Message
}

// Bytes converts a cache entry to bytes.
func (e *UpdRequestCacheEntry) Bytes() []byte {
	ct := e.Source.Bytes()
	buf := make([]byte, len(ct)+4)
	idx := 0

	binary.LittleEndian.PutUint64(buf[idx:idx+4], e.SeqCounter)
	idx += 4

	copy(buf[idx:], ct)
	return buf
}

// ParseCacheEntry parses a cache entry and returns the remaining buffer.
func ParseCacheEntry(buf []byte) (UpdRequestCacheEntry, []byte) {
	re := UpdRequestCacheEntry{}
	idx := 0

	re.SeqCounter = binary.LittleEndian.Uint64(buf[idx : idx+4])
	idx += 4

	ct, rbuf := ParseContact(buf[idx:])
	re.Source = ct

	return re, rbuf
}

// Bytes converts an update protocol request to a byte buffer.
func (r *UpdPayloadCacheRequest) Bytes() []byte {
	var ent [][]byte
	tlen := 0

	for _, c := range r.Entries {
		e := c.Bytes()
		tlen += len(e)
		ent = append(ent, e)
	}

	buf := make([]byte, tlen+4)
	idx := 0

	binary.LittleEndian.PutUint32(buf[idx:], r.NumEntries)
	idx += 4

	for _, e := range ent {
		// put entries
		copy(buf[idx:], e)

		idx += len(e)
	}

	return buf
}

// Bytes converts the cache response to a byte buffer.
func (r *UpdPayloadCacheResponse) Bytes() []byte {
	var mbufs [][]byte
	tlen := 0

	for _, v := range r.Entries {
		b := v.Bytes()
		mbufs = append(mbufs, b)
		tlen += len(b)
	}

	buf := make([]byte, tlen+4)
	idx := 0

	binary.LittleEndian.PutUint32(buf[idx:idx+4], r.NumEntries)
	idx += 4

	for _, e := range mbufs {
		copy(buf[idx:], e)
		idx += len(e)
	}

	return buf
}

// ParsePayloadCacheResponse parses a cache response.
func ParsePayloadCacheResponse(buf []byte) UpdPayloadCacheResponse {
	idx := 0
	pcr := UpdPayloadCacheResponse{}

	pcr.NumEntries = binary.LittleEndian.Uint32(buf[idx : idx+4])
	idx += 4

	for i := 0; i < int(pcr.NumEntries); i++ {
		m, rbuf := ParseMessage(buf[idx:])
		pcr.Entries = append(pcr.Entries, m)
		idx = 0
		buf = rbuf
	}

	return pcr
}

// ParsePayloadCacheRequest parses a cache request.
func ParsePayloadCacheRequest(buf []byte) UpdPayloadCacheRequest {
	cr := UpdPayloadCacheRequest{}
	idx := 0

	cr.NumEntries = binary.LittleEndian.Uint32(buf[idx : idx+4])
	idx += 4

	buf = buf[idx:]
	for i := 0; i < int(cr.NumEntries); i++ {
		e, rbuf := ParseCacheEntry(buf)
		cr.Entries = append(cr.Entries, e)

		buf = rbuf

		if len(buf) == 0 {
			break
		}
	}

	return cr
}

// Bytes converts a update protocol payload to a byte buffer.
func (u *UpdPayload) Bytes() []byte {
	buf := make([]byte, 3+int(u.DataLength))
	idx := 0

	buf[idx] = u.Operation
	idx++

	binary.LittleEndian.PutUint16(buf[idx:idx+2], u.DataLength)
	idx += 2

	copy(buf[idx:], u.Data)
	return buf
}

// ParseUpdPayload parses an update protocol payload.
func ParseUpdPayload(buf []byte) (*UpdPayload, error) {
	upd := UpdPayload{}
	idx := 0

	upd.Operation = buf[idx]
	idx++

	upd.DataLength = binary.LittleEndian.Uint16(buf[idx : idx+2])
	idx += 2

	if (len(buf) - int(idx)) < int(upd.DataLength) {
		return nil, errors.New("payload invalid")
	}

	upd.Data = buf[idx : idx+int(upd.DataLength)]
	idx += int(upd.DataLength)

	return &upd, nil
}
