package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserDB interface{}

type EchoHandlers struct {
	db             UserDB
	ServiceStorage ServiceStorage
}

type ServiceStorage interface {
	SaveFile(file []byte, title string) (string, error)
}

func New(db UserDB, service ServiceStorage) (*EchoHandlers, error) {
	return &EchoHandlers{db: db}, nil
}

func (e *EchoHandlers) MainPage(c echo.Context) error {
	c.String(http.StatusOK, "let's gooo")
	return nil
}
