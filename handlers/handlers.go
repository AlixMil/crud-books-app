package handlers

import (
	"bytes"
	"crud-books/models"
	jwt_package "crud-books/pkg/jwt"
	"fmt"
	"io"
	"net/http"
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
	UploadFile(file []byte) (string, error)
	GetBook(bookToken string) (*models.GetBookResponse, error)
	GetBooks(filter models.Filter, sorting models.Sort) (*[]models.BookData, error)
	UpdateBook(bookId string, updater models.BookDataUpdater) error
	DeleteBook(tokenBook string) error
	GetUserByInsertedId(userId string) (*models.UserData, error)
}

func (e *EchoHandlers) UploadFile(c echo.Context) error {
	c.Request().ParseMultipartForm(32 << 20)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return fmt.Errorf("read multipart form failed, error: %w", err)
	}
	buf := bytes.NewBuffer(nil)
	byteCont, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file failed, error: %w", err)
	}
	buf.Write(byteCont)

	fileToken, err := e.Services.UploadFile(buf.Bytes())
	if err != nil {
		return fmt.Errorf("upload file failed, error: %w", err)
	}

	return c.String(http.StatusOK, fileToken)
}

func getUserIdFromCtx(context echo.Context) (string, error) {
	user, ok := context.Get("user").(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("get user id from ctx failed, jwt token missed")
	}
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)
	if userId == "" {
		return "", fmt.Errorf("userId is empty")
	}
	return userId, nil
}

func (e *EchoHandlers) CreateBook(c echo.Context) error {
	var reqBody models.CreateBookRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return fmt.Errorf("failed of reading (binding) request body in create book func of handlers. Error: %w", err)
	}
	userId, err := getUserIdFromCtx(c)
	if err != nil {
		return err
	}
	userData, err := e.Services.GetUserByInsertedId(userId)
	if err != nil {
		return fmt.Errorf("getting user data by jwt user id failed, error: %w", err)
	}

	fileToken, err := e.Services.CreateBook(reqBody.Title, reqBody.Description, reqBody.FileToken, userData.Email)
	if err != nil {
		return fmt.Errorf("create book failed, error: %w", err)
	}
	return c.String(http.StatusOK, fileToken)
}

func (e *EchoHandlers) GetBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookToken := strings.Replace(path, "/books/", "", 1)
	bookData, err := e.Services.GetBook(bookToken)
	if err != nil {
		return fmt.Errorf("attempting to receive book data from db failed, error: %w", err)
	}

	return c.JSON(http.StatusOK, bookData)
}

func getBooksParamsFieldFiller(c echo.Context, userEmail string) (*models.Filter, *models.Sort, error) {

	filter := models.Filter{
		Email:  userEmail,
		Search: c.QueryParams().Get("search"),
	}

	intLimit, err := strconv.Atoi(c.QueryParams().Get("limit"))
	if err != nil {
		return nil, nil, fmt.Errorf("limit param should be a number, error: %w", err)
	}
	var offset string
	if c.QueryParams().Get("offset") == "" {
		offset = "0"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, nil, fmt.Errorf("offset param should be a number")
	}

	sorting := models.Sort{
		SortField: c.QueryParams().Get("sort"),
		Limit:     intLimit,
		Direction: c.QueryParams().Get("direction"),
		Offset:    intOffset,
	}

	return &filter, &sorting, nil
}

func (e EchoHandlers) GetBooksPublic(c echo.Context) error {
	filter, sort, err := getBooksParamsFieldFiller(c, "")
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return fmt.Errorf("service layer get books public error: %w", err)
	}
	return c.JSON(http.StatusOK, books)
}

func (e EchoHandlers) GetBooksOfUser(c echo.Context) error {
	userId, err := getUserIdFromCtx(c)
	if err != nil {
		return err
	}
	userData, err := e.Services.GetUserByInsertedId(userId)
	if err != nil {
		return fmt.Errorf("getting user data by jwt user id failed, error: %w", err)
	}

	filter, sort, err := getBooksParamsFieldFiller(c, userData.Email)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	books, err := e.Services.GetBooks(*filter, *sort)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, books)
}

func (e *EchoHandlers) SignUp(c echo.Context) error {
	var signUpData models.UserDataInput
	err := c.Bind(&signUpData)
	if err != nil {
		return fmt.Errorf("sign up handler: %w", err)
	}

	jwtTok, err := e.Services.SignUp(signUpData)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	bearerAuth := fmt.Sprintf("bearer %s", jwtTok)

	c.Response().Header().Add("Authorization", bearerAuth)

	return c.String(http.StatusOK, fmt.Sprintf("Thanks for registration!\nYour auth token: %s", bearerAuth))
}

func (e *EchoHandlers) SignIn(c echo.Context) error {
	var signInData models.UserDataInput
	err := c.Bind(&signInData)
	if err != nil {
		return fmt.Errorf("sign in handler: %w", err)
	}
	token, err := e.Services.SignIn(signInData)
	if err != nil {
		return fmt.Errorf("failed in signin service, error: %w", err)
	}

	c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (e *EchoHandlers) UpdateBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookId := strings.Replace(path, "/books/", "", 1)
	var updater models.BookDataUpdater
	c.Bind(&updater)
	err := e.Services.UpdateBook(bookId, updater)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return c.String(http.StatusOK, "Book data successfully updated!")
}

func (e *EchoHandlers) DeleteBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookId := strings.Replace(path, "/books/", "", 1)
	err := e.Services.DeleteBook(bookId)
	if err != nil {
		return c.String(http.StatusNoContent, "Books with provided ID not founded")
	}

	return c.JSON(http.StatusOK, "")
}

func (e *EchoHandlers) TestAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt_package.JwtCustomClaims)
	return c.String(http.StatusOK, fmt.Sprintf("your userId: %s", claims.UserId))
}

func New(serviceLayer Service) *EchoHandlers {
	return &EchoHandlers{
		Services: serviceLayer,
	}
}
