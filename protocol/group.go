package protocol

// GroupMessageSeverity defines a severity for a group message.
type GroupMessageSeverity uint8

// Message severities for group messages.
const (
	SeverityDebug     = 0
	SeverityInfo      = 1
	SeverityWarn      = 2
	SeverityCritical  = 3
	SeverityEmergency = 4
	SeverityNormal    = 5
)

// Group operations.
const (
	GroupOperationMessage    = 0
	GroupOperationMembership = 1
)

// GroupPayload defines a payload for a group message.
type GroupPayload struct {
	GroupLength uint8
	Group       string
	Timestamp   uint32
	Operation   uint8
	Payload     []byte
}

// GroupMessagePayload defines a payload that transports messages over the group payload.
type GroupMessagePayload struct {
	Type          GroupMessageSeverity
	MessageLength uint16
	Message       string
}

// GroupMembershipPayload is used for the group membership protocol.
type GroupMembershipPayload struct {
}
