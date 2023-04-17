package mongodb

import (
	"context"
	"crud-books/config"
	"crud-books/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	usersCollectionName = "users"
	booksCollectionName = "books"
	filesCollectionName = "files"
	ctxTimeout          = 3
)

type MongoDB struct {
	booksCollection *mongo.Collection
	usersCollection *mongo.Collection
	filesCollection *mongo.Collection
	db              *mongo.Database
}

func (m *MongoDB) GetBook(bookToken string) (*models.BookData, error) {
	filter := bson.M{"fileToken": bookToken}
	result := m.booksCollection.FindOne(context.TODO(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	var bookData models.BookData
	err := result.Decode(&bookData)
	if err != nil {
		return nil, fmt.Errorf("decode book data in get book func of db failed, error: %w", err)
	}

	return &bookData, nil
}

func (m *MongoDB) getUserData(filter bson.M) (*models.UserData, error) {
	userCur := m.usersCollection.FindOne(context.TODO(), filter)

	if userCur.Err() == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	userData := new(models.UserData)
	err := userCur.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("userdata decoding failed, error: %w", err)
	}

	return userData, nil
}

func (m *MongoDB) CreateUser(email, passwordHash string) (string, error) {
	filter := bson.M{"email": email}
	_, err := m.getUserData(filter)

	if err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("user with email %v already created", email)
	}

	doc := bson.M{"email": email, "passwordHash": passwordHash, "booksIds": []bson.M{}}

	cur, err := m.usersCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		return "", fmt.Errorf("failed collection insert in add user func, error: %w", err)
	}

	id, ok := cur.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("converting insertedid of user in db failed, err: %w", err)
	}

	return id.Hex(), nil
}

func (m *MongoDB) CreateBook(title, description, fileToken, emailOwner string) (string, error) {
	_, err := m.GetBook(fileToken)
	if err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("book from file token already exist")
	}

	fileData, err := m.GetFileData(fileToken)
	if err != nil {
		return "", fmt.Errorf("get file data failed, error: %w", err)
	}

	dBook := models.BookData{
		Title:       title,
		Description: description,
		FileToken:   fileToken,
		OwnerEmail:  emailOwner,
		Url:         fileData.DownloadPage,
	}

	insertResult, err := m.booksCollection.InsertOne(context.TODO(), dBook)
	if err != nil {
		return "", fmt.Errorf("failed collection insert in create book func, error: %w", err)
	}

	filter := bson.M{"email": emailOwner}
	update := bson.M{"$set": bson.M{"$push": bson.M{
		"books": insertResult.InsertedID,
	}}}

	m.usersCollection.FindOneAndUpdate(context.TODO(), filter, update)

	return dBook.FileToken, nil
}

func prepareCollectionToSearch(collection *mongo.Collection) (string, error) {
	modelForSearch := mongo.IndexModel{Keys: bson.M{"title": "text"}}
	indexName, err := collection.Indexes().CreateOne(context.TODO(), modelForSearch)
	if err != nil {
		return "", fmt.Errorf("preparing collection to search failed, error: %w", err)
	}

	return indexName, nil
}

func getFindOptions(params models.ValidateDataInGetLists) *options.FindOptions {
	fOpt := options.Find()
	fOpt.SetLimit(int64(params.Limit))
	fOpt.SetSort(bson.M{params.SortField: params.Direction})
	fOpt.SetSkip(int64(params.Offset))

	return fOpt
}

func (m *MongoDB) GetListBooksPublic(paramsOfBooks *models.ValidateDataInGetLists) ([]models.BookData, error) {
	var filter bson.M

	indexName, err := prepareCollectionToSearch(m.booksCollection)
	if err != nil {
		return nil, err
	}

	opts := getFindOptions(*paramsOfBooks)

	if paramsOfBooks.Search != "" {
		filter = bson.M{"$text": bson.M{"$search": paramsOfBooks.Search}}
	}

	cursor, err := m.booksCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find in books collection failed, error: %w", err)
	}

	var res []models.BookData
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		return nil, fmt.Errorf("decoding cursor failed, error: %w", err)
	}
	_, err = m.booksCollection.Indexes().DropOne(context.TODO(), indexName)
	if err != nil {
		return nil, fmt.Errorf("dropping index failed, error: %w", err)
	}

	return res, nil
}

