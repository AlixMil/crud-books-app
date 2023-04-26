package handlers

import (
	"bytes"
	"crud-books/models"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	searchParam    = "search"
	limitParam     = "limit"
	offsetParam    = "offset"
	sortParam      = "sort"
	directionParam = "direction"

	fileFieldKey = "file"

	bindingError              = "binding request body error: %s"
	getUserIdFromCtxError     = "getting user id error: %s"
	createBookError           = "create book error: %s"
	parseMultipartError       = "read multipart error: %s"
	extensionCheckError       = "provided file should be PDF"
	readFileContentError      = "reading file content error: %s"
	serviceUploadFileError    = "upload file error: %s"
	serviceGetUserByIdError   = "get user by id failed, err: %s"
	getBookError              = "bookdata error: %s"
	getBooksParamsError       = "get books params error: %s"
	getBooksServiceError      = "get books service error: %s"
	getBooksParamsFillerError = "get books params filler error: %s"
	getBooksPrivateError      = "get books private error: %s"
	getBooksPublicError       = "get books public error: %s"
	serviceSignUpError        = "service signup error: %s"
	serviceSignInError        = "service signin error: %s"
	serviceUpdateBookError    = "service layer updatebook error: %s"
	serviceDeleteBookError    = "service delete book error: %s"
)

type Handlers struct {
	Services Service
}

//go:generate mockgen -source=handlers.go -destination=./mocks/handlers_mock.go
type Service interface {
	SignIn(user models.UserDataInput) (string, error)
	SignUp(user models.UserDataInput) (string, error)

	UploadFile(file []byte, fileName string) (string, error)

	GetBooks(filter models.Filter, sorting models.Sort) (*[]models.BookData, error)
	GetBook(bookToken string) (*models.GetBookResponse, error)
	CreateBook(title, description, fileToken, userEmail string) (string, error)
	UpdateBook(bookFileToken string, updater models.BookDataUpdater) error
	DeleteBook(tokenBook string) error

	GetUserById(userId string) (*models.UserData, error)
}

func New(serviceLayer Service) *Handlers {
	return &Handlers{
		Services: serviceLayer,
	}
}

func (e *Handlers) UploadFile(c echo.Context) error {
	c.Request().ParseMultipartForm(32 << 20)
	file, fileHeader, err := c.Request().FormFile(fileFieldKey)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(parseMultipartError, err.Error()))
	}

	fileName := fileHeader.Filename

	ext := filepath.Ext(fileName)
	if ext != ".pdf" {
		return c.String(http.StatusBadRequest, extensionCheckError)
	}

	buf := bytes.NewBuffer(nil)
	byteCont, err := io.ReadAll(file)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(readFileContentError, err.Error()))
	}
	buf.Write(byteCont)

	fileToken, err := e.Services.UploadFile(buf.Bytes(), fileName)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(serviceUploadFileError, err.Error()))
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
		return "", fmt.Errorf("convert user item to jwt format failed")
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("claims type assertion failed")
	}
	userId := claims["id"].(string)
	if userId == "" {
		return "", fmt.Errorf("userId is empty")
	}
	return userId, nil
}

func (e *Handlers) CreateBook(c echo.Context) error {
	var reqBody models.CreateBookRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(bindingError, err.Error()))
	}

	userId, err := getUserIdFromCtx(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(getUserIdFromCtxError, err.Error()))
	}

	userData, err := e.Services.GetUserById(userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(serviceGetUserByIdError, err.Error()))
	}

	fileToken, err := e.Services.CreateBook(reqBody.Title, reqBody.Description, reqBody.FileToken, userData.Email)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(createBookError, err.Error()))
	}
	return c.String(http.StatusOK, fileToken)
}

func (e *Handlers) GetBook(c echo.Context) error {
	bookToken := c.Param("id")
	bookData, err := e.Services.GetBook(bookToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(getBookError, err.Error()))
	}

	return c.JSON(http.StatusOK, bookData)
}

func getBooksParamsFieldFiller(c echo.Context, userEmail string) (*models.Filter, *models.Sort, error) {
	defaultLimit := 10
	defaultOffset := 0

	filter := models.Filter{
		Email:  userEmail,
		Search: c.QueryParams().Get(searchParam),
	}

	intLimit, err := strconv.Atoi(c.QueryParams().Get(limitParam))
	if err != nil {
		intLimit = defaultLimit
	}

	intOffset, err := strconv.Atoi(c.QueryParams().Get(offsetParam))
	if err != nil {
		intOffset = defaultOffset
	}

	sorting := models.Sort{
		SortField: c.QueryParams().Get(sortParam),
		Limit:     intLimit,
		Direction: c.QueryParams().Get(directionParam),
		Offset:    intOffset,
	}

	return &filter, &sorting, nil
}

func (e Handlers) getBooksPublic(c echo.Context) (*[]models.BookData, error) {
	filter, sort, err := getBooksParamsFieldFiller(c, "")
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksParamsError, err.Error()))
	}

	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksServiceError, err.Error()))
	}

	return books, nil

}

func (e Handlers) getBooksPrivate(c echo.Context, userId string) (*[]models.BookData, error) {
	userData, err := e.Services.GetUserById(userId)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf(serviceGetUserByIdError, err.Error()))
	}

	filter, sort, err := getBooksParamsFieldFiller(c, userData.Email)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksParamsFillerError, err.Error()))
	}

	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksServiceError, err.Error()))
	}

	return books, nil

}

func (e Handlers) GetBooks(c echo.Context) error {
	userId, err := getUserIdFromCtx(c)
	if err != nil {
		log.Errorf("error: %w", err)
		books, err := e.getBooksPublic(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksPublicError, err.Error()))
		}

		return c.JSON(http.StatusOK, books)
	}
	books, err := e.getBooksPrivate(c, userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(getBooksPrivateError, err.Error()))
	}

	return c.JSON(http.StatusOK, books)
}

func (e *Handlers) SignUp(c echo.Context) error {
	var signUpData models.UserDataInput
	err := c.Bind(&signUpData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(bindingError, err.Error()))
	}

	jwtTok, err := e.Services.SignUp(signUpData)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(serviceSignUpError, err.Error()))
	}

	c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTok))

	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtTok,
	})
}

func (e *Handlers) SignIn(c echo.Context) error {
	var signInData models.UserDataInput
	err := c.Bind(&signInData)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(bindingError, err.Error()))
	}
	jwtTok, err := e.Services.SignIn(signInData)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(serviceSignInError, err.Error()))
	}

	c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTok))

	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtTok,
	})
}

func (e *Handlers) UpdateBook(c echo.Context) error {
	bookToken := c.Param("id")

	var updater models.BookDataUpdater
	err := c.Bind(&updater)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(bindingError, err.Error()))
	}

	err = e.Services.UpdateBook(bookToken, updater)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf(serviceUpdateBookError, err.Error()))
	}

	return c.String(http.StatusOK, "")
}

func (e *Handlers) DeleteBook(c echo.Context) error {
	bookToken := c.Param("id")

	err := e.Services.DeleteBook(bookToken)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf(serviceDeleteBookError, err.Error()))
	}

	return c.JSON(http.StatusOK, "")
}
