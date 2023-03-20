package gofile

import (
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

	var v UploadServerSummary
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

func TestService_UploadFile_Success(t *testing.T) {
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
