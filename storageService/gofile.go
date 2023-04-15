package storageService

import (
	"bytes"
	"crud-books/config"
	"crud-books/models"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	paramTokenName     = "token"
	paramContentIdName = "contentsId"
	paramFileName      = "file"
	paramFolderIdName  = "folderId"
)

var (
	urlGetServer  = "https://api.gofile.io/getServer"
	urlDeleteFile = "https://api.gofile.io/deleteContent"
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

type DeleteFileScheme struct {
	ContentsId string `json:"contentsId"`
	Token      string `json:"token"`
}

type Storage struct {
	StorageService string
	apiKey         string
	client         *http.Client
	folderId       string
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

	return jBody.Data.Server, nil
}

func getBodyWriter(file []byte, apiKey, folderId string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	fileWriter, err := writer.CreateFormFile(paramFileName, "file")
	if err != nil {
		return body, "", err
	}

	_, err = fileWriter.Write(file)
	if err != nil {
		return body, "", err
	}

	err = writer.WriteField(paramTokenName, apiKey)
	if err != nil {
		return body, "", err
	}

	err = writer.WriteField(paramFolderIdName, folderId)
	if err != nil {
		return body, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func (s Storage) UploadFile(servForUpload string, file []byte) (*models.FileData, error) {
	body, contentType, err := getBodyWriter(file, s.apiKey, s.folderId)
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
		return nil, fmt.Errorf("error when uploadfile func sending request: %w", err)
	}

	jsonResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var uploadResp UploadFileResponse
	err = json.Unmarshal(jsonResponse, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal of uploadfile: %w", err)
	}

	return &models.FileData{
		Token:        uploadResp.Data.FileID,
		DownloadPage: uploadResp.Data.DownloadPage,
	}, nil
}

func (s Storage) DeleteFile(fileToken string) error {
	j := DeleteFileScheme{
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

	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	_, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when sending response of delete file: %w", err)
	}

	return nil
}

func New(cfg config.Config) *Storage {
	s := Storage{
		StorageService: "gofile",
		apiKey:         cfg.GoFileServiceApiKey,
		client:         http.DefaultClient,
		folderId:       cfg.GoFileFolderToken,
	}
	return &s
}