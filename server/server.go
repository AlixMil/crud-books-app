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
}

func (s Server) InitHandlers(handlers Handlers) {
	s.server.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("SECRET!"),
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" || c.Request().URL.Path == "/register" {
				return true
			}
			return false
		},
	}))

	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	s.server.Add(http.MethodPost, "/login", handlers.SignIn)
	s.server.Add(http.MethodPost, "/register", handlers.SignUp)
	s.server.Add(http.MethodPost, "/files", handlers.UploadFile)
}

func (s Server) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}
