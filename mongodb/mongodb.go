package mongodb

import (
	"context"
	"crud-books/config"
	"crud-books/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollectionName = "users"
	booksCollectionName = "books"
	filesCollectionName = "files"
	ctxTimeout          = 3
)

type MongoDB struct {
	booksCollection mongo.Collection
	usersCollection mongo.Collection
	filesCollection mongo.Collection
}

type InsertedID struct {
	InsertedID interface{}
}

func (m *MongoDB) GetBook(bookToken string) (*models.BookData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"fileToken": bookToken}
	result := m.booksCollection.FindOne(ctx, filter)
	var bookData models.BookData
	err := result.Decode(&bookData)
	if err != nil {
		return nil, fmt.Errorf("decode book data in get book func of db failed, error: %w", err)
	}
	return &bookData, nil
}

func (m *MongoDB) CreateUser(email, passwordHash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	res := m.usersCollection.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return "", fmt.Errorf("user with email %v already created", email)
	}

	doc := models.UserData{Email: email, PasswordHash: passwordHash, BooksIds: []string{}}

	cur, err := m.usersCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed collection insert in add user func, error: %w", err)
	}
	objId := fmt.Sprintf("%v", cur.InsertedID)
	return objId, nil
}

func (m *MongoDB) CreateBook(title, description, fileToken, emailOwner string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
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

	insertResult, err := m.booksCollection.InsertOne(ctx, dBook)
	if err != nil {
		return "", fmt.Errorf("failed collection insert in create book func, error: %w", err)
	}

	filter := bson.M{"email": emailOwner}
	update := bson.M{"$set": bson.M{"$push": bson.M{
		"books": insertResult.InsertedID,
	}}}

	m.usersCollection.FindOneAndUpdate(ctx, filter, update)

	return dBook.FileToken, nil
}

func (m *MongoDB) GetListBooksPublic(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeout)
	defer cancel()

	var filter bson.M
	modelForSearch := mongo.IndexModel{Keys: bson.M{"title": "text"}}
	indexName, err := m.booksCollection.Indexes().CreateOne(ctx, modelForSearch)
	if err != nil {
		return nil, fmt.Errorf("creating index in get list books of user failed, error: %w", err)
	}
	if paramsOfBooks.Search != "" {
		filter = bson.M{"$text": bson.M{"$search": paramsOfBooks.Search}}
	}

	limitOpt := options.Find().SetLimit(int64(paramsOfBooks.Limit))
	sortOpt := options.Find().SetSort(bson.M{paramsOfBooks.SortField: paramsOfBooks.Direction})
	cursor, err := m.booksCollection.Find(ctx, filter, sortOpt, limitOpt)
	if err != nil {
		return nil, fmt.Errorf("find in books collection failed, error: %w", err)
	}

	var res []models.BookData
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, fmt.Errorf("decoding cursor failed, error: %w", err)
	}
	_, err = m.booksCollection.Indexes().DropOne(context.TODO(), indexName)
	if err != nil {
		return nil, fmt.Errorf("dropping index failed, error: %w", err)
	}
	return &res, nil
}

func (m *MongoDB) GetListBooksOfUser(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeout)
	defer cancel()

	emailFilter := bson.M{"email": paramsOfBooks.Email}

	modelForSearch := mongo.IndexModel{Keys: bson.M{"title": "text"}}
	indexName, err := m.booksCollection.Indexes().CreateOne(ctx, modelForSearch)
	if err != nil {
		return nil, fmt.Errorf("creating index in get list books of user failed, error: %w", err)
	}

	userCur, err := m.usersCollection.Find(ctx, emailFilter)
	if err != nil {
		return nil, fmt.Errorf("user search failed, error: %w", err)
	}
	var userData models.UserData
	err = userCur.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("decoding user data failed, error: %w", err)
	}

	bookFilter := bson.M{"owner": userData.Id, "$text": bson.M{"$search": paramsOfBooks.Search}}

	sortOpt := options.Find().SetSort(bson.M{paramsOfBooks.SortField: paramsOfBooks.Direction})
	limitOpt := options.Find().SetLimit(int64(paramsOfBooks.Limit))
	booksCur, err := m.booksCollection.Find(ctx, bookFilter, sortOpt, limitOpt)
	if err != nil {
		return nil, fmt.Errorf("book collection search failed, error: %w", err)
	}
	var books *[]models.BookData
	err = booksCur.All(ctx, &books)
	if err != nil {
		return nil, fmt.Errorf("decoding cursor failed, error: %w", err)
	}
	_, err = m.booksCollection.Indexes().DropOne(context.TODO(), indexName)
	if err != nil {
		return nil, fmt.Errorf("dropping index failed, error: %w", err)
	}
	return books, nil
}

func (m *MongoDB) ChangeFieldOfBook(id, fieldName, fieldValue string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{fieldName: fieldValue}}

	m.booksCollection.FindOneAndUpdate(ctx, filter, update)

	return nil
}

func (m *MongoDB) UploadFileData(fileToken, downloadPage string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	doc := bson.M{"token": fileToken, "downloadPage": downloadPage}

	_, err := m.filesCollection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insertone in uploadfiledata failed, error: %w", err)
	}
	return nil
}

func (m *MongoDB) GetFileData(fileToken string) (*models.FileData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"token": fileToken}

	res := m.filesCollection.FindOne(ctx, filter)
	var f models.FileData
	err := res.Decode(&f)
	if err != nil {
		return nil, fmt.Errorf("decoding result of getfiledata failed, error: %w", err)
	}
	return &models.FileData{Id: f.Id, Token: f.Token, DownloadPage: f.DownloadPage}, nil
}

func (m *MongoDB) GetUserData(email string) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	res := m.usersCollection.FindOne(ctx, filter)
	var userData models.UserData
	err := res.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding finded user data in getUserData, error: %w", err)
	}

	return &userData, nil
}

func (m *MongoDB) GetUserDataByInsertedId(userId string) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"_id": userId}
	res := m.usersCollection.FindOne(ctx, filter)
	var userData models.UserData
	err := res.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("failed while decoding finded user data in getUserData, error: %w", err)
	}

	return &userData, nil
}

func (m MongoDB) DeleteBook(tokenBook string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	filter := bson.M{"fileToken": tokenBook}
	m.usersCollection.FindOneAndDelete(ctx, filter)
	return nil
}

func pingToDb(ctx context.Context, db mongo.Database) string {
	var result bson.M
	err := db.RunCommand(context.Background(), bson.M{"ping": 1}).Decode(&result)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return ""
}

func NewClient(cfg config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    cfg.DatabaseName,
		Username:      cfg.DatabaseLogin,
		Password:      cfg.DatabasePwd,
	}
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", cfg.DatabaseHost, cfg.DatabasePort)).SetAuth(credential)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DatabaseName)

	ctxPing, cancelPing := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelPing()

	pingSummary := pingToDb(ctxPing, *db)
	if pingSummary != "" {
		return nil, fmt.Errorf("ping error: %s", pingSummary)
	}

	return &MongoDB{
		booksCollection: *db.Collection(booksCollectionName),
		usersCollection: *db.Collection(usersCollectionName),
		filesCollection: *db.Collection(filesCollectionName),
	}, nil
}
