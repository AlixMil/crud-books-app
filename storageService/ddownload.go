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
	urlService           = "https://ddownload.com/"
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

type queryParams struct {
	queryParamName string
	queryParamVal  string
}

type getFileInfoResponse struct {
	Msg        string `json:"msg"`
	ServerTime string `json:"server_time"`
	Status     int    `json:"status"`
	Result     []struct {
		Uploaded string `json:"uploaded"`
		Status   int    `json:"status"`
		Filecode string `json:"filecode"`
		Name     string `json:"name"`
		Size     string `json:"size"`
	} `json:"result"`
}

func doRequest(method, path string, queryParams []queryParams, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return &http.Request{}, err
	}

	for _, v := range queryParams {
		q := req.URL.Query()
		q.Add(v.queryParamName, v.queryParamVal)
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (s Service) GetServerToUpload() (*UploadServerSummary, error) {
	queryParams := []queryParams{{queryParamName: apiKeyParamName, queryParamVal: s.apiKey}}
	req, err := doRequest(http.MethodGet, urlGetServerToUpload, queryParams, nil)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("json unmarshal error: %w", err)
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

func getBodyWriter(file []byte, sessId, utype string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileWriter, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return body, "", err
	}

	_, err = fileWriter.Write(file)
	if err != nil {
		return body, "", err
	}

	err = writer.WriteField("sess_id", sessId)
	if err != nil {
		return body, "", err
	}
	err = writer.WriteField("utype", utype)
	if err != nil {
		return body, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func (s Service) UploadFile(file []byte) (string, error) {
	uploadServerSummary, err := s.GetServerToUpload()
	if err != nil {
		return "", fmt.Errorf("get server error: %v", err)
	}

	body, contentType, err := getBodyWriter(file, uploadServerSummary.SessId, "prem")
	if err != nil {
		return "", err
	}

	req, err := doRequest(http.MethodPost, uploadServerSummary.Server, []queryParams{}, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", contentType)

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

	if len(uploadResp) == 0 {
		return "", fmt.Errorf("empty upload response, error: %w", err)
	}

	return uploadResp[0].FileCode, nil

}

func (s Service) getFileInfo() (*getFileInfoResponse, error) {
	req, err := doRequest(http.MethodGet, urlGetFileInfo, []queryParams{{queryParamName: apiKeyParamName, queryParamVal: s.apiKey}}, nil)
	if err != nil {
		return nil, fmt.Errorf("error in request of getFileLink func: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in response of getFileLink func: %w", err)
	}

	defer res.Body.Close()
	jsonRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error in jsonresponse of getFileLink func: %w", err)
	}

	var getFileRes getFileInfoResponse
	json.Unmarshal(jsonRes, &getFileRes)

	return &getFileRes, nil
}

func New(apikey string) Service {
	s := Service{
		apiKey: apikey,
		client: http.DefaultClient,
	}

	return s
}
