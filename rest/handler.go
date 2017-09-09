package rest

import "github.com/donothingloop/hamgo/node"

// Handler stores handlers for the rest server.
type Handler struct {
	node *node.Node
}

// NewHandler creates a new handler for the REST server.
func NewHandler(n *node.Node) *Handler {
	return &Handler{
		node: n,
	}
}
