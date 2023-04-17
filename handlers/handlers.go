package handlers

import (
	"bytes"
	"crud-books/models"
	jwt_package "crud-books/pkg/jwt"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type EchoHandlers struct {
	Services Service
}

//go:generate mockgen -source=handlers.go -destination=./mocks/handlers_mock.go
type Service interface {
	SignIn(user models.UserDataInput) (string, error)
	SignUp(user models.UserDataInput) (string, error)
	CreateBook(title, description, fileToken, userEmail string) (string, error)
	UploadFile(file []byte, fileHeader *multipart.FileHeader) (string, error)
	GetBook(bookToken string) (*models.GetBookResponse, error)
	GetBooks(filter models.Filter, sorting models.Sort) (*[]models.BookData, error)
	UpdateBook(bookFileToken string, updater models.BookDataUpdater) error
	DeleteBook(tokenBook string) error
	GetUserByInsertedId(userId string) (*models.UserData, error)
}

func (e *EchoHandlers) UploadFile(c echo.Context) error {
	c.Request().ParseMultipartForm(32 << 20)
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("read multipart form failed, error: %s", err.Error()))
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".pdf" {
		return c.String(http.StatusBadRequest, "provided file should be PDF")
	}

	buf := bytes.NewBuffer(nil)
	byteCont, err := io.ReadAll(file)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("reading file failed, error: %s", err.Error()))
	}
	buf.Write(byteCont)

	fileToken, err := e.Services.UploadFile(buf.Bytes(), fileHeader)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("upload file failed, error: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"fileToken": fileToken,
	})
}

func getUserIdFromCtx(c echo.Context) (string, error) {
	uCtx := c.Get("user")
	if uCtx == nil {
		return "", fmt.Errorf("user context is null")
	}
	user, ok := uCtx.(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("convert ctx user to jwt format failed")
	}
	claims, ok := user.Claims.(*jwt_package.JwtCustomClaims)
	if !ok {
		return "", fmt.Errorf("claims type assertion falied")
	}
	userId := claims.UserId
	if userId == "" {
		return "", fmt.Errorf("userId is empty")
	}
	return userId, nil
}

func (e *EchoHandlers) CreateBook(c echo.Context) error {
	var reqBody models.CreateBookRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to check request body, err: %s", err.Error()))
	}

	userId, err := getUserIdFromCtx(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("get user id from ctx error: %s", err.Error()))
	}

	userData, err := e.Services.GetUserByInsertedId(userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("get user by inserted id failed, err: %s", err.Error()))
	}

	fileToken, err := e.Services.CreateBook(reqBody.Title, reqBody.Description, reqBody.FileToken, userData.Email)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("create book error: %s", err.Error()))
	}
	return c.String(http.StatusOK, fileToken)
}

func (e *EchoHandlers) GetBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookToken := strings.Replace(path, "/books/", "", 1)
	bookData, err := e.Services.GetBook(bookToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("attempting to receive book data from db failed, error: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, bookData)
}

func getBooksParamsFieldFiller(c echo.Context, userEmail string) (*models.Filter, *models.Sort, error) {
	defaultLimit := 10
	defaultOffset := 0

	filter := models.Filter{
		Email:  userEmail,
		Search: c.QueryParams().Get("search"),
	}

	intLimit, err := strconv.Atoi(c.QueryParams().Get("limit"))
	if err != nil {
		intLimit = defaultLimit
	}

	intOffset, err := strconv.Atoi(c.QueryParams().Get("offset"))
	if err != nil {
		intOffset = defaultOffset
	}

	sorting := models.Sort{
		SortField: c.QueryParams().Get("sort"),
		Limit:     intLimit,
		Direction: c.QueryParams().Get("direction"),
		Offset:    intOffset,
	}

	return &filter, &sorting, nil
}

func (e EchoHandlers) getBooksPublic(c echo.Context) (*[]models.BookData, error) {
	filter, sort, err := getBooksParamsFieldFiller(c, "")
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("get books params field filler error: %s", err.Error()))
	}

	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("get books service layer error: %s", err.Error()))
	}

	return books, nil

}

func (e EchoHandlers) getBooksPrivate(c echo.Context, userId string) (*[]models.BookData, error) {
	userData, err := e.Services.GetUserByInsertedId(userId)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("get user by inserted id failed, error: %s", err.Error()))
	}

	filter, sort, err := getBooksParamsFieldFiller(c, userData.Email)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("get books params field filler error: %s", err.Error()))
	}

	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("get books service layer error: %s", err.Error()))
	}

	return books, nil

}

func (e EchoHandlers) GetBooks(c echo.Context) error {
	userId, err := getUserIdFromCtx(c)
	if err != nil {
		books, err := e.getBooksPublic(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("get books public error: %s", err.Error()))
		}

		return c.JSON(http.StatusOK, books)
	}
	books, err := e.getBooksPrivate(c, userId)

	return c.JSON(http.StatusOK, books)
}

func (e *EchoHandlers) SignUp(c echo.Context) error {
	var signUpData models.UserDataInput
	err := c.Bind(&signUpData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to check request body, err: %s", err.Error()))
	}

	jwtTok, err := e.Services.SignUp(signUpData)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("service layer signup error: %s", err.Error()))
	}

	c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTok))

	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtTok,
	})
}

func (e *EchoHandlers) SignIn(c echo.Context) error {
	var signInData models.UserDataInput
	err := c.Bind(&signInData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to check request body, err: %s", err.Error()))
	}
	jwtTok, err := e.Services.SignIn(signInData)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("service layer sign in error: %s", err.Error()))
	}

	c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTok))

	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtTok,
	})
}

func (e *EchoHandlers) UpdateBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookFileToken := strings.Replace(path, "/books/", "", 1)

	var updater models.BookDataUpdater
	err := c.Bind(&updater)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to check request body, err: %s", err.Error()))
	}

	err = e.Services.UpdateBook(bookFileToken, updater)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("service layer updatebook error: %s", err.Error()))
	}

	return c.String(http.StatusOK, "Book data successfully updated!")
}

func (e *EchoHandlers) DeleteBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookId := strings.Replace(path, "/books/", "", 1)

	err := e.Services.DeleteBook(bookId)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("service layer delete book error: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, "")
}

func New(serviceLayer Service) *EchoHandlers {
	return &EchoHandlers{
		Services: serviceLayer,
	}
}
