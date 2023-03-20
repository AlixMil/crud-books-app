package models

type User struct {
	Id       string   `bson:"_id,omitempty"`
	Email    string   `bson:"email"`
	Password string   `bson:"password"`
	Books    []string `bson:"books"`
}

type Book struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	FileToken   string `bson:"fileToken"`
	Owner       string `bson:"owner"`
}
