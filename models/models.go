package models

type UserData struct {
	Id           string   `bson:"_id,omitempty"`
	Email        string   `bson:"email"`
	PasswordHash string   `bson:"passwordHash"`
	BooksIds     []string `bson:"booksIds"`
}

type BookData struct {
	Id          string `bson:"_id,omitempty"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	FileToken   string `bson:"fileToken"`
	Url         string `bson:"url"`
	OwnerEmail  string `bson:"owner"`
}

type BookDataUpdater struct {
	FileToken   string `json:"fileToken"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UserDataInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FileData struct {
	Id           string `bson:"_id,omitempty"`
	Token        string `bson:"token"`
	DownloadPage string `bson:"downloadPage"`
}

type Filter struct {
	Email  string
	Search string
}

type Sort struct {
	SortField string
	Limit     int
	Direction string
	Offset    int
}

type ValidateDataInGetLists struct {
	Email  string
	Search string

	SortField string
	Limit     int
	Direction int

	Offset int
}

// STORAGE RESPONSES

type UploadFileReturn struct {
	DownloadPage string
	FileToken    string
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

type UploadServerSummary struct {
	Status string `json:"status"`
	Data   struct {
		Server string `json:"server"`
	} `json:"data"`
}

type CreateBookRequest struct {
	FileToken   string `json:"fileToken"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetBookResponse struct {
	FileURL     string `json:"fileURL"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
