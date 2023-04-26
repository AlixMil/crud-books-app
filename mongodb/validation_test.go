package mongodb

import (
	"crud-books/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateParametersFunc(t *testing.T) {
	filter := models.Filter{
		Email:  "jjojas@gmail.com",
		Search: "1asdagasd",
	}
	sort := models.Sort{
		SortField: "title",
		Direction: "asc",
		Limit:     12,
		Offset:    12,
	}
	t.Run("complete", func(t *testing.T) {
		params := ValidateParams(filter, sort)
		assert.Equal(t, filter.Email, params.Email)
		assert.Equal(t, sort.SortField, params.SortField)
		assert.Equal(t, sort.Limit, params.Limit)
	})
	t.Run("direction_is_default", func(t *testing.T) {
		filter = models.Filter{
			Email:  "aksldkals@gmail.com",
			Search: "",
		}
		sort = models.Sort{
			SortField: "sdasd",
			Limit:     1,
			Direction: "asdfasa",
			Offset:    0,
		}
		params := ValidateParams(filter, sort)

		assert.Equal(t, directionDefaultParam, params.Direction)
	})

	t.Run("SortField_is_default", func(t *testing.T) {
		filter = models.Filter{
			Email:  "aksldkals@gmail.com",
			Search: "",
		}
		sort = models.Sort{
			SortField: "sdasd",
			Limit:     1,
			Direction: "asdfasa",
			Offset:    0,
		}

		params := ValidateParams(filter, sort)

		assert.Equal(t, sortFieldDefaultParam, params.SortField)
	})

	t.Run("limit_is_default", func(t *testing.T) {
		filter = models.Filter{
			Email:  "aksldkals@gmail.com",
			Search: "",
		}
		sort = models.Sort{
			SortField: "sdasd",
			Limit:     10000,
			Direction: "asdfasa",
			Offset:    0,
		}

		params := ValidateParams(filter, sort)

		assert.Equal(t, maxSizeOfLimitParam, params.Limit)
	})
}
