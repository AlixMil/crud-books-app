package mongodb

import (
	"crud-books/models"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_GetBook(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("getBookSuccess", func(m *mtest.T) {
		bookColl := m.Client.Database("crudbooks").Collection(booksCollectionName)
		book := models.BookData{
			Id:          "123",
			Title:       "TITLE",
			Description: "DESCRIPTION",
			FileToken:   "FILETOKEN",
			Url:         "URL",
			OwnerEmail:  "OWNER",
		}

		m.AddMockResponses(mtest.CreateCursorResponse(1, "crubooks.books", mtest.FirstBatch, bson.D{
			{"_id", book.Id},
			{"title", book.Title},
			{"description", book.Description},
			{"fileToken", book.FileToken},
			{"url", book.Url},
			{"owner", book.OwnerEmail},
		}))
		db := New()
		db.booksCollection = bookColl

		res, err := db.GetBook("123")
		fmt.Println(res)
		assert.Equal(t, book.Title, res.Title)
		require.NoError(m, err)

		assert.Equal(m, true, reflect.DeepEqual(book, *res))
	})
}
