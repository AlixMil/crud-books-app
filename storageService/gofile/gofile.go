package gofile

import (
	"bytes"
	storageService_helpers "crud-books/storageService"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

type Service struct {
	ServiceName string
	apiKey      string
	client      *http.Client
	folderId    string
}

type UploadServerSummary struct {
	Status string `json:"status"`
	Data   struct {
		Server string `json:"server"`
	} `json:"data"`
}

type queryParam = storageService_helpers.QueryParams

type uploadFileResponse struct {
	Status string `json:"status"`
	Data   struct {
		DownloadPage string `json:"downloadPage"`
		Code         string `json:"code"`
		ParentFolder string `json:"parentFolder"`
		FileID       string `json:"fileId"`
		FileName     string `json:"fileName"`
		Md5          string `json:"md5"`
	} `json:"data"`
}

func (s Service) getServerToUpload() (string, error) {
	req, err := storageService_helpers.DoRequest(
		http.MethodGet,
		urlGetServer,
		[]queryParam{},
		&bytes.Buffer{},
	)
	if err != nil {
		return "", fmt.Errorf("error in request of getServerToUpload : %w", err)
	}
	jsonRes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading body of req in getServerToUpload: %w", err)
	}

	jBody := UploadServerSummary{}
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

func (s Service) UploadFile(file []byte, isTest bool) (string, error) {
	serverToUpload, err := s.getServerToUpload()
	if err != nil {
		return "", fmt.Errorf("error in serverToUpload getting of UploadFile: %w", err)
	}

	body, contentType, err := getBodyWriter(file, s.apiKey, s.folderId)
	if err != nil {
		return "", fmt.Errorf("getbodywrite in uploadfile throw error: %w", err)
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
		return "", fmt.Errorf("error in req of UploadFile: %w", err)
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

	var uploadResp uploadFileResponse
	err = json.Unmarshal(jsonResponse, &uploadResp)
	if err != nil {
		return "", fmt.Errorf("error in unmarshal of uploadfile: %w", err)
	}

	return uploadResp.Data.DownloadPage, nil
}

func (s Service) DeleteFile(fileToken string) error {
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
	req.Header.Add("Content-Type", "application/json")
	_, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when sendind response of delete file: %w", err)
	}
	return nil
}

func New(apiKey, folderId string) *Service {
	s := Service{
		ServiceName: "gofile",
		apiKey:      apiKey,
		client:      http.DefaultClient,
		folderId:    folderId,
	}
	return &s
}
