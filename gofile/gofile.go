package gofile

import "crud-books/storage"

type goFile struct {
	apiKey string
}

func New(apiKey string) goFile {
	service := goFile{apiKey: apiKey}
	return service
}

type SaveFileResponse struct {
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

func (s *SaveFileResponse) Token() string {
	return s.Data.FileID
}

func (g goFile) SaveFile(file []byte, title string) (storage.TokenGetter, error) {
	// do request
	return &SaveFileResponse{}, nil
}

func (g goFile) getFile()
