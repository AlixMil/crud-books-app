package services

import (
	"crud-books/models"
	"crud-books/storageService/gofile/gofile_responses"
	"fmt"
)

//go:generate mockgen -source=services.go -destination=mocks/services_mock.go
type Tokener interface {
	GenerateTokens(userId string) (string, string, error)
}

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

type Hasher interface {
	GetNewHash(password string) (string, error)
	CompareHashWithPassword(password, hash string) error
}

type Services struct {
	db      DB
	tokener Tokener
	storage Storager
	hasher  Hasher
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

func getParamsWValidation(email, search, sortField, direction string, limit, offset int) *models.ValidateDataInGetLists {
	var res models.ValidateDataInGetLists

	res.Email = email
	res.Search = search
	res.Limit = limit
	res.Offset = offset

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
	token, _, err := s.tokener.GenerateTokens(userCred.Id)
	if err != nil {
		return "", fmt.Errorf("failed generate token in sign in method, error: %w", err)
	}

	return token, nil
}

func (s Services) GetUserByInsertedId(userId string) (*models.UserData, error) {
	usData, err := s.db.GetUserDataByInsertedId(userId)
	if err != nil {
		return nil, fmt.Errorf("userData by user id fetch failed, error: %w", err)
	}
	return usData, nil
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

	token, _, err := s.tokener.GenerateTokens(userId)
	if err != nil {
		return "", fmt.Errorf("generation token failed, error: %w", err)
	}

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
	validateParams := getParamsWValidation("", filter.Search, sorting.SortField, sorting.Direction, sorting.Limit, sorting.Offset)
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
	validParams := getParamsWValidation(filter.Email, filter.Search, sorting.SortField, sorting.Direction, sorting.Limit, sorting.Offset)

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

func New(db DB, tokener Tokener, storage Storager, hasher Hasher) *Services {
	return &Services{
		db:      db,
		tokener: tokener,
		storage: storage,
		hasher:  hasher,
	}
}
