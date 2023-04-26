package services

import (
	"crud-books/models"
	"fmt"
)

//go:generate mockgen -source=services.go -destination=mocks/services_mock.go
type Tokener interface {
	GenAccessToken(userId string) (string, error)
	GenRefreshToken() (string, error)
}

type DB interface {
	CreateUser(email, passwordHash string) (string, error)
	GetUserData(email string) (*models.UserData, error)
	GetUserDataById(userId string) (*models.UserData, error)

	UploadFileData(fileData *models.FileData) error
	GetFileData(fileToken string) (*models.FileData, error)

	GetListBooks(filter models.Filter, sort models.Sort) ([]models.BookData, error)
	GetListBooksUser(filter models.Filter, sort models.Sort) ([]models.BookData, error)

	CreateBook(title, description, fileToken, emailOwner string) (string, error)
	GetBook(bookToken string) (*models.BookData, error)
	UpdateBook(bookFileToken string, updater models.BookDataUpdater) error
	DeleteBook(tokenBook string) error
}

type Storager interface {
	GetServerToUpload() (string, error)
	UploadFile(servForUpload string, file []byte, fileName string) (*models.FileData, error)
	DeleteFile(fileToken string) error
}

type Hasher interface {
	GetNewHash(password string) (string, error)
	CompareHashWithPassword(password, hash string) error
}

type Services struct {
	db          DB
	tokenEngine Tokener
	hashEngine  Hasher
	storage     Storager
}

func New(db DB, tokener Tokener, storage Storager, hasher Hasher) *Services {
	return &Services{
		db:          db,
		tokenEngine: tokener,
		hashEngine:  hasher,
		storage:     storage,
	}
}

func (s *Services) SignIn(user models.UserDataInput) (string, error) {
	userCred, err := s.db.GetUserData(user.Email)
	if err != nil {
		return "", fmt.Errorf("failed find user in db in sign in method, error: %w", err)
	}
	err = s.hashEngine.CompareHashWithPassword(user.Password, userCred.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("password isn't equal")
	}
	accessToken, err := s.tokenEngine.GenAccessToken(userCred.Id)
	if err != nil {
		return "", fmt.Errorf("failed generate token in sign in method, error: %w", err)
	}

	return accessToken, nil
}

func (s *Services) SignUp(user models.UserDataInput) (string, error) {
	passwordHash, err := s.hashEngine.GetNewHash(user.Password)
	if err != nil {
		return "", fmt.Errorf("getting new hash of password failed, error: %w", err)
	}

	userId, err := s.db.CreateUser(user.Email, passwordHash)
	if err != nil {
		return "", fmt.Errorf("create user DB proccess failed, error: %w", err)
	}

	accessToken, err := s.tokenEngine.GenAccessToken(userId)
	if err != nil {
		return "", fmt.Errorf("generation token failed, error: %w", err)
	}

	return accessToken, nil
}

func (s *Services) GetUserById(userId string) (*models.UserData, error) {
	usData, err := s.db.GetUserDataById(userId)
	if err != nil {
		return nil, fmt.Errorf("userData by user id fetch failed, error: %w", err)
	}
	return usData, nil
}

func (s *Services) CreateBook(title, description, fileToken, userEmail string) (string, error) {
	fileToken, err := s.db.CreateBook(title, description, fileToken, userEmail)
	if err != nil {
		return "", fmt.Errorf("create book failed, error: %w", err)
	}

	return fileToken, nil
}

func (s *Services) UploadFile(file []byte, fileName string) (string, error) {
	servForUpload, err := s.storage.GetServerToUpload()
	if err != nil {
		return "", fmt.Errorf("getting cell server for upload file failed, error: %w", err)
	}

	fileData, err := s.storage.UploadFile(servForUpload, file, fileName)
	if err != nil {
		return "", fmt.Errorf("upload file to service failed, error: \n%w", err)
	}

	err = s.db.UploadFileData(fileData)
	if err != nil {
		return "", fmt.Errorf("recording to db of uploaded file failed, error: %w", err)
	}
	return fileData.Token, nil
}

func (s *Services) GetBook(bookToken string) (*models.GetBookResponse, error) {
	bookData, err := s.db.GetBook(bookToken)
	if err != nil {
		return nil, fmt.Errorf("get book in get book of service failed, error: %w", err)
	}
	return &models.GetBookResponse{
		FileURL:     bookData.Url,
		Title:       bookData.Title,
		Description: bookData.Description,
	}, nil
}

func (s *Services) GetBooks(filter models.Filter, sorting models.Sort) (*[]models.BookData, error) {
	var books []models.BookData

	if filter.Email == "" {
		booksArr, err := s.db.GetListBooks(filter, sorting)
		if err != nil {
			return nil, fmt.Errorf("get list of books public, error: %w", err)
		}
		books = booksArr
	}
	if filter.Email != "" {
		booksArr, err := s.db.GetListBooksUser(filter, sorting)
		if err != nil {
			return nil, fmt.Errorf("get list of books error: %w", err)
		}
		books = booksArr
	}
	return &books, nil
}

func (s *Services) UpdateBook(bookFileToken string, updater models.BookDataUpdater) error {
	err := s.db.UpdateBook(bookFileToken, updater)
	if err != nil {
		return fmt.Errorf("updating failed, error: %w", err)
	}
	return nil
}

func (s *Services) DeleteBook(bookId string) error {
	err := s.db.DeleteBook(bookId)
	if err != nil {
		return fmt.Errorf("delete book was failed, error: %w", err)
	}
	return nil
}
