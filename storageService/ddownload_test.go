package storageService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestService_GetServerToUpload(t *testing.T) {
// 	// setup envirenment
// 	apiKeyService := "SDjkasjdkh$U%(U$%hkdjaks"
// 	service := New(apiKeyService)
// 	// execute
// 	getServerResponse, err := service.GetServerToUpload()
// 	if err != nil {
// 		t.Fail()
// 	}
// 	// assert
// 	_ = getServerResponse
// }

func TestService_GetServerToUpload_NotFound(t *testing.T) {
	// setup environment
	apiKey := "some api key of the ddownload service"
	s := New(apiKey)
	mockServer := givenTestServer(alwaysNotFoundHandler())

	urlGetServerToUpload = mockServer.URL
	// execution
	_, err := s.GetServerToUpload()
	// assert
	assert.EqualError(t, err, "status code 404")
}

func givenTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func alwaysNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

func parsedRequestHandler(apiKeyVal string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get(apiKeyParamName) != apiKeyVal {
			log.Fatalf("Apikey doesn't exist in query of request")
		}
	})
}

func TestService_UploadFile_GetServerToUpload_RequestError(t *testing.T) {
	s := New("jfaksjdkjas")
	mockServer := givenTestServer(alwaysNotFoundHandler())
	urlGetServerToUpload = mockServer.URL

	_, err := s.UploadFile(nil)
	assert.EqualError(t, err, "get server error: status code 404")
}

// test for func doReqeustWQuery
func Test_DoRequestWQuery(t *testing.T) {
	const apiKeyVal = "12312"
	mockServer := givenTestServer(parsedRequestHandler(apiKeyVal))
	urlForTest := mockServer.URL

	req, err := doRequestWQuery(http.MethodGet, urlForTest, apiKeyParamName, apiKeyVal)
	if err != nil {
		t.Errorf("request forming was error: %s", err.Error())
	}
	c := http.Client{}
	_, err = c.Do(req)
	if err != nil {
		log.Fatalf("err borned when request was sent: %v", err.Error())
	}
}

func TestService_UploadFile_GetServerToUploadRequestError(t *testing.T) {
	s := New("jfaksjdkjas")
	mockServer := givenTestServer(alwaysNotFoundHandler())
	urlGetServerToUpload = mockServer.URL

	_, err := s.UploadFile([]byte{1, 2, 3, 4})

	assert.Error(t, err, "get server error: status code 404")
}

// func TestService_UploadFile_GetServerToUpload(t *testing.T) {
// 	s := New("jfaksjdkjas")
// 	mockServer := givenTestServer(alwaysNotFoundHandler())
// 	urlGetServerToUpload = mockServer.URL

// 	_, err := s.UploadFile([]byte{1, 2, 3, 4})

// 	assert.Error(t, err, "get server error: status code 404")
// }

// {
//     "msg": "OK",
//     "server_time": "2017-08-11 04:29:54",
//     "status": 200,
//     "sess_id": "1cvxk3uo91qtafcx8k7vptntn65nek1z2xybg4sicmk7jjbl5n",
//     "result": "https://wwwNNN.ucdn.to/cgi-bin/upload.cgi"
// }

type handlerOption func(w http.ResponseWriter, r *http.Request) error

func checkMultipartBody(t *testing.T, expectFileData []byte) handlerOption {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseMultipartForm(0)

		buff := bytes.NewBuffer([]byte{})
		for _, fileHeaders := range r.MultipartForm.File {
			for _, fh := range fileHeaders {
				f, err := fh.Open()
				if err != nil {
					return fmt.Errorf("open file `%s`: %w", fh.Filename, err)
				}
				_, err = buff.ReadFrom(f)
				if err != nil {
					return fmt.Errorf("read from file `%s`: %w", fh.Filename, err)
				}

			}
		}
		require.Equal(t, string(expectFileData), buff.String())
		return nil
	}
}

func handlerHelper(t *testing.T, path string, options ...handlerOption) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, f := range options {
			err := f(w, r)
			assert.NoError(t, err)
		}
		file, err := os.OpenFile(path, os.O_RDONLY, 0777)
		assert.NoError(t, err)
		defer file.Close()
		responseFromFile, err := io.ReadAll(file)
		assert.NoError(t, err)
		w.Write(responseFromFile)
	})
}

func getServerToUploadHandler(t *testing.T, url string) http.Handler {
	fileName := "./testsData/getServerToUploadSuccessResponse.json"

	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	require.NoError(t, err)
	defer file.Close()
	j, err := io.ReadAll(file)
	require.NoError(t, err)

	var v getServerUploadResponse
	err = json.Unmarshal(j, &v)
	require.NoError(t, err)
	v.Result = url

	r, err := json.Marshal(v)
	require.NoError(t, err)
	err = os.WriteFile(fileName, r, 0777)
	require.NoError(t, err)

	return handlerHelper(t, fileName)
}

func uploadFileServerHandler(t *testing.T, body []byte) http.Handler {
	return handlerHelper(t, "./testsData/uploadFileServerResponse.json", checkMultipartBody(t, body))
}

func TestService_UploadFile_Success(t *testing.T) {
	s := New("asjdkasjdk")
	uploadBody := []byte("Hello, World!")
	mockUploadServer := givenTestServer(uploadFileServerHandler(t, uploadBody))
	mockGetServer := givenTestServer(getServerToUploadHandler(t, mockUploadServer.URL))

	urlGetServerToUpload = mockGetServer.URL

	got, err := s.UploadFile(uploadBody)
	require.NoError(t, err)
	assert.Equal(t, "yzanp0ps7sgl", got)
}
