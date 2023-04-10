package handlers

import (
	"bytes"
	mock_handlers "crud-books/handlers/mocks"
	"crud-books/models"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userEmail        = "jjkfsjdkfj@gmail.com"
	userPassword     = "4136425"
	hash             = "fjakskajdskjakdsjdkajdkjKJ@"
	defaultUserToken = "51351jkjfkdjakJKJJUAOSd#A"
	defaultUserId    = "4131351351351523"
)

type mocks struct {
	serviceLayer *mock_handlers.MockService
}

func getMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		serviceLayer: mock_handlers.NewMockService(ctrl),
	}
}

func getRequestWRecJson(path, method string, buf *bytes.Buffer) (httptest.ResponseRecorder, echo.Context) {
	serv := echo.New()
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := serv.NewContext(req, rec)
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": defaultUserId,
	})
	c.Set("user", auth)
	return *rec, c
}

func getRequestWRecFormFile(path, method string, file []byte) (httptest.ResponseRecorder, echo.Context) {
	serv := echo.New()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "filename")
	r := bytes.NewReader(file)
	io.Copy(part, r)

	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	writer.Close()

	return *rec, serv.NewContext(req, rec)
}

func Test_EchoHandlers_SignUp(t *testing.T) {
	inp := models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	}
	jwtTok := "askdlaksdlasd"
	body, _ := json.Marshal(&inp)
	buf := bytes.NewBuffer(body)
	rec, c := getRequestWRecJson("/register", http.MethodPost, buf)
	mocks := getMocks(t)
	mocks.serviceLayer.EXPECT().SignUp(inp).Return(jwtTok, nil)

	h, _ := New(mocks.serviceLayer)
	err := h.SignUp(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, fmt.Sprintf("Thanks for registration!\nYour auth token: bearer %s", jwtTok), rec.Body.String())
}

func TestEchoHandlers_SignIn(t *testing.T) {
	inp := models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	}
	body, _ := json.Marshal(&inp)
	buf := bytes.NewBuffer(body)

	rec, c := getRequestWRecJson("/login", http.MethodPost, buf)

	mocks := getMocks(t)
	jwtToken := "41351adas"

	h, _ := New(mocks.serviceLayer)
	mocks.serviceLayer.EXPECT().SignIn(inp).Return(jwtToken, nil)

	err := h.SignIn(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, fmt.Sprintf("Bearer %s", jwtToken), rec.Header().Get("Authorization"))
}

func Test_UploadFile(t *testing.T) {
	mocks := getMocks(t)
	fileToken := "4819uikjsdkas"

	file, err := os.Open("./testsData/file.pdf")
	require.NoError(t, err)
	defer file.Close()
	fileByte, err := io.ReadAll(file)
	require.NoError(t, err)
	rec, c := getRequestWRecFormFile("/files", http.MethodPost, fileByte)
	mocks.serviceLayer.EXPECT().UploadFile(fileByte).Return(fileToken, nil)

	h, _ := New(mocks.serviceLayer)

	err = h.UploadFile(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ := io.ReadAll(rec.Body)
	assert.Equal(t, fileToken, string(res))
}

func Test_GetUserIdFromCtx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": defaultUserId,
		})
		c := srv.NewContext(req, rec)
		c.Set("user", tok)

		usId, err := getUserIdFromCtx(c)
		require.NoError(t, err)
		assert.Equal(t, defaultUserId, usId)
	})

	t.Run("failed", func(t *testing.T) {
		srv := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := srv.NewContext(req, rec)

		_, err := getUserIdFromCtx(c)
		assert.EqualError(t, err, "get user id from ctx failed, jwt token missed")
	})

}

func Test_CreateBook(t *testing.T) {
	mocks := getMocks(t)
	reqBody := models.CreateBookRequest{
		FileToken:   "513513513",
		Title:       "TITLE",
		Description: "description",
	}
	userData := models.UserData{
		Id:           defaultUserId,
		Email:        "jfkaljsdj@gmail.com",
		PasswordHash: "kfaksjdkasd2415",
		BooksIds:     []string{},
	}
	buf := bytes.NewBuffer(nil)
	b, _ := json.Marshal(&reqBody)
	buf.Write(b)
	rec, c := getRequestWRecJson("/books", http.MethodPost, buf)
	h, _ := New(mocks.serviceLayer)

	mocks.serviceLayer.EXPECT().GetUserByInsertedId(defaultUserId).Return(&userData, nil)
	mocks.serviceLayer.EXPECT().CreateBook(reqBody.Title, reqBody.Description, reqBody.FileToken, userData.Email).Return(reqBody.FileToken, nil)

	err := h.CreateBook(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	assert.Equal(t, reqBody.FileToken, string(res))
}

func Test_GetBook(t *testing.T) {
	bookId := "1513"
	bookResp := models.GetBookResponse{
		Title:       "Story of islands",
		Description: "Terrible island stories",
		FileURL:     "google.com",
	}
	rec, c := getRequestWRecJson(fmt.Sprintf("/books/%s", bookId), http.MethodGet, &bytes.Buffer{})
	mocks := getMocks(t)

	mocks.serviceLayer.EXPECT().GetBook(bookId).Return(&bookResp, nil)

	h, _ := New(mocks.serviceLayer)

	err := h.GetBook(c)
	require.NoError(t, err)
	var res *models.GetBookResponse
	fmt.Println(rec.Body)
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, bookResp, *res)
}
