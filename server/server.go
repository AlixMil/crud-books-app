package server

import (
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

type jwtCustomClaims struct {
	Email  string
	UserId string
	jwt.RegisteredClaims
}

func New(port string, JwtSecret string) Server {
	server := Server{
		port:      port,
		server:    echo.New(),
		jwtSecret: JwtSecret,
	}

	return server
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
}

func (s Server) InitHandlers(handlers Handlers) {
	// JWT Auth settings
	s.server.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: s.jwtSecret,
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" || c.Request().URL.Path == "/register" || c.Request().URL.Path == "/books" {
				return true
			}
			return false
		},
	}))
	// CORS settings
	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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
}

func (s Server) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}
