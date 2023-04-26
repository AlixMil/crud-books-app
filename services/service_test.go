package services

import (
	"crud-books/models"
	mock_services "crud-books/services/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	db       *mock_services.MockDB
	hasher   *mock_services.MockHasher
	tokener  *mock_services.MockTokener
	storager *mock_services.MockStorager
}

func getMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		db:       mock_services.NewMockDB(ctrl),
		hasher:   mock_services.NewMockHasher(ctrl),
		tokener:  mock_services.NewMockTokener(ctrl),
		storager: mock_services.NewMockStorager(ctrl),
	}
}

func Test_SignIn(t *testing.T) {
	mocks := getMocks(t)
	const token = "alsdkalsdkjl22412"
	userInp := models.UserDataInput{
		Email:    "jojek.as@gmail.com",
		Password: "41uijaksjd",
	}
	userData := models.UserData{
		Id:           "12",
		Email:        userInp.Email,
		PasswordHash: "kasldjlajsd@2413",
	}

	mocks.db.EXPECT().GetUserData(userInp.Email).Return(&userData, nil)
	mocks.hasher.EXPECT().CompareHashWithPassword(userInp.Password, userData.PasswordHash).Return(nil)
	mocks.tokener.EXPECT().GenAccessToken(userData.Id).Return(token, nil)

	s := New(mocks.db, mocks.tokener, nil, mocks.hasher)
	tokenRes, err := s.SignIn(userInp)

	require.NoError(t, err)
	assert.Equal(t, token, tokenRes)
}

func Test_SignUp(t *testing.T) {
	mocks := getMocks(t)
	const (
		hash   = "jfkahsjkdhakjbjahsdk@&%*#Y*y58h"
		token  = "alsdkalsdkjl22412"
		userId = "4513513"
	)
	userInp := models.UserDataInput{
		Email:    "jojek.as@gmail.com",
		Password: "41uijaksjd",
	}

	mocks.hasher.EXPECT().GetNewHash(userInp.Password).Return(hash, nil)
	mocks.db.EXPECT().CreateUser(userInp.Email, hash).Return(userId, nil)
	mocks.tokener.EXPECT().GenAccessToken(userId).Return(token, nil)

	s := New(mocks.db, mocks.tokener, nil, mocks.hasher)
	tokenRes, err := s.SignUp(userInp)

	require.NoError(t, err)
	assert.Equal(t, token, tokenRes)
}

func Test_GetUserByInsertedId(t *testing.T) {
	mocks := getMocks(t)
	const userId = "51324124"
	userData := models.UserData{
		Id:           userId,
		Email:        "jojekalsjd@gmail.com",
		PasswordHash: "fjakjsdk2412",
	}

	mocks.db.EXPECT().GetUserDataById(userId).Return(&userData, nil)
	s := New(mocks.db, nil, nil, nil)
	usData, err := s.GetUserById(userId)
	require.NoError(t, err)
	assert.Equal(t, &userData, usData)
}

func Test_CreateBook(t *testing.T) {
	mocks := getMocks(t)
	book := models.BookData{
		Id:          "142",
		Title:       "bookTitle",
		Description: "descriptionBook",
		FileToken:   "41213513513614",
		Url:         "https://google.com",
		OwnerEmail:  "jaksdjkasd@gmail.com",
	}
	mocks.db.EXPECT().CreateBook(book.Title, book.Description, book.FileToken, book.OwnerEmail).Return(book.FileToken, nil)

	s := New(mocks.db, nil, nil, nil)
	fToken, err := s.CreateBook(book.Title, book.Description, book.FileToken, book.OwnerEmail)
	require.NoError(t, err)
	assert.Equal(t, book.FileToken, fToken)
}

func Test_UploadFile(t *testing.T) {
	mocks := getMocks(t)
	defServToUpload := "google.com"
	file := []byte("dkalsjdkjasdkas")

	fileName := "namefile"

	fileRet := models.FileData{
		DownloadPage: "http://download.com",
		Token:        "jfkajsdkj413513",
	}
	mocks.storager.EXPECT().GetServerToUpload().Return(defServToUpload, nil)
	mocks.storager.EXPECT().UploadFile(defServToUpload, file, fileName).Return(&fileRet, nil)
	mocks.db.EXPECT().UploadFileData(&fileRet).Return(nil)

	s := New(mocks.db, nil, mocks.storager, nil)
	res, err := s.UploadFile(file, fileName)
	require.NoError(t, err)
	assert.Equal(t, fileRet.Token, res)
}

func Test_GetBook(t *testing.T) {
	mocks := getMocks(t)
	const bookToken = "1893859195"
	bookData := models.BookData{
		Id:          "413313",
		Title:       "Title",
		Description: "descriptio",
		FileToken:   bookToken,
		Url:         "aksldkasl",
		OwnerEmail:  "kfaljsdlj@gmail.com",
	}
	want := &models.GetBookResponse{
		FileURL:     bookData.Url,
		Title:       bookData.Title,
		Description: bookData.Description,
	}
	mocks.db.EXPECT().GetBook(bookToken).Return(&bookData, nil)

	s := New(mocks.db, nil, nil, nil)
	res, err := s.GetBook(bookToken)
	require.NoError(t, err)
	assert.Equal(t, want, res)
}

func Test_GetBooks(t *testing.T) {
	mocks := getMocks(t)
	filter := models.Filter{
		Email:  "jaksjdksjad@gmail.com",
		Search: "skdjsk",
	}
	sort := models.Sort{
		SortField: "title",
		Limit:     10,
		Direction: "asc",
		Offset:    10,
	}
	want := []models.BookData{
		{
			Id:          "id",
			Title:       "title",
			Description: "descript",
			FileToken:   "file token",
			Url:         "URL",
			OwnerEmail:  "owner",
		},
		{
			Id:          "id2",
			Title:       "title2",
			Description: "descript2",
			FileToken:   "file token2",
			Url:         "URL2",
			OwnerEmail:  "owner2",
		},
	}
	mocks.db.EXPECT().GetListBooksUser(filter, sort).Return(want, nil)

	s := New(mocks.db, nil, nil, nil)
	books, err := s.GetBooks(filter, sort)
	require.NoError(t, err)
	assert.Equal(t, &want, books)
}
