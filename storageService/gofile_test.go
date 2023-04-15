package storageService

import (
	"crud-books/config"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultConfig = config.Config{
	GoFileServiceApiKey: "123",
	GoFileFolderToken:   "123",
}

var testDeleteFileData = DeleteFileScheme{
	ContentsId: "15133513",
	Token:      defaultConfig.GoFileFolderToken,
}

func givenTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func handlerHelper[T UploadFileResponse | ServerToUploadResponse](t *testing.T, testData T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonByte, err := json.Marshal(testData)
		assert.NoError(t, err)
		_, err = w.Write(jsonByte)
		require.NoError(t, err)
	})
}

func getServerToUploadHandler(t *testing.T, url string) http.Handler {
	testData := ServerToUploadResponse{
		Status: "ok",
		Data: DataFromServerToUploadResponse{
			Server: url,
		},
	}
	return handlerHelper(t, testData)
}

func uploadFileServerHandler(t *testing.T) http.Handler {
	testData := UploadFileResponse{
		Status: "ok",
		Data: DataFromUploadFileResponse{
			DownloadPage: "https://gofile.io/d/Z19n9a",
			Code:         "Z19n9a",
			ParentFolder: "3dbc2f87-4c1e-4a81-badc-af004e61a5b4",
			FileID:       "4991e6d7-5217-46ae-af3d-c9174adae924",
			FileName:     "example.mp4",
			Md5:          "10c918b1d01aea85864ee65d9e0c2305",
		},
	}

	return handlerHelper(t, testData)
}

func deleteFileHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		reqBody := new(DeleteFileScheme)
		err = json.Unmarshal(body, &reqBody)
		require.NoError(t, err)
		assert.Equal(t, testDeleteFileData.ContentsId, reqBody.ContentsId)
		assert.Equal(t, testDeleteFileData.Token, reqBody.Token)
	})
	// handlerHelper(t, testData)
}

func Test_Service_GetServerToUpload(t *testing.T) {
	s := New(defaultConfig)

	const urlToUpload = "https://google.com"

	mockGetToUploadServ := givenTestServer(getServerToUploadHandler(t, urlToUpload))

	urlGetServer = mockGetToUploadServ.URL
	serv, err := s.GetServerToUpload()
	require.NoError(t, err)
	assert.Equal(t, urlToUpload, serv)
}

func Test_Service_UploadFile(t *testing.T) {
	s := New(defaultConfig)
	file := []byte("Hello, World!")

	mockUploadServer := givenTestServer(uploadFileServerHandler(t))

	got, err := s.UploadFile(mockUploadServer.URL, file)

	require.NoError(t, err)
	assert.Equal(t, "https://gofile.io/d/Z19n9a", got.DownloadPage)
}

func Test_DeleteFile(t *testing.T) {
	s := New(defaultConfig)

	mockServ := givenTestServer(deleteFileHandler(t))

	urlDeleteFile = mockServ.URL
	err := s.DeleteFile(testDeleteFileData.ContentsId)
	require.NoError(t, err)
}
