package services

import (
	"crud-books/models"
	"crud-books/mongodb"
	"crud-books/pkg/hasher"
	jwt_package "crud-books/pkg/jwt"
	"crud-books/storageService/gofile"
	"crud-books/storageService/gofile/gofile_responses"
	"fmt"
)

type Tokener interface {
	GenerateToken(userId string) (string, error)
	ParseToken(token string) (string, error)
}

type DB interface {
	CreateUser(email, passwordHash string) error
	CreateBook(title, description, fileToken, emailOwner string) (string, error)
	ChangeFieldOfBook(collectionName, id, fieldName, fieldValue string) error
	UploadFileData(fileToken, downloadPage string) error
	GetUserData(email string) (*models.UserData, error)
	GetBook(bookToken string) (*models.BookData, error)
	GetFileData(fileToken string) (*models.FileData, error)
	GetListOfBooks(userEmail string) (*[]models.BookData, error)
	DeleteBook(tokenBook string) error
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

type getBookResponse struct {
	FileUrl     string `json:"fileUrl"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func New(db *mongodb.MongoDB, tokener jwt_package.JwtTokener, storage gofile.Storage, hasher hasher.Hasher) *Services {
	return &Services{
		db:      db,
		tokener: tokener,
		storage: storage,
		hasher:  hasher,
	}
}

func (s Services) SignIn(user models.UserDataInput) (string, error) {
	userCred, err := s.db.GetUserData(user.Email)
	if err != nil {
		return "", fmt.Errorf("failed find user in db in sign in method, error: %w", err)
	}
	err = s.hasher.CompareHashWithPassword(user.Password, userCred.Password)
	if err != nil {
		return "", fmt.Errorf("password isn't equal")
	}
	token, err := s.tokener.GenerateToken(userCred.Id)
	if err != nil {
		return "", fmt.Errorf("failed generate token in sign in method, error: %w", err)
	}

	return token, nil
}

func (s Services) SugnUp(user models.UserDataInput) error {
	passwordHash, err := s.hasher.GetNewHash(user.Password)
	if err != nil {
		return fmt.Errorf("getting new hash of password failed, error: %w", err)
	}

	return s.db.CreateUser(user.Email, passwordHash)
}

func (s Services) ParseToken(token string) (string, error) {
	userId, err := s.tokener.ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("failed parsing token, error: %w", err)
	}

	return userId, nil
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

func (s Services) GetBook(bookToken string) (*getBookResponse, error) {
	bookData, err := s.db.GetBook(bookToken)
	if err != nil {
		return nil, fmt.Errorf("get book in get book of service failed, error: %w", err)
	}
	return &getBookResponse{
		FileUrl:     bookData.DownloadUrl,
		Title:       bookData.Title,
		Description: bookData.Description,
	}, nil
}

func (s Services) GetListBooksOfUser(userEmail string) (*[]models.BookData, error) {
	res, err := s.db.GetListOfBooks(userEmail)
	if err != nil {
		return nil, fmt.Errorf("get list of books failed, error: %w", err)
	}
	return res, nil
}

func (s Services) UpdateBook(bookField, tokenBook, fieldName, fieldValue string) error {
	err := s.db.ChangeFieldOfBook("books", tokenBook, fieldName, fieldValue)
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
