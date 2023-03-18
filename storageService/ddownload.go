package storageService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	apiKeyParamName = "key"
)

// URLs for requests
var (
	urlGetServerToUpload = "https://api-v2.ddownload.com/api/upload"
	urlAccountInfo       = "https://api-v2.ddownload.com/api/account/info"
	urlGetFileInfo       = "https://api-v2.ddownload.com/api/file/info"
	urlGetFilesList      = "https://api-v2.ddownload.com/api/file/list"
	urlRenameFile        = "https://api-v2.ddownload.com/api/file/rename"
)

type Service struct {
	ServiceName string
	apiKey      string
	client      *http.Client
}

type getServerUploadResponse struct {
	Msg        string `json:"msg"`
	ServerTime string `json:"server_time"`
	Status     int    `json:"status"`
	SessID     string `json:"sess_id"`
	Result     string `json:"result"`
}

type UploadServerSummary struct {
	SessId string
	Server string
}

func doRequestWQuery(method, path, queryParamName, queryParamVal string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return &http.Request{}, err
	}
	q := req.URL.Query()
	q.Add(queryParamName, queryParamVal)
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func doRequestWBody(method, path string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func (s Service) GetServerToUpload() (*UploadServerSummary, error) {
	// Do request to server for getting actual server and session. Needed for upload file to correct server with session
	req, err := doRequestWQuery(http.MethodGet, urlGetServerToUpload, apiKeyParamName, s.apiKey)
	if err != nil {
		return &UploadServerSummary{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("status code %d", http.StatusNotFound)
	}
	defer res.Body.Close()

	jsonResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var serverUpload getServerUploadResponse
	err = json.Unmarshal(jsonResponse, &serverUpload)
	if err != nil {
		return nil, err
	}

	return &UploadServerSummary{
		SessId: serverUpload.SessID,
		Server: serverUpload.Result,
	}, nil
}

type UploadFileResponse []struct {
	FileCode   string `json:"file_code"`
	FileStatus string `json:"file_status"`
}

func getBodyWriter(file []byte, sessId, utype string) (*bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormField("file")
	if err != nil {
		return body, err
	}
	_, err = fileWriter.Write(file)
	if err != nil {
		return body, err
	}

	err = writer.WriteField("sess_id", sessId)
	if err != nil {
		return body, err
	}
	err = writer.WriteField("utype", utype)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (s Service) UploadFile(file []byte) (string, error) {
	uploadServerSummary, err := s.GetServerToUpload()
	if err != nil {
		return "", fmt.Errorf("get server error: %v", err)
	}

	body, err := getBodyWriter(file, uploadServerSummary.SessId, "prem")
	if err != nil {
		return "", err
	}

	req, err := doRequestWBody(http.MethodPost, uploadServerSummary.Server, body)
	if err != nil {
		return "", err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	jsonResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var uploadResp UploadFileResponse

	err = json.Unmarshal(jsonResponse, &uploadResp)
	if err != nil {
		return "", err
	}

	return uploadResp[0].FileCode, nil

}

func New(apikey string) Service {
	s := Service{
		apiKey: apikey,
		client: http.DefaultClient,
	}

	return s
}
