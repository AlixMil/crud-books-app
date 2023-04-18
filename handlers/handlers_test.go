package handlers

import (
	"bytes"
	mock_handlers "crud-books/handlers/mocks"
	"crud-books/models"
	jwt_package "crud-books/pkg/jwt"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"reflect"
	"testing"
	"time"

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

func getReqWithJson(path, method string, buf *bytes.Buffer) (httptest.ResponseRecorder, echo.Context) {
	serv := echo.New()
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := serv.NewContext(req, rec)
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt_package.JwtCustomClaims{
		UserId: defaultUserId,
	})
	c.Set("user", auth)
	return *rec, c
}

func getHeaderForFile(fileName string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			"file",
			fileName,
		))
	h.Set("Content-Type", "application/pdf")
	return h
}

func getReqWithFormFile(path, method string, file []byte, fileName string) (httptest.ResponseRecorder, echo.Context) {
	serv := echo.New()
	multipartCore := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartCore)

	h := getHeaderForFile(fileName)

	part, _ := writer.CreatePart(h)
	r := bytes.NewReader(file)
	io.Copy(part, r)

	req := httptest.NewRequest(method, path, multipartCore)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	writer.Close()

	return *rec, serv.NewContext(req, rec)
}

func Test_EchoHandlers_SignUp(t *testing.T) {
	type SignUpResponse struct {
		Token string `json:"token"`
	}
	inp := models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	}
	jwtTok := "askdlaksdlasd"
	body, _ := json.Marshal(&inp)
	buf := bytes.NewBuffer(body)
	rec, c := getReqWithJson("/register", http.MethodPost, buf)
	mocks := getMocks(t)
	mocks.serviceLayer.EXPECT().SignUp(inp).Return(jwtTok, nil)

	h := New(mocks.serviceLayer)
	err := h.SignUp(c)
	require.NoError(t, err)

	res := new(SignUpResponse)

	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, jwtTok, res.Token)
}

func TestEchoHandlers_SignIn(t *testing.T) {
	inp := models.UserDataInput{
		Email:    userEmail,
		Password: userPassword,
	}
	body, _ := json.Marshal(&inp)
	buf := bytes.NewBuffer(body)

	rec, c := getReqWithJson("/login", http.MethodPost, buf)

	mocks := getMocks(t)
	jwtToken := "41351adas"

	h := New(mocks.serviceLayer)
	mocks.serviceLayer.EXPECT().SignIn(inp).Return(jwtToken, nil)

	err := h.SignIn(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, fmt.Sprintf("Bearer %s", jwtToken), rec.Header().Get("Authorization"))
}

func Test_UploadFile(t *testing.T) {
	mocks := getMocks(t)
	fileToken := "4819uikjsdkas"

	file, _ := os.Open("./testsData/file.pdf")
	defer file.Close()

	bytesOfFile, _ := io.ReadAll(file)

	rec, c := getReqWithFormFile("/files", http.MethodPost, bytesOfFile, file.Name())

	fileFromForm, fileHeader, _ := c.Request().FormFile("file")
	bytesOfFileForm, _ := io.ReadAll(fileFromForm)
	mocks.serviceLayer.EXPECT().UploadFile(bytesOfFileForm, fileHeader.Filename).Return(fileToken, nil)

	h := New(mocks.serviceLayer)

	err := h.UploadFile(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_GetUserIdFromCtx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt_package.JwtCustomClaims{
			UserId: defaultUserId,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			},
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
		assert.EqualError(t, err, "user context is null")
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
	}
	buf := bytes.NewBuffer(nil)
	b, _ := json.Marshal(&reqBody)
	buf.Write(b)
	rec, c := getReqWithJson("/books", http.MethodPost, buf)
	h := New(mocks.serviceLayer)

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
	rec, c := getReqWithJson(fmt.Sprintf("/books/%s", bookId), http.MethodGet, &bytes.Buffer{})
	mocks := getMocks(t)

	mocks.serviceLayer.EXPECT().GetBook(bookId).Return(&bookResp, nil)

	h := New(mocks.serviceLayer)

	err := h.GetBook(c)
	require.NoError(t, err)
	var res *models.GetBookResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, bookResp, *res)
}

func Test_GetBooksParamsFieldFiller(t *testing.T) {
	wantFilt := models.Filter{
		Email:  userEmail,
		Search: "asdjkasjdkajs",
	}
	wantSort := models.Sort{
		SortField: "title",
		Limit:     5,
		Direction: "desc",
		Offset:    0,
	}
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books?search=%s&limit=%d&sort=%s&direction=%s", wantFilt.Search, wantSort.Limit, wantSort.SortField, wantSort.Direction), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	filter, sorting, err := getBooksParamsFieldFiller(c, userEmail)
	require.NoError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(&wantFilt, filter))
	assert.Equal(t, true, reflect.DeepEqual(&wantSort, sorting))
}

func getReqForGetBooksParams(wantFilt models.Filter, wantSort models.Sort, isAuth bool) (*httptest.ResponseRecorder, *echo.Context) {
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/books?search=%s&limit=%d&sort=%s&direction=%s",
			wantFilt.Search,
			wantSort.Limit,
			wantSort.SortField,
			wantSort.Direction,
		),
		nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	if isAuth {
		auth := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt_package.JwtCustomClaims{
			UserId: defaultUserId,
		})
		c.Set("user", auth)
	}

	return rec, &c
}

