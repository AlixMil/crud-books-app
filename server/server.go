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
	GetBooks(c echo.Context) error
	UpdateBook(c echo.Context) error
	DeleteBook(c echo.Context) error
}

func (s Server) UseRouters(handlers Handlers) {
	s.server.Add(http.MethodPost, "/login", handlers.SignIn)
	s.server.Add(http.MethodPost, "/register", handlers.SignUp)
	s.server.Add(http.MethodGet, "/books", handlers.GetBooks)
	s.server.Add(http.MethodPost, "/files", handlers.UploadFile)
	s.server.Add(http.MethodPost, "/books", handlers.CreateBook)
	s.server.Add(http.MethodGet, "/books/:id", handlers.GetBook)
	s.server.Add(http.MethodPatch, "/books/:id", handlers.UpdateBook)
	s.server.Add(http.MethodDelete, "/books/:id", handlers.DeleteBook)
}

func (s Server) InitMiddlewares() {
	s.server.Use(echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			if c.Request().URL.Path == "/books" && c.Request().Header.Get("Authorization") == "" {
				return nil
			}

			return fmt.Errorf("JWT token invalid or expired")
		},
		ContinueOnIgnoredError: true,
		SigningKey:             []byte(s.jwtSecret),
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			if path == "/login" || path == "/register" {
				return true
			}
			return false
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt_package.JwtCustomClaims)
		},
	}))
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
