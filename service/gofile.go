package gofile

import (
	"bytes"
	"crud-books/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	getServerAddress = "https://api.gofile.io/getServer"
	mainApiAddress   = "https://api.gofile.io"
)

type goFile struct {
	apiKey string
}

type UploadFileResponse struct {
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

type UloadFileRequest struct {
	file  []byte
	token string
}

func (g goFile) UploadFile(file []byte) (*UploadFileResponse, error) {
	// receive most speed server
	serverResponse, err := g.GetBestServer()
	if err != nil {
		return &UploadFileResponse{}, err
	}

	// do request

	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	part, _ := writer.CreateFormFile("file", "filename")
	part.Write(file)

	res, err := http.Post(helpers.UploadAddress(serverResponse.Data.Server), "multipart/form-data", buffer)
	if err != nil {
		return &UploadFileResponse{}, err
	}

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &UploadFileResponse{}, err
	}
	fmt.Println(response)
	return &UploadFileResponse{}, nil
}

func (u UploadFileResponse) UploadFileSummary() map[string]string {
	return map[string]string{
		"DownloadPage": u.Data.DownloadPage,
		"FileId":       u.Data.FileID,
		"FileName":     u.Data.FileName,
		"Status":       u.Status,
	}
}

func (g goFile) GetBestServer() (*GetServerResponse, error) {
	res, err := http.Get("https://api.gofile.io/getServer")
	if err != nil {
		return &GetServerResponse{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &GetServerResponse{}, err
	}

	responseBody := GetServerResponse{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return &GetServerResponse{}, err
	}

	fmt.Printf("Response body from GetBest: %v", responseBody)
	// fetch GET https://api.gofile.io/getServer
	return &responseBody, nil
}

type GetServerResponse struct {
	Status string `json:"status"`
	Data   struct {
		Server string `json:"server"`
	} `json:"data"`
}

func (s GetServerResponse) GetServer() string {
	return s.Data.Server
}

func (g goFile) DeleteFile(contentId string) (*DeleteFileResponse, error) {
	// do request
	return &DeleteFileResponse{}, nil
}

func (g goFile) GetFile(contentId string) (*GetFileResponse, error) {
	// do request
	return &GetFileResponse{}, nil
}

type GetFileResponse struct {
	Status string `json:"status"`
	Data   struct {
		IsOwner            bool     `json:"isOwner"`
		ID                 string   `json:"id"`
		Type               string   `json:"type"`
		Name               string   `json:"name"`
		ParentFolder       string   `json:"parentFolder"`
		Code               string   `json:"code"`
		CreateTime         int      `json:"createTime"`
		Public             bool     `json:"public"`
		Childs             []string `json:"childs"`
		TotalDownloadCount int      `json:"totalDownloadCount"`
		TotalSize          int      `json:"totalSize"`
		Contents           struct {
			ItemId struct {
				ID            string `json:"id"`
				Type          string `json:"type"`
				Name          string `json:"name"`
				ParentFolder  string `json:"parentFolder"`
				CreateTime    int    `json:"createTime"`
				Size          int    `json:"size"`
				DownloadCount int    `json:"downloadCount"`
				Md5           string `json:"md5"`
				Mimetype      string `json:"mimetype"`
				ServerChoosen string `json:"serverChoosen"`
				DirectLink    string `json:"directLink"`
				Link          string `json:"link"`
			} `json:""`
		} `json:"contents"`
	} `json:"data"`
}

type DeleteFileResponse struct {
	Status string `json: status`
	Data   string `json: data`
}

func (d DeleteFileResponse) DeleteStatus() bool {
	if d.Status == "ok" {
		return true
	}
	return false
}

func New(apiKey string) goFile {
	service := goFile{apiKey: apiKey}
	return service
}
