package handlers

import (
	"bytes"
	"crud-books/models"
	"crud-books/pkg/hasher"
	"crud-books/services"
	"crud-books/storageService/gofile/gofile_responses"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type EchoHandlers struct {
	db             UserDB
	ServiceStorage ServiceStorage
	Hasher         hasher.Hasher
	JWTSecret      string
	Services       Service
	Tokener        tokener
}

type tokener interface {
	GenerateTokens(userId string) (string, string, error)
}

type jwtCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type UserDB interface {
	GetBook(bookToken string) (*models.BookData, error)
	CreateUser(email, passwordHash string) (string, error)
	CreateBook(title, description, fileToken, emailOwner string) (string, error)
	GetListBooksPublic(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error)
	GetListBooksOfUser(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error)
	ChangeFieldOfBook(id, fieldName, fieldValue string) error
	UploadFileData(fileToken, downloadPage string) error
	GetFileData(fileToken string) (*models.FileData, error)
	GetUserData(email string) (*models.UserData, error)
	GetUserDataByInsertedId(userId string) (*models.UserData, error)
	DeleteBook(tokenBook string) error
}

type Service interface {
	SignIn(user models.UserDataInput) (string, error)
	SignUp(user models.UserDataInput) (string, error)
	CreateBook(title, description, fileToken, userEmail string) (string, error)
	UploadFile(file []byte) (string, error)
	GetBook(bookToken string) (*services.GetBookResponse, error)
	GetBooksPublic(filter models.Filter, sorting models.Sort) (*[]models.BookData, error)
	GetListBooksOfUser(filter models.Filter, sorting models.Sort) (*[]models.BookData, error)
	UpdateBook(bookField, tokenBook, fieldName, fieldValue string) error
	DeleteBook(tokenBook string) error
}

type ServiceStorage interface {
	UploadFile(file []byte, isTest bool) (*gofile_responses.UploadFileReturn, error)
	DeleteFile(fileToken string) error
}

type UploadFormData struct {
	File []byte `form:"file"`
}

type FormFileData struct {
	ContentType string
	Content     []byte
	Name        string
}

type GetBookResponse struct {
	FileURL     string `json:"fileURL"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (e *EchoHandlers) UploadFile(c echo.Context) error {
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

	fileToken, err := e.Services.UploadFile(buf.Bytes())
	if err != nil {
		return fmt.Errorf("upload file failed, error: %w", err)
	}

	return c.String(http.StatusOK, fileToken)
}

type CreateBookRequest struct {
	FileToken   string `json:"fileToken"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getUserId(context echo.Context) (string, error) {
	user := context.Get("user").(*jwt.Token)
	if user == nil {
		return "", fmt.Errorf("getting user from context in getUserId failed")
	}
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)
	if userId == "" {
		return "", fmt.Errorf("userId is empty")
	}
	return userId, nil
}

func (e *EchoHandlers) CreateBook(c echo.Context) error {
	var reqBody CreateBookRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return fmt.Errorf("failed of reading (binding) request body in create book func of handlers. Error: %w", err)
	}
	userId, err := getUserId(c)
	if err != nil {
		return err
	}
	userData, err := e.db.GetUserDataByInsertedId(userId)
	if err != nil {
		return fmt.Errorf("getting user data by jwt user id failed, error: %w", err)
	}

	bookToken, err := e.Services.CreateBook(reqBody.Title, reqBody.Description, reqBody.FileToken, userData.Email)
	if err != nil {
		return fmt.Errorf("create book failed, error: %w", err)
	}
	return c.String(http.StatusOK, bookToken)
}

func (e *EchoHandlers) GetBook(c echo.Context) error {
	path := c.Request().URL.Path
	bookToken := strings.Replace(path, "/books/", "", 1)
	bookData, err := e.Services.GetBook(bookToken)
	if err != nil {
		return fmt.Errorf("attempt to receive book data from db failed, error: %w", err)
	}
	response := GetBookResponse{
		FileURL:     bookData.FileUrl,
		Title:       bookData.Title,
		Description: bookData.Description,
	}
	marshResp, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshaling after get book failed, error: %w", err)
	}
	return c.JSON(http.StatusOK, marshResp)
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
	books, err := e.Services.GetBooksPublic(*filter, *sort)
	if err != nil {
		return fmt.Errorf("service layer get books public error: %w", err)
	}
	return c.JSON(http.StatusOK, books)
}

func (e EchoHandlers) GetBooksOfUser(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}
	userData, err := e.db.GetUserDataByInsertedId(userId)
	if err != nil {
		return fmt.Errorf("getting user data by jwt user id failed, error: %w", err)
	}

	filter, sort, err := getBooksParamsFieldFiller(c, userData.Email)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	books, err := e.Services.GetListBooksOfUser(*filter, *sort)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, books)
}

func (e *EchoHandlers) SignUp(c echo.Context) error {
	var signUpData models.UserDataInput
	err := c.Bind(&signUpData)
	log.Println("sign up handler")
	if err != nil {
		return fmt.Errorf("sign up handler: %w", err)
	}

	jwtTok, err := e.Services.SignUp(models.UserDataInput(signUpData))
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
	token, err := e.Services.SignIn(models.UserDataInput(signInData))
	if err != nil {
		return fmt.Errorf("failed in signin service, error: %w", err)
	}

	if c.Request().Header.Get("Authorization") == "" {
		c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (e *EchoHandlers) UpdateBook(c echo.Context) error {

	return c.JSON(http.StatusOK, "sda")
}

func (e *EchoHandlers) DeleteBook(c echo.Context) error {

	return c.JSON(http.StatusOK, "sda")
}

func (e *EchoHandlers) TestAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return c.String(http.StatusOK, fmt.Sprintf("your userId: %s", claims.UserId))
}

func New(db UserDB, storage ServiceStorage, jwtSecret string, serviceLayer Service, Tokener tokener) (*EchoHandlers, error) {
	return &EchoHandlers{
		db:        db,
		Hasher:    *hasher.New(),
		JWTSecret: jwtSecret,
		Services:  serviceLayer,
		Tokener:   Tokener,
	}, nil
}
