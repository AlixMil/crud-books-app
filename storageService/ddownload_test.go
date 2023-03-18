package storageService

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestService_UploadFile_GetServerToUpload(t *testing.T) {
	s := New("jfaksjdkjas")
	mockServer := givenTestServer(alwaysNotFoundHandler())
	urlGetServerToUpload = mockServer.URL

	_, err := s.UploadFile([]byte{1, 2, 3, 4})

	assert.Error(t, err, "get server error: status code 404")
}
