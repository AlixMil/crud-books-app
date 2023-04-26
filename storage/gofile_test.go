package storage

import (
	"crud-books/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultConfig = config.Config{
	GoFileServiceApiKey: "123",
	GoFileFolderToken:   "123",
}

var testDeleteFileData = DeleteFileRequest{
	ContentsId: "15133513",
	Token:      defaultConfig.GoFileFolderToken,
}

func getTestServer(handler http.Handler) *httptest.Server {
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

func Test_Service_GetServerToUpload(t *testing.T) {
	s := New(&defaultConfig)

	const urlToUpload = "store3"

	mockGetToUploadServ := getTestServer(getServerToUploadHandler(t, urlToUpload))

	urlGetServer = mockGetToUploadServ.URL
	serv, err := s.GetServerToUpload()
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("https://%s.gofile.io/uploadFile", urlToUpload), serv)
}

func uploadFileServerHandler(t *testing.T, fileBytesToCompare []byte, testDataStruct UploadFileResponse) http.Handler {

	testDataJson, _ := json.Marshal(testDataStruct)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, fileHeader, _ := r.FormFile("file")
		assert.Equal(t, "file.pdf", fileHeader.Filename)

		file, _ := fileHeader.Open()
		fileBody, _ := io.ReadAll(file)
		assert.Equal(t, fileBytesToCompare, fileBody)

		_, err := w.Write(testDataJson)
		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		require.NoError(t, err)
	})

	// handlerHelper(t, testData)
}

func Test_Service_UploadFile(t *testing.T) {
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

	s := New(&defaultConfig)
	file, err := os.Open("./testsData/file.pdf")

	require.NoError(t, err)
	fileBytes, err := io.ReadAll(file)
	require.NoError(t, err)

	fileId := "4991e6d7-5217-46ae-af3d-c9174adae924"

	mockUploadServer := getTestServer(uploadFileServerHandler(t, fileBytes, testData))

	got, err := s.UploadFile(mockUploadServer.URL, fileBytes, file.Name())

	require.NoError(t, err)
	assert.Equal(t, "https://gofile.io/d/Z19n9a", got.DownloadPage)
	assert.Equal(t, fileId, got.Token)
}

func deleteFileHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		reqBody := new(DeleteFileRequest)
		err = json.Unmarshal(body, &reqBody)
		require.NoError(t, err)
		assert.Equal(t, testDeleteFileData.ContentsId, reqBody.ContentsId)
		assert.Equal(t, testDeleteFileData.Token, reqBody.Token)
	})
	// handlerHelper(t, testData)
}

func Test_DeleteFile(t *testing.T) {
	s := New(&defaultConfig)

	mockServ := getTestServer(deleteFileHandler(t))

	urlDeleteFile = mockServ.URL
	err := s.DeleteFile(testDeleteFileData.ContentsId)
	require.NoError(t, err)
}
