package mongodb

import (
	"context"
	"crud-books/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) GetBook(bookToken string) (*models.BookData, error) {
	searchByFileToken := bson.M{"fileToken": bookToken}
	result := m.booksCollection.FindOne(context.TODO(), searchByFileToken)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	var bookData models.BookData
	err := result.Decode(&bookData)
	if err != nil {
		return nil, fmt.Errorf("decoding book failed: %w", err)
	}

	return &bookData, nil
}

func (m *MongoDB) CreateUser(email, passwordHash string) (string, error) {
	doc := bson.M{"email": email, "passwordHash": passwordHash, "booksIds": []bson.M{}}

	cur, err := m.usersCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		return "", fmt.Errorf("insert user error: %w", err)
	}

	id, ok := cur.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("converting id error: %w", err)
	}

	return id.Hex(), nil
}

func (m *MongoDB) CreateBook(title, description, fileToken, emailOwner string) (string, error) {
	fileData, err := m.GetFileData(fileToken)
	if err != nil {
		return "", fmt.Errorf("getting file data error: %w", err)
	}

	dBook := models.BookData{
		Title:       title,
		Description: description,
		FileToken:   fileToken,
		OwnerEmail:  emailOwner,
		Url:         fileData.DownloadPage,
	}

	_, err = m.booksCollection.InsertOne(context.TODO(), dBook)
	if err != nil {
		return "", fmt.Errorf("inserting book error: %w", err)
	}

	return dBook.FileToken, nil
}

func (m *MongoDB) GetListBooks(filter models.Filter, sort models.Sort) ([]models.BookData, error) {

	params := ValidateParams(filter, sort)

	search := bson.M{}
	opts := getFindOptions(params)

	if params.Search != "" {
		search = bson.M{"$text": bson.M{"$search": params.Search}}
	}

	cursor, err := m.booksCollection.Find(context.TODO(), search, opts)
	if err != nil {
		return nil, fmt.Errorf("find books error: %w", err)
	}

	var res []models.BookData
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		return nil, fmt.Errorf("decoding bookData's error: %w", err)
	}

	return res, nil
}

func (m *MongoDB) GetListBooksUser(filter models.Filter, sort models.Sort) ([]models.BookData, error) {
	params := ValidateParams(filter, sort)

	userData, err := m.GetUserData(params.Email)
	if err != nil {
		return nil, err
	}

	var bookFilter bson.M

	bookFilter = bson.M{"owner": userData.Email, "$text": bson.M{"$search": params.Search}}

	if params.Search == "" {
		bookFilter = bson.M{"owner": userData.Email}
	}

	opts := getFindOptions(params)

	booksCur, err := m.booksCollection.Find(context.TODO(), bookFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("book find error: %w", err)
	}

	var books []models.BookData

	err = booksCur.All(context.TODO(), &books)
	if err != nil {
		return nil, fmt.Errorf("decoding bookData's error: %w", err)
	}

	return books, nil
}

func (m *MongoDB) UpdateBook(bookFileToken string, updater models.BookDataUpdater) error {
	_, err := m.GetFileData(updater.FileToken)
	if err != nil {
		return fmt.Errorf("filetoken should be only as existed fileTokens")
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
		return fmt.Errorf("upload file error: %w", err)
	}

	return nil
}

func (m *MongoDB) GetFileData(fileToken string) (*models.FileData, error) {
	filter := bson.M{"token": fileToken}

	var f models.FileData
	cur := m.filesCollection.FindOne(context.TODO(), filter)
	if cur.Err() != nil {
		return nil, fmt.Errorf("find files error: %w", cur.Err())
	}
	err := cur.Decode(&f)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("decoding filedata error: %w", err)
	}

	return &models.FileData{Id: f.Id, Token: f.Token, DownloadPage: f.DownloadPage}, nil
}

func (m *MongoDB) GetUserData(email string) (*models.UserData, error) {
	searchByEmail := bson.M{"email": email}

	userData, err := m.getUserData(searchByEmail)
	if err != nil {
		return nil, fmt.Errorf("getUserData error: %w", err)
	}

	return userData, nil
}

func (m *MongoDB) GetUserDataById(userId string) (*models.UserData, error) {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("convert userId from hex failed: %w", err)
	}
	searchById := bson.M{"_id": objId}
	userData, err := m.getUserData(searchById)
	if err != nil {
		return nil, fmt.Errorf("getting uData error: %w", err)
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

func getFindOptions(params *models.ParamsAfterValidation) *options.FindOptions {
	fOpt := options.Find()
	fOpt.SetLimit(int64(params.Limit))
	fOpt.SetSort(bson.M{params.SortField: params.Direction})
	fOpt.SetSkip(int64(params.Offset))

	return fOpt
}
