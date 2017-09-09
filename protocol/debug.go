package protocol

// Debug operations.
const (
	DebugOperationVersion   = 0
	DebugOperationBroadcast = 1
)

// Debug defines a query message for the protocol.
type Debug struct {
	Operation uint8
}
