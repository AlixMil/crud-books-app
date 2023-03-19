package gofile

import (
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

func getServerToUploadHandler(t *testing.T) http.Handler {
	return handlerHelper(t, "./testsData/getServerToUploadResponse.json")
}

func uploadFileServerHandler(t *testing.T) http.Handler {
	return handlerHelper(t, "./testsData/uploadFileServerResponse.json")
}

func TestService_UploadFile_Success(t *testing.T) {
	s := New(fakeApikey, "b0rELG")
	uploadBody := []byte("Hello, World!")

	mockUploadServer := givenTestServer(uploadFileServerHandler(t))
	mockGetServer := givenTestServer(getServerToUploadHandler(t))
	urlUploadServer = mockUploadServer.URL
	urlGetServer = mockGetServer.URL

	got, err := s.UploadFile(uploadBody, true)
	require.NoError(t, err)
	assert.Equal(t, "https://gofile.io/d/Z19n9a", got)
}
