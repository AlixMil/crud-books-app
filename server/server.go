package server

import (
	jwt_package "crud-books/pkg/jwt"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port      string
	server    *echo.Echo
	jwtSecret string
}

type Handlers interface {
	UploadFile(c echo.Context) error
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	CreateBook(c echo.Context) error
	GetBook(c echo.Context) error
	GetBooksPublic(c echo.Context) error
	GetBooksOfUser(c echo.Context) error
	UpdateBook(c echo.Context) error
	DeleteBook(c echo.Context) error
	TestAuth(c echo.Context) error
}

func (s Server) UseRouters(handlers Handlers) {
	// public paths
	s.server.Add(http.MethodPost, "/login", handlers.SignIn)
	s.server.Add(http.MethodPost, "/register", handlers.SignUp)
	s.server.Add(http.MethodPost, "/books", handlers.GetBooksPublic)
	// private paths
	s.server.Add(http.MethodPost, "/files", handlers.UploadFile)
	s.server.Add(http.MethodPost, "/books", handlers.CreateBook)
	s.server.Add(http.MethodPost, "/books/:id", handlers.GetBook)
	s.server.Add(http.MethodPost, "/books", handlers.GetBooksOfUser)
	s.server.Add(http.MethodPatch, "/books/:id", handlers.UpdateBook)
	s.server.Add(http.MethodDelete, "/books/:id", handlers.DeleteBook)

	s.server.Add(http.MethodGet, "/testAuth", handlers.TestAuth)
}

func (s Server) InitMiddlewares() {
	// JWT Auth settings
	s.server.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(s.jwtSecret),
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			if path == "/login" || path == "/register" || path == "/books" {
				return true
			}
			return false
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt_package.JwtCustomClaims)
		},
	}))
	// CORS settings
	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}

func (s Server) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}

func New(port string, JwtSecret string) Server {
	server := Server{
		port:      port,
		server:    echo.New(),
		jwtSecret: JwtSecret,
	}

	return server
}
