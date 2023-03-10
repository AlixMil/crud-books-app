package storageservice

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type Service struct {
	ServiceName      string
	apiKey           string
	pathsForRequests map[string]string
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

func (s Service) getServerToUpload() (UploadServerSummary, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, s.pathsForRequests["GetServerToUpload"], nil)
	if err != nil {
		return UploadServerSummary{}, err
	}
	q := req.URL.Query()
	q.Add("key", s.apiKey)
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return UploadServerSummary{}, err
	}
	defer res.Body.Close()
	jsonResponse, err := io.ReadAll(req.Body)
	if err != nil {
		return UploadServerSummary{}, err
	}

	var serverUpload getServerUploadResponse
	err = json.Unmarshal(jsonResponse, &serverUpload)
	if err != nil {
		return UploadServerSummary{}, err
	}

	return UploadServerSummary{
		SessId: serverUpload.SessID,
		Server: serverUpload.Result,
	}, nil
}

type UploadFileResponse []struct {
	FileCode   string `json:"file_code"`
	FileStatus string `json:"file_status"`
}

func (s Service) UploadFile(file []byte) (string, error) {
	server, err := s.getServerToUpload()
	if err != nil {
		return "", err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormField("file")
	if err != nil {
		return "", err
	}
	fileWriter.Write(file)
	writer.WriteField("sess_id", server.SessId)
	writer.WriteField("utype", "prem")

	// data := url.Values{}
	// data.Set("sess_id", server.SessId)
	// data.Set("utype", "prem")

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, server.Server, body)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
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
		ServiceName: "ddownload",
		apiKey:      apikey,
		pathsForRequests: map[string]string{
			"GetServerToUpload": "https://api-v2.ddownload.com/api/upload",
			"accountInfo":       "https://api-v2.ddownload.com/api/account/info",
			"getFileInfo":       "https://api-v2.ddownload.com/api/file/info",
			"getFilesList":      "https://api-v2.ddownload.com/api/file/list",
			"renameFile":        "https://api-v2.ddownload.com/api/file/rename",
		},
	}

	return s
}
