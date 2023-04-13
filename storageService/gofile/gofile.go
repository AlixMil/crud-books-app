package gofile

import (
	"bytes"
	"crud-books/models"
	storageService_helpers "crud-books/storageService"
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
	urlGetServer    = "https://api.gofile.io/getServer"
	urlDeleteFile   = "https://api.gofile.io/deleteContent"
	urlUploadServer = ""
)

type Storage struct {
	StorageService string
	apiKey         string
	client         *http.Client
	folderId       string
}

type queryParam = storageService_helpers.QueryParams

func (s Storage) getServerToUpload() (string, error) {
	req, err := storageService_helpers.DoRequest(
		http.MethodGet,
		urlGetServer,
		[]queryParam{},
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

	jBody := models.UploadServerSummary{}
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

func (s Storage) UploadFile(file []byte, isTest bool) (*models.UploadFileReturn, error) {
	serverToUpload, err := s.getServerToUpload()
	if err != nil {
		return nil, fmt.Errorf("error in serverToUpload getting of UploadFile: %w", err)
	}

	body, contentType, err := getBodyWriter(file, s.apiKey, s.folderId)
	if err != nil {
		return nil, fmt.Errorf("getbodywrite in uploadfile throw error: %w", err)
	}

	if !isTest {
		urlUploadServer = fmt.Sprintf("https://%v.gofile.io/UploadFile", serverToUpload)
	}

	req, err := storageService_helpers.DoRequest(
		http.MethodPost,
		urlUploadServer,
		[]queryParam{},
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

	var uploadResp models.UploadFileResponse
	err = json.Unmarshal(jsonResponse, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal of uploadfile: %w", err)
	}

	return &models.UploadFileReturn{DownloadPage: uploadResp.Data.DownloadPage, FileToken: uploadResp.Data.FileID}, nil
}

func (s Storage) DeleteFile(fileToken string) error {
	type jsonScheme struct {
		ContentsId string `json:"contentsId"`
		Token      string `json:"token"`
	}
	j := jsonScheme{
		ContentsId: fileToken,
		Token:      s.apiKey,
	}
	jsonBody, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("error in marshal of delete file: %w", err)
	}
	req, err := storageService_helpers.DoRequest(
		http.MethodDelete,
		urlDeleteFile,
		[]queryParam{},
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("error in request of delete file: %w", err)
	}
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	_, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when sendind response of delete file: %w", err)
	}
	return nil
}

func New(apiKey, folderId string) *Storage {
	s := Storage{
		StorageService: "gofile",
		apiKey:         apiKey,
		client:         http.DefaultClient,
		folderId:       folderId,
	}
	return &s
}
