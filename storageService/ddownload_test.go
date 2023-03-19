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

func Test_DoRequestWQuery(t *testing.T) {
	const apiKeyVal = "12312"
	mockServer := givenTestServer(parsedRequestHandler(apiKeyVal))
	urlForTest := mockServer.URL

	req, err := doRequest(http.MethodGet, urlForTest, []queryParams{{queryParamName: apiKeyParamName, queryParamVal: apiKeyVal}}, &bytes.Buffer{})
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

type handlerOption func(w http.ResponseWriter, r *http.Request) error

func checkMultipartBody(t *testing.T, expectFileData []byte, sessId string) handlerOption {
	return func(w http.ResponseWriter, r *http.Request) error {
		// 10 mb
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
		var (
			actualSessId string
			actualUType  string
		)

		if len(r.MultipartForm.Value["sess_id"]) != 0 {
			actualSessId = r.MultipartForm.Value["sess_id"][0]
		} else {
			actualSessId = ""
		}

		if len(r.MultipartForm.Value["utype"]) != 0 {
			actualUType = r.MultipartForm.Value["utype"][0]
		} else {
			actualUType = ""
		}

		require.Equal(t, sessId, actualSessId)
		require.Equal(t, "prem", actualUType)

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
	fileName := "./testsData/getServerToUploadSuccessResponse.json"

	file, err := os.OpenFile(fileName, os.O_RDWR, 0777)
	require.NoError(t, err)
	defer file.Close()
	j, err := io.ReadAll(file)
	require.NoError(t, err)

	var v getServerUploadResponse
	err = json.Unmarshal(j, &v)
	require.NoError(t, err)
	v.SessID = "1cvxk3uo91qtafcx8k7vptntn65nek1z2xybg4sicmk7jj3413"

	r, err := json.Marshal(v)
	require.NoError(t, err)
	err = os.WriteFile(fileName, r, 0777)
	require.NoError(t, err)
	return handlerHelper(t, "./testsData/uploadFileServerResponse.json", checkMultipartBody(t, body, v.SessID))
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

func checkGetFileInfoRequest(t *testing.T, expectQueryData []queryParams) handlerOption {
	return func(w http.ResponseWriter, r *http.Request) error {
		for _, q := range expectQueryData {
			require.Equal(t, q.queryParamVal, r.URL.Query().Get(q.queryParamName))
		}
		return nil
	}
}

func succesGetFileInfoHandler(t testing.T, fileCode, apiKey string) http.Handler {
	return handlerHelper(&t, "./testsData/getFileInfoResponse.json", checkGetFileInfoRequest(&t, []queryParams{
		{queryParamName: fileCodeParamName, queryParamVal: fileCode},
		{queryParamName: apiKeyParamName, queryParamVal: apiKey},
	}))
}

func TestService_GetFileInfo_Success(t *testing.T) {
	s := New("asjdkasjdk")
	fileCode := "1ahye98t2y6r"
	mockGetFileInfoServer := givenTestServer(succesGetFileInfoHandler(*t, fileCode, s.apiKey))

	urlGetFileInfo = mockGetFileInfoServer.URL

	got, err := s.getFileInfo(fileCode)
	require.NoError(t, err)

	assert.Equal(t, "1ahye98t2y6r", got.Result[0].Filecode)
}
