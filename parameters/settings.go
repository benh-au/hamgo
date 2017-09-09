package parameters

// PeerSettings stores the settings for a peer.
type PeerSettings struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

// LogicSettings provides settings for the node logic component.
type LogicSettings struct {
	// CacheSize in messages
	CacheSize uint `json:"cacheSize"`
	ReadOnly  bool `json:"readonly,omitempty"`
}

// Settings stores the settings of the node.
type Settings struct {
	Port             uint           `json:"port"`
	PeerQueueSize    uint           `json:"peerQueueSize"`
	Retries          uint           `json:"retries"`
	Peers            []PeerSettings `json:"peers"`
	ReconnectTimeout uint           `json:"reconnectTimeout"`
	LogicSettings    LogicSettings  `json:"logic"`
}
