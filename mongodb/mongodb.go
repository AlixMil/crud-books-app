package mongodb

import (
	"context"
	"crud-books/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollectionName = "users"
	booksCollectionName = "books"
	filesCollectionName = "files"
)

type MongoDB struct {
	DB *mongo.Database
}

type InsertedID struct {
	InsertedID interface{}
}

type UserProperties struct {
	Id          string   `bson:"_id,omitempty"`
	Email       string   `bson:"email"`
	Password    string   `bson:"password"`
	linkToBooks []string `bson:"linkToBooks"`
}

type BooksProperties struct {
	Id          string `bson:"_id,omitempty"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	FileToken   string `bson:"fileToken"`
}

func (m MongoDB) GetBook(bookToken string) (*models.BookData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(booksCollectionName)

	filter := bson.M{"fileToken": bookToken}
	result := coll.FindOne(ctx, filter)
	var bookData models.BookData
	err := result.Decode(&bookData)
	if err != nil {
		return nil, fmt.Errorf("decode book data in get book func of db failed, error: %w", err)
	}
	return &bookData, nil
}

func (m *MongoDB) CreateUser(email, passwordHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection("users")

	filter := bson.M{"email": email}
	res := coll.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("user with email %v already created", email)
	}

	doc := models.UserData{Email: email, PasswordHash: passwordHash, BooksIds: []string{}}

	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed collection insert in add user func, error: %w", err)
	}
	return nil
}

func (m *MongoDB) CreateBook(title, description, fileToken, emailOwner string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	books := m.DB.Collection(booksCollectionName)
	bookData, err := m.GetFileData(fileToken)
	if err != nil {
		return "", fmt.Errorf("get file data failed, error: %w", err)
	}
	dBook := models.BookData{
		Title:       title,
		Description: description,
		FileToken:   fileToken,
		Owner:       emailOwner, // email of owner
		Url:         bookData.DownloadPage,
	}

	insertResult, err := books.InsertOne(ctx, dBook)
	if err != nil {
		return "", fmt.Errorf("failed collection insert in create book func, error: %w", err)
	}
	users := m.DB.Collection(usersCollectionName)

	filter := bson.M{"email": emailOwner}
	update := bson.M{"$set": bson.M{"$push": bson.M{
		"books": insertResult.InsertedID,
	}}}

	users.FindOneAndUpdate(ctx, filter, update)

	return dBook.FileToken, nil
}

type ValidateData struct {
	Email  string
	Search string

	SortField string
	Limit     int
	Direction int
}

func GetParamsWValidate(email, search, sortField, direction string, limit int) (*ValidateData, error) {
	var res ValidateData
	if email == "" {
		return nil, fmt.Errorf("email is not provided")
	} else {
		res.Email = email
	}
	if sortField == "title" {
		res.SortField = "title"
	} else if sortField == "date" {
		res.SortField = "date"
	} else {
		return nil, fmt.Errorf("provided sorting parameter incorrect")
	}
	if direction == "asc" {
		res.Direction = 1
	} else if direction == "desc" {
		res.Direction = -1
	} else {
		return nil, fmt.Errorf("provided direction parameter incorrect")
	}

	if limit > 100 || limit < 5 {
		return nil, fmt.Errorf("you should provide limit parameter in range between from 5 to 100")
	}
	res.Limit = limit
	res.Search = search
	return &res, nil
}

func (m *MongoDB) GetListOfBooks(filter models.Filter, sorting models.Sort) (*[]models.BookData, error) {

	validParams, err := GetParamsWValidate(filter.Email, filter.Search, sorting.SortField, sorting.Direction, sorting.Limit)
	if err != nil {
		return nil, fmt.Errorf("get params with validate failed, error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userCollection := m.DB.Collection(usersCollectionName)
	emailFilter := bson.M{"email": validParams.Email}
	userCur, err := userCollection.Find(ctx, emailFilter)
	if err != nil {
		return nil, fmt.Errorf("user search failed, error: %w", err)
	}
	var userData models.UserData
	err = userCur.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("decoding user data failed, error: %w", err)
	}

	bookCollection := m.DB.Collection(booksCollectionName)
	bookReqFilter := bson.M{"owner": userData.Id}

	sortOpt := options.Find().SetSort(bson.M{validParams.SortField: validParams.Direction})
	searchOpt := options.Find().SetLimit(int64(validParams.Limit))
	booksCur, err := bookCollection.Find(ctx, bookReqFilter, sortOpt, searchOpt)
	if err != nil {
		return nil, fmt.Errorf("book collection search failed, error: %w", err)
	}
	var books *[]models.BookData
	booksCur.All(ctx, &books)
	return books, nil
}

func (m *MongoDB) ChangeFieldOfBook(collectionName, id, fieldName, fieldValue string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(collectionName)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{fieldName: fieldValue}}

	coll.FindOneAndUpdate(ctx, filter, update)

	return nil
}

func (m *MongoDB) UploadFileData(fileToken, downloadPage string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(filesCollectionName)

	doc := bson.M{"token": fileToken, "downloadPage": downloadPage}

	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insertone in uploadfiledata failed, error: %w", err)
	}
	return nil
}

func (m MongoDB) GetFileData(fileToken string) (*models.FileData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(filesCollectionName)

	filter := bson.M{"token": fileToken}

	res := coll.FindOne(ctx, filter)
	var f models.FileData
	err := res.Decode(&f)
	if err != nil {
		return nil, fmt.Errorf("decoding result of getfiledata failed, error: %w", err)
	}
	return &models.FileData{Id: f.Id, Token: f.Token, DownloadPage: f.DownloadPage}, nil
}

func (m *MongoDB) GetUserData(email string) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(usersCollectionName)

	filter := bson.M{"email": email}
	res := coll.FindOne(ctx, filter)
	var userData models.UserData
	err := res.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding finded user data in getUserData, error: %w", err)
	}

	return &userData, nil
}

func (m *MongoDB) GetUserDataByInsertedId(userId string) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(usersCollectionName)

	filter := bson.M{"_id": userId}
	res := coll.FindOne(ctx, filter)
	var userData models.UserData
	err := res.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding finded user data in getUserData, error: %w", err)
	}

	return &userData, nil
}

func (m MongoDB) DeleteBook(tokenBook string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := m.DB.Collection(usersCollectionName)

	filter := bson.M{"fileToken": tokenBook}
	coll.FindOneAndDelete(ctx, filter)
	return nil
}

func NewClient(login, pwd, dbName, host, port string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/", login, pwd, host, port))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	log.Println("Client of mongodb successfully connected!")

	db := client.Database(dbName)
	log.Printf("Name of DB: %s", db.Name())
	return &MongoDB{DB: db}, nil
}
