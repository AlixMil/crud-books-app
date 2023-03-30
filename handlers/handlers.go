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
	"net/http"
	"strconv"
	"strings"

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
	GenerateToken(userId string) (string, error)
	ParseToken(token string) (string, error)
}

type UserDB interface {
	CreateUser(email, passwordHash string) error
	CreateBook(title, description, fileToken, emailOwner string) (string, error)
	ChangeFieldOfBook(collectionName, id, fieldName, fieldValue string) error
	UploadFileData(fileToken, downloadPage string) error
	GetUserData(email string) (*models.UserData, error)
	GetBook(bookToken string) (*models.BookData, error)
	GetFileData(fileToken string) (*models.FileData, error)
	GetListOfBooks(userEmail string) (*[]models.BookData, error)
	DeleteBook(tokenBook string) error
	GetUserDataByInsertedId(userId string) (*models.UserData, error)
}

type Service interface {
	SignIn(user models.UserDataInput) (string, error)
	SugnUp(user models.UserDataInput) error
	CreateBook(title, description, fileToken, userEmail string) (string, error)
	UploadFile(file []byte) (string, error)
	GetBook(bookToken string) (*services.GetBookResponse, error)
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
	// load file to storage -> return token of file
	// record to db token
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

func (e *EchoHandlers) CreateBook(c echo.Context) error {
	var reqBody CreateBookRequest
	err := c.Bind(&reqBody)
	if err != nil {
		return fmt.Errorf("failed of reading (binding) request body in create book func of handlers. Error: %w", err)
	}
	jwtToken := c.Request().Header.Get("Authorization")
	userId, err := e.Tokener.ParseToken(jwtToken)
	if err != nil {
		return fmt.Errorf("failed parse jwt, error: %w", err)
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

func (e EchoHandlers) GetUserBooks(c echo.Context) error {
	if c.QueryParams().Get("search") == "" {
		return c.String(http.StatusBadRequest, "search query parameter is mandatory for this request")
	}
	jwtToken := c.Request().Header.Get("Authorization")
	userId, err := e.Tokener.ParseToken(jwtToken)
	if err != nil {
		return fmt.Errorf("failed parse jwt, error: %w", err)
	}
	userData, err := e.db.GetUserDataByInsertedId(userId)
	if err != nil {
		return fmt.Errorf("getting user data by jwt user id failed, error: %w", err)
	}
	search := c.QueryParams().Get("search")
	limit := c.QueryParams().Get("limit")
	sort := c.QueryParams().Get("sort")
	direction := c.QueryParams().Get("direction")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.String(http.StatusBadRequest, "limit parameter should be number")
	}

	filter := models.Filter{
		Email:  userData.Email,
		Search: search,
	}
	sorting := models.Sort{
		SortField: sort,
		Limit:     limitInt,
		Direction: direction,
	}

	books, err := e.Services.GetListBooksOfUser(filter, sorting)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	c.JSON(http.StatusOK, books)

}

func (e *EchoHandlers) SignUp(c echo.Context) error {
	var signUpData models.UserDataInput
	err := c.Bind(&signUpData)
	if err != nil {
		return fmt.Errorf("sign up handler: %w", err)
	}
	e.Services.SugnUp(models.UserDataInput(signUpData))

	h, err := e.Hasher.GetNewHash(signUpData.Password)
	if err != nil {
		return fmt.Errorf("hashing of password failed in sign up handler, error: %w", err)
	}
	err = e.db.CreateUser(signUpData.Email, h)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Thanks for registration!")
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
		c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
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
