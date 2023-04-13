package gofile

import (
	"crud-books/models"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeApikey = "ajskdhasSHdjashdjh@HJ#HJ$H#JH$(!Y#(HdsaD3413KAS42ND5K))"

func givenTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func handlerHelper(t *testing.T, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile(path, os.O_RDONLY, 0777)
		assert.NoError(t, err)
		defer file.Close()
		responseFromFile, err := io.ReadAll(file)
		assert.NoError(t, err)
		_, err = w.Write(responseFromFile)
		require.NoError(t, err)
	})
}

func getServerToUploadHandler(t *testing.T, url string) http.Handler {
	fileName := "./testsData/getServerToUploadResponse.json"

	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	require.NoError(t, err)
	defer file.Close()
	j, err := io.ReadAll(file)
	require.NoError(t, err)

	var v models.UploadServerSummary
	err = json.Unmarshal(j, &v)
	require.NoError(t, err)
	v.Data.Server = url

	r, err := json.Marshal(v)
	require.NoError(t, err)
	err = os.WriteFile(fileName, r, 0777)
	require.NoError(t, err)

	return handlerHelper(t, fileName)
}

func uploadFileServerHandler(t *testing.T) http.Handler {
	return handlerHelper(t, "./testsData/uploadFileServerResponse.json")
}

func Test_Service_UploadFile(t *testing.T) {
	s := New(fakeApikey, "b0rELG")
	uploadBody := []byte("Hello, World!")

	mockUploadServer := givenTestServer(uploadFileServerHandler(t))
	mockGetServer := givenTestServer(getServerToUploadHandler(t, mockUploadServer.URL))

	urlGetServer = mockGetServer.URL
	urlUploadServer = mockUploadServer.URL

	got, err := s.UploadFile(uploadBody, true)

	require.NoError(t, err)
	assert.Equal(t, "https://gofile.io/d/Z19n9a", got.DownloadPage)
}

func Test_DeleteFile(t *testing.T) {
	fileToken := "123"
	s := New(fakeApikey, "b0rELG")
	type jsonScheme struct {
		ContentsId string `json:"contentsId"`
		Token      string `json:"token"`
	}
	var bodyJson jsonScheme

	mockServ := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			errMsg := []byte("method doesn't founded")
			w.Write(errMsg)
		}
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &bodyJson)
		require.NoError(t, err)
		assert.Equal(t, fileToken, bodyJson.ContentsId)
		assert.Equal(t, fakeApikey, bodyJson.Token)
	}))

	urlDeleteFile = mockServ.URL
	err := s.DeleteFile(fileToken)
	require.NoError(t, err)
}