func (m *MongoDB) GetListBooksOfUser(paramsOfBooks *models.ValidateDataInGetLists) ([]models.BookData, error) {
	userData, err := m.GetUserData(paramsOfBooks.Email)
	if err != nil {
		return nil, err
	}

	var bookFilter bson.M

	bookFilter = bson.M{"owner": userData.Email, "$text": bson.M{"$search": paramsOfBooks.Search}}

	if paramsOfBooks.Search == "" {
		bookFilter = bson.M{"owner": userData.Email}
	}

	opts := getFindOptions(*paramsOfBooks)
	indexName, err := prepareCollectionToSearch(m.booksCollection)
	if err != nil {
		return nil, fmt.Errorf("prepare collection to search error: %w", err)
	}

	booksCur, err := m.booksCollection.Find(context.TODO(), bookFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("book collection search failed, error: %w", err)
	}

	var books []models.BookData

	err = booksCur.All(context.TODO(), &books)
	if err != nil {
		return nil, fmt.Errorf("decoding cursor failed, error: %w", err)
	}

	_, err = m.booksCollection.Indexes().DropOne(context.TODO(), indexName)
	if err != nil {
		return nil, fmt.Errorf("dropping index failed, error: %w", err)
	}

	return books, nil
}

func (m *MongoDB) UpdateBook(bookFileToken string, updater models.BookDataUpdater) error {
	_, err := m.GetFileData(updater.FileToken)
	if err != nil {
		return fmt.Errorf("filetoken should be only exist in system")
	}
	filter := bson.M{"fileToken": bookFileToken}
	update := bson.M{
		"$set": bson.M{
			"title":       updater.Title,
			"description": updater.Description,
			"fileToken":   updater.FileToken,
		}}

	res := m.booksCollection.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("book doesn't exist")
	}

	return nil
}

func (m *MongoDB) UploadFileData(fileData *models.FileData) error {
	doc := bson.M{
		"token":        fileData.Token,
		"downloadPage": fileData.DownloadPage,
	}

	_, err := m.filesCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		return fmt.Errorf("uploadFileData failed, error: %w", err)
	}

	return nil
}

func (m *MongoDB) GetFileData(fileToken string) (*models.FileData, error) {
	filter := bson.M{"token": fileToken}

	var f models.FileData
	err := m.filesCollection.FindOne(context.TODO(), filter).Decode(&f)
	if err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("decoding result of getfiledata failed, error: %w", err)
	}

	return &models.FileData{Id: f.Id, Token: f.Token, DownloadPage: f.DownloadPage}, nil
}

func (m *MongoDB) GetUserData(email string) (*models.UserData, error) {
	filter := bson.M{"email": email}

	userData, err := m.getUserData(filter)
	if err != nil {
		return nil, fmt.Errorf("getUserData error: %w", err)
	}

	log.Printf("userData: %s", userData)

	return userData, nil
}

func (m *MongoDB) GetUserDataByInsertedId(userId string) (*models.UserData, error) {

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("retrievev obj id from hex failed, error: %w", err)
	}
	filter := bson.M{"_id": objId}
	userData, err := m.getUserData(filter)
	if err != nil {
		return nil, fmt.Errorf("getUserData error: %w", err)
	}

	return userData, nil
}

func (m MongoDB) DeleteBook(bookId string) error {
	filter := bson.M{"fileToken": bookId}
	res := m.booksCollection.FindOneAndDelete(context.TODO(), filter)
	if res.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("deleting file failed, file doesn't exist")
	}

	return nil
}

func (m *MongoDB) Connect(cfg config.Config) error {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    cfg.DatabaseName,
		Username:      cfg.DatabaseLogin,
		Password:      cfg.DatabasePwd,
	}
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", cfg.DatabaseHost, cfg.DatabasePort))
	opts.SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("connecting to db failed, error: %w", err)
	}

	db := client.Database(cfg.DatabaseName)
	m.booksCollection = db.Collection(booksCollectionName)
	m.usersCollection = db.Collection(usersCollectionName)
	m.filesCollection = db.Collection(filesCollectionName)
	m.db = db

	return nil
}

func (m MongoDB) Ping() error {
	ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelPing()

	err := m.db.Client().Ping(ctxPing, readpref.Primary())

	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	return nil
}

func New() *MongoDB {
	return &MongoDB{}
}
