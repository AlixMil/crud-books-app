package services

import (
	"crud-books/models"
	"crud-books/pkg/hasher"
	"crud-books/storageService/gofile/gofile_responses"
	"fmt"
	"log"
)

type Tokener interface {
	GenerateToken(userId string) (string, error)
	// ParseToken(token string) (string, error)
}

//go:generate mockgen -source=services.go -destination=mocks/mock.go
type DB interface {
	CreateUser(email, passwordHash string) (string, error)
	CreateBook(title, description, fileToken, emailOwner string) (string, error)
	ChangeFieldOfBook(id, fieldName, fieldValue string) error
	UploadFileData(fileToken, downloadPage string) error
	GetUserData(email string) (*models.UserData, error)
	GetBook(bookToken string) (*models.BookData, error)
	GetFileData(fileToken string) (*models.FileData, error)
	DeleteBook(tokenBook string) error
	GetUserDataByInsertedId(userId string) (*models.UserData, error)
	GetListBooksPublic(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error)
	GetListBooksOfUser(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error)
}

type Storager interface {
	UploadFile(file []byte, isTest bool) (*gofile_responses.UploadFileReturn, error)
	DeleteFile(fileToken string) error
}

type Services struct {
	db      DB
	tokener Tokener
	storage Storager
	hasher  hasher.Hasher
}

type GetBookResponse struct {
	FileUrl     string `json:"fileUrl"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

const (
	sortFieldDefaultParam = "title"
	directionDefaultParam = 1
	limitDefaultParam     = 10
	maxSizeOfLimitParam   = 100
)

func getParamsWValidate(email, search, sortField, direction string, limit int) *models.ValidateDataInGetLists {
	var res models.ValidateDataInGetLists

	res.Email = email
	res.Search = search
	res.Limit = limit

	if sortField != "title" && sortField != "date" {
		res.SortField = sortFieldDefaultParam
	} else if sortField == "title" {
		res.SortField = "title"
	} else if sortField == "date" {
		res.SortField = "date"
	}

	if direction != "desc" && direction != "asc" {
		res.Direction = directionDefaultParam
	} else if direction == "asc" {
		res.Direction = 1
	} else if direction == "desc" {
		res.Direction = -1
	}

	if limit > maxSizeOfLimitParam {
		res.Limit = maxSizeOfLimitParam
	}

	if limit == 0 || limit < 0 {
		res.Limit = limitDefaultParam
	}

	return &res
}

func (s Services) SignIn(user models.UserDataInput) (string, error) {
	userCred, err := s.db.GetUserData(user.Email)
	if err != nil {
		return "", fmt.Errorf("failed find user in db in sign in method, error: %w", err)
	}
	err = s.hasher.CompareHashWithPassword(user.Password, userCred.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("password isn't equal")
	}
	token, err := s.tokener.GenerateToken(userCred.Id)
	if err != nil {
		return "", fmt.Errorf("failed generate token in sign in method, error: %w", err)
	}

	return token, nil
}

func (s Services) SignUp(user models.UserDataInput) (string, error) {
	passwordHash, err := s.hasher.GetNewHash(user.Password)
	if err != nil {
		return "", fmt.Errorf("getting new hash of password failed, error: %w", err)
	}

	userId, err := s.db.CreateUser(user.Email, passwordHash)
	if err != nil {
		return "", fmt.Errorf("create user DB proccess failed, error: %w", err)
	}

	token, err := s.tokener.GenerateToken(userId)
	if err != nil {
		return "", fmt.Errorf("generation token failed, error: %w", err)
	}
	log.Println("services signup")

	return token, nil
}

func (s Services) ParseToken(token string) (string, error) {
	// userId, err := s.tokener.ParseToken(token)
	// if err != nil {
	// 	return "", fmt.Errorf("failed parsing token, error: %w", err)
	// }

	return "userId", nil
}

func (s Services) CreateBook(title, description, fileToken, userEmail string) (string, error) {
	bookToken, err := s.db.CreateBook(title, description, fileToken, userEmail)
	if err != nil {
		return "", fmt.Errorf("create book failed, error: %w", err)
	}

	return bookToken, nil
}

func (s Services) UploadFile(file []byte) (string, error) {
	res, err := s.storage.UploadFile(file, false)
	if err != nil {
		return "", fmt.Errorf("upload file to service failed, error: %w", err)
	}
	err = s.db.UploadFileData(res.FileToken, res.DownloadPage)
	if err != nil {
		return "", fmt.Errorf("recording to db of uploaded file failed, error: %w", err)
	}
	return res.FileToken, nil
}

func (s Services) GetBook(bookToken string) (*GetBookResponse, error) {
	bookData, err := s.db.GetBook(bookToken)
	if err != nil {
		return nil, fmt.Errorf("get book in get book of service failed, error: %w", err)
	}
	return &GetBookResponse{
		FileUrl:     bookData.Url,
		Title:       bookData.Title,
		Description: bookData.Description,
	}, nil
}

func (s Services) GetBooksPublic(filter models.Filter, sorting models.Sort) (*[]models.BookData, error) {
	var books *[]models.BookData
	validateParams := getParamsWValidate("", filter.Search, sorting.SortField, sorting.Direction, sorting.Limit)
	if validateParams.Email == "" {
		_, err := s.db.GetListBooksPublic(validateParams)
		if err != nil {
			return nil, fmt.Errorf("get list of books public, error: %w", err)
		}
	}
	if validateParams.Email != "" {
		_, err := s.db.GetListBooksOfUser(validateParams)
		if err != nil {
			return nil, fmt.Errorf("get list of books error: %w", err)
		}
	}
	return books, nil
}

func (s Services) GetListBooksOfUser(filter models.Filter, sorting models.Sort) (*[]models.BookData, error) {
	validParams := getParamsWValidate(filter.Email, filter.Search, sorting.SortField, sorting.Direction, sorting.Limit)

	books, err := s.db.GetListBooksOfUser(validParams)
	if err != nil {
		return nil, fmt.Errorf("get list of books method failed, error: %w", err)
	}
	return books, nil
}

func (s Services) UpdateBook(bookField, tokenBook, fieldName, fieldValue string) error {
	err := s.db.ChangeFieldOfBook(tokenBook, fieldName, fieldValue)
	if err != nil {
		return fmt.Errorf("updating of fields book was failed, error: %w", err)
	}
	return nil
}

func (s Services) DeleteBook(tokenBook string) error {
	err := s.db.DeleteBook(tokenBook)
	if err != nil {
		return fmt.Errorf("delete book was failed, error: %w", err)
	}
	return nil
}

func New(db DB, tokener Tokener, storage Storager, hasher hasher.Hasher) *Services {
	return &Services{
		db:      db,
		tokener: tokener,
		storage: storage,
		hasher:  hasher,
	}
}
