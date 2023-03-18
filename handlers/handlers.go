package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoHandlers struct {
	db             UserDB
	ServiceStorage ServiceStorage
}

type UserDB interface{}

type ServiceStorage interface {
	UploadFile(file []byte) (string, error)
}

func New(db UserDB, service ServiceStorage) (*EchoHandlers, error) {
	return &EchoHandlers{db: db}, nil
}

func (e *EchoHandlers) MainPage(c echo.Context) error {
	c.String(http.StatusOK, "let's gooo")
	return nil
}

type UploadRequest struct {
}

type UploadFormData struct {
	File []byte `form: file`
}

type FormFileData struct {
	ContentType string
	Content     []byte
	Name        string
}

func (e *EchoHandlers) UploadFile(c echo.Context) error {
	// load file to storage -> return token of file
	// record to db file name and token
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.String(http.StatusBadRequest, "Please attach PDF file")
	}

	fileMultipart, err := file.Open()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, fileMultipart); err != nil {
		return err
	}

	_, err = e.ServiceStorage.UploadFile(buf.Bytes())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("error in fileUpload: %v", err))
	}

	return nil
}
