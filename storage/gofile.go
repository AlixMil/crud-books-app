package storage

import (
	"bytes"
	"crud-books/config"
	"crud-books/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"reflect"
)

const (
	tokenName     = "token"
	contentIdName = "contentsId"
	fileName      = "file"

	сontentTypeHeaderName = "Content-Type"
	contentTypeJSON       = "application/json"
)

var (
	urlGetServer          = "https://api.gofile.io/getServer"
	urlDeleteFile         = "https://api.gofile.io/deleteContent"
	urlUploadFileTemplate = "https://%s.gofile.io/uploadFile"
)

type DataFromUploadFileResponse struct {
	DownloadPage string `json:"downloadPage"`
	Code         string `json:"code"`
	ParentFolder string `json:"parentFolder"`
	FileID       string `json:"fileId"`
	FileName     string `json:"fileName"`
	Md5          string `json:"md5"`
}

type UploadFileResponse struct {
	Status string                     `json:"status"`
	Data   DataFromUploadFileResponse `json:"data"`
}

type DataFromServerToUploadResponse struct {
	Server string `json:"server"`
}

type ServerToUploadResponse struct {
	Status string                         `json:"status"`
	Data   DataFromServerToUploadResponse `json:"data"`
}

type DeleteFileRequest struct {
	ContentsId string `json:"contentsId"`
	Token      string `json:"token"`
}

type Storage struct {
	apiKey string
	client *http.Client
}

type QueryParams struct {
	QueryParamName string
	QueryParamVal  string
}

func DoRequest(method, path string, queryParams []QueryParams, body *bytes.Buffer) (*http.Request, error) {
	if body == nil {
		body = nil
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return &http.Request{}, err
	}

	for _, v := range queryParams {
		q := req.URL.Query()
		q.Add(v.QueryParamName, v.QueryParamVal)
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (s Storage) GetServerToUpload() (string, error) {
	req, err := DoRequest(
		http.MethodGet,
		urlGetServer,
		[]QueryParams{},
		&bytes.Buffer{},
	)
	if err != nil {
		return "", fmt.Errorf("error in request of getServerToUpload : %w", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error when client send request in getServerToUpload: %w", err)
	}

	jsonRes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading body of req in getServerToUpload: %w", err)
	}

	jBody := ServerToUploadResponse{}
	err = json.Unmarshal(jsonRes, &jBody)
	if err != nil {
		return "", fmt.Errorf("error in unmarshalling of req body in getServerToUpload: %w", err)
	}

	serverAddress := fmt.Sprintf(urlUploadFileTemplate, jBody.Data.Server)

	return serverAddress, nil
}

func createPdfFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", "application/pdf")

	return w.CreatePart(h)
}

func getBodyWriter(fileName string, file []byte, apiKey string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileWriter, err := createPdfFormFile(writer, fileName)
	if err != nil {
		return body, "", err
	}

	_, err = fileWriter.Write(file)
	if err != nil {
		return body, "", err
	}

	err = writer.WriteField(tokenName, apiKey)
	if err != nil {
		return body, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func (s Storage) UploadFile(servForUpload string, file []byte, fileName string) (*models.FileData, error) {
	body, contentType, err := getBodyWriter(filepath.Base(fileName), file, s.apiKey)
	if err != nil {
		return nil, fmt.Errorf("getbodywrite in uploadfile throw error: %w", err)
	}

	req, err := DoRequest(
		http.MethodPost,
		servForUpload,
		[]QueryParams{},
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("error in req of UploadFile: %w", err)
	}

	req.Header.Add("Content-Type", contentType)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("uploadFile func sending request failed: %w", err)
	}

	jsonResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("upload file response : %s", jsonResponse)

	var uploadResp UploadFileResponse
	err = json.Unmarshal(jsonResponse, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal of uploadfile: %w", err)
	}

	if reflect.ValueOf(uploadResp.Data).IsZero() {
		return nil, fmt.Errorf("upload file error, service storage unexpected response")
	}

	return &models.FileData{
		Token:        uploadResp.Data.FileID,
		DownloadPage: uploadResp.Data.DownloadPage,
	}, nil
}

func (s Storage) DeleteFile(fileToken string) error {
	j := DeleteFileRequest{
		ContentsId: fileToken,
		Token:      s.apiKey,
	}
	jsonBody, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("error in marshal of delete file: %w", err)
	}

	req, err := DoRequest(
		http.MethodDelete,
		urlDeleteFile,
		[]QueryParams{},
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error in request of delete file: %w", err)
	}

	req.Header.Add(сontentTypeHeaderName, contentTypeJSON)

	_, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when sending response of delete file: %w", err)
	}

	return nil
}

func New(cfg *config.Config) *Storage {
	s := Storage{
		apiKey: cfg.GoFileServiceApiKey,
		client: http.DefaultClient,
	}
	return &s
}
