package rest

import (
	"fmt"

	"github.com/donothingloop/hamgo/node"
	"github.com/donothingloop/hamgo/parameters"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

// Server provides a server for accessing the hamgo protocol.
type Server struct {
	settings parameters.RESTSettings
}

// NewServer creates a new rest server.
func NewServer(sett parameters.RESTSettings) *Server {
	return &Server{
		settings: sett,
	}
}

// Init the rest server.
func (r *Server) Init(n *node.Node) {
	logrus.Debug("RESTServer: starting")

	e := echo.New()

	if !r.settings.CORS {
		logrus.Debug("RESTServer: cors middleware enabled")
		e.Use(middleware.CORS())
	}

	hndlr := NewHandler(n)

	e.Use(middleware.Recover())
	hndlr.registerAPI(e.Group("/api"))

	e.Static("/", r.settings.Frontend)

	port := r.settings.Port
	err := e.Start(fmt.Sprintf(":%d", port))

	logrus.Debugf("RESTServer: listening on port %d", port)

	if err != nil {
		logrus.WithError(err).Warn("REST server error")
	}
}
