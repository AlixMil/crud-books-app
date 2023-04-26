package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port          string
	server        *echo.Echo
	jwtSecret     string
	jwtMiddleware echo.MiddlewareFunc
}

type Handlers interface {
	UploadFile(c echo.Context) error
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	CreateBook(c echo.Context) error
	GetBook(c echo.Context) error
	GetBooks(c echo.Context) error
	UpdateBook(c echo.Context) error
	DeleteBook(c echo.Context) error
}

func (s Server) UseRouters(handlers Handlers) {
	s.server.Add(http.MethodPost, "/login", handlers.SignIn)
	s.server.Add(http.MethodPost, "/register", handlers.SignUp)

	s.server.Add(http.MethodPost, "/files", handlers.UploadFile)

	books := s.server.Group("/books")
	books.Add(http.MethodPost, "", handlers.CreateBook)
	books.Add(http.MethodGet, "", handlers.GetBooks)
	books.Add(http.MethodGet, "/:id", handlers.GetBook)
	books.Add(http.MethodPut, "/:id", handlers.UpdateBook)
	books.Add(http.MethodDelete, "/:id", handlers.DeleteBook)
}

func (s Server) InitMiddlewares() {
	s.server.Use(s.jwtMiddleware)
	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}

func (s Server) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}

func New(port string, JwtSecret string, JwtMiddleware echo.MiddlewareFunc) Server {
	server := Server{
		port:          port,
		server:        echo.New(),
		jwtSecret:     JwtSecret,
		jwtMiddleware: JwtMiddleware,
	}

	return server
}