func Test_GetBooksPublic(t *testing.T) {
	wantFilt := models.Filter{
		Email:  "",
		Search: "asdjkasjdkajs",
	}
	wantSort := models.Sort{
		SortField: "title",
		Limit:     5,
		Direction: "desc",
		Offset:    0,
	}
	wantBooks := []models.BookData{
		{
			Id:          "1",
			Title:       "Title1",
			Description: "description1",
			FileToken:   "FileToken1",
			Url:         "URL1",
			OwnerEmail:  "email1",
		},
		{
			Id:          "2",
			Title:       "Title2",
			Description: "description2",
			FileToken:   "FileToken2",
			Url:         "URL2",
			OwnerEmail:  "email2",
		},
	}

	rec, c := getReqForGetBooksParams(wantFilt, wantSort, false)

	mocks := getMocks(t)
	mocks.serviceLayer.EXPECT().GetBooks(wantFilt, wantSort).Return(&wantBooks, nil)
	h := New(mocks.serviceLayer)
	err := h.GetBooks(*c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var received []models.BookData
	json.Unmarshal(rec.Body.Bytes(), &received)
	assert.Equal(t, true, reflect.DeepEqual(wantBooks, received))
}

func TestGetBooksOfUser(t *testing.T) {
	userData := models.UserData{
		Id:           defaultUserId,
		Email:        "keer@gmail.com",
		PasswordHash: "581937ikajsdkajsf",
	}
	mocks := getMocks(t)
	wantFilt := models.Filter{
		Email:  userData.Email,
		Search: "sosage",
	}
	wantSort := models.Sort{
		SortField: "title",
		Limit:     5,
		Direction: "desc",
		Offset:    0,
	}

	booksResponse := []models.BookData{
		{
			Id:          "123",
			Title:       "title1",
			Description: "desc",
			FileToken:   "1241",
			Url:         "url",
			OwnerEmail:  "joajsd@gmail.com",
		},
		{
			Id:          "1233",
			Title:       "title14",
			Description: "desc4",
			FileToken:   "12414",
			Url:         "url4",
			OwnerEmail:  "joajsd@gmail.com4",
		},
	}
	rec, c := getReqForGetBooksParams(wantFilt, wantSort, true)
	h := New(mocks.serviceLayer)
	mocks.serviceLayer.EXPECT().GetUserByInsertedId(defaultUserId).Return(&userData, nil)
	mocks.serviceLayer.EXPECT().GetBooks(wantFilt, wantSort).Return(&booksResponse, nil)

	err := h.GetBooks(*c)
	require.NoError(t, err)
	var res *[]models.BookData

	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, true, reflect.DeepEqual(booksResponse, *res))

}

func Test_SignUp(t *testing.T) {
	usInp := models.UserDataInput{
		Email:    userEmail,
		Password: "153513",
	}
	mocks := getMocks(t)
	jwtTok := "ksaljlsajl4252"

	rec := httptest.NewRecorder()
	body, _ := json.Marshal(usInp)
	buf := bytes.NewBuffer(body)
	req := httptest.NewRequest(http.MethodPost, "/register", buf)
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echo.New().NewContext(req, rec)

	h := New(mocks.serviceLayer)
	mocks.serviceLayer.EXPECT().SignUp(usInp).Return(jwtTok, nil)

	err := h.SignUp(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	authHead := rec.Header().Get("Authorization")
	assert.Equal(t, fmt.Sprintf("Bearer %s", jwtTok), authHead)
}

func Test_SignIn(t *testing.T) {
	mocks := getMocks(t)
	usInp := models.UserDataInput{
		Email:    userEmail,
		Password: "153513",
	}
	jwtTok := "ksaljlsajl4252"

	rec := httptest.NewRecorder()
	body, _ := json.Marshal(usInp)
	buf := bytes.NewBuffer(body)
	req := httptest.NewRequest(http.MethodPost, "/login", buf)
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echo.New().NewContext(req, rec)

	mocks.serviceLayer.EXPECT().SignIn(usInp).Return(jwtTok, nil)

	h := New(mocks.serviceLayer)

	err := h.SignIn(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	authHead := rec.Header().Get("Authorization")
	assert.Equal(t, fmt.Sprintf("Bearer %s", jwtTok), authHead)
}

func Test_DeleteBook(t *testing.T) {
	mocks := getMocks(t)
	bookId := "123"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/books/%s", bookId), &bytes.Buffer{})
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echo.New().NewContext(req, rec)

	mocks.serviceLayer.EXPECT().DeleteBook(bookId).Return(nil)

	h := New(mocks.serviceLayer)
	err := h.DeleteBook(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_UpdateBook(t *testing.T) {
	mocks := getMocks(t)
	bookId := "123"
	updater := models.BookDataUpdater{
		FileToken:   "filetoken",
		Title:       "Title",
		Description: "description",
	}

	body, _ := json.Marshal(updater)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/books/%s", bookId), bytes.NewBuffer(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := echo.New().NewContext(req, rec)

	mocks.serviceLayer.EXPECT().UpdateBook(bookId, updater).Return(nil)

	h := New(mocks.serviceLayer)
	err := h.UpdateBook(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, "Book data successfully updated!", rec.Body.String())
}
