package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateParametersFunc(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "title"
		direction := "asc"
		limit := 12

		params := getParamsWValidate(email, search, sort, direction, limit)
		assert.Equal(t, params.Email, email)
		assert.Equal(t, params.SortField, sort)
		assert.Equal(t, params.Limit, limit)
	})
	t.Run("direction_is_default", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "title"
		direction := "dfajks"
		limit := 12

		params := getParamsWValidate(email, search, sort, direction, limit)

		assert.Equal(t, directionDefaultParam, params.Direction)
	})

	t.Run("SortField_is_default", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "ajskda"
		direction := "asc"
		limit := 12

		params := getParamsWValidate(email, search, sort, direction, limit)

		assert.Equal(t, sortFieldDefaultParam, params.SortField)
	})

	t.Run("limit_is_default", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "date"
		direction := "asc"
		limit := 12343

		params := getParamsWValidate(email, search, sort, direction, limit)

		assert.Equal(t, maxSizeOfLimitParam, params.Limit)
	})
}
