package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	port   string
	server *echo.Echo
}

func New(port string) Server {
	server := Server{
		port:   port,
		server: echo.New(),
	}

	return server
}

type Handlers interface {
	MainPage(c echo.Context) error
}

func (s Server) InitHandlers(handlers Handlers) {
	s.server.Add(http.MethodGet, "/", handlers.MainPage)
}

func (s Server) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}
