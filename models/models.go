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
	Owner       string `bson:"owner"`
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
}

type ValidateDataInGetLists struct {
	Email  string
	Search string

	SortField string
	Limit     int
	Direction int
}
