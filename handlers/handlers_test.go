package handlers

import (
	"bytes"
	"crud-books/models"
	"crud-books/services"
	mock_services "crud-books/services/mocks"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userEmail    = "jjkfsjdkfj@gmail.com"
	userPassword = "4136425"
	hash         = "fjakskajdskjakdsjdkajdkjKJ@"
	token        = "jdaksjdkasjdkasjdkas"
)

func getBufJson(inp models.UserDataInput) *bytes.Buffer {
	body, _ := json.Marshal(&inp)
	return bytes.NewBuffer(body)
}

func TestEchoHandlers_SignUp(t *testing.T) {
	buf := getBufJson(models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	})

	serv := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := serv.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	dbMock := mock_services.NewMockDB(mockCtrl)
	hashMock := mock_services.NewMockHasher(mockCtrl)
	tokenerMock := mock_services.NewMockTokener(mockCtrl)

	handlers := EchoHandlers{
		Services: services.New(dbMock, tokenerMock, nil, hashMock),
	}
	hashMock.EXPECT().GetNewHash(userPassword).Return(hash, nil)
	dbMock.EXPECT().CreateUser(userEmail, hash).Return("userId", nil)
	tokenerMock.EXPECT().GenerateTokens("userId").Return(token, "", nil)

	err := handlers.SignUp(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, fmt.Sprintf("Thanks for registration!\nYour auth token: bearer %s", token), rec.Body.String())
}

func TestEchoHandlers_SignIn(t *testing.T) {
	buf := getBufJson(models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	})
	userCred := models.UserData{
		Id:           "5246231",
		Email:        userEmail,
		PasswordHash: hash,
		BooksIds:     []string{"513513", "131351"},
	}

	serv := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := serv.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	dbMock := mock_services.NewMockDB(mockCtrl)
	hashMock := mock_services.NewMockHasher(mockCtrl)
	tokenerMock := mock_services.NewMockTokener(mockCtrl)

	handlers := EchoHandlers{
		Services: services.New(dbMock, tokenerMock, nil, hashMock),
	}
	dbMock.EXPECT().GetUserData(userEmail).Return(&userCred, nil)
	hashMock.EXPECT().CompareHashWithPassword(userPassword, hash).Return(nil)
	tokenerMock.EXPECT().GenerateTokens(userCred.Id).Return(token, "", nil)

	err := handlers.SignIn(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, fmt.Sprintf("Bearer %s", token), rec.Header().Get("Authorization"))
}
