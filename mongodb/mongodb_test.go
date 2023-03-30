package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ValidateFunc(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "title"
		direction := "asc"
		limit := 12

		params, err := GetParamsWValidate(email, search, sort, direction, limit)
		require.NoError(t, err)
		assert.Equal(t, params.Email, email)
		assert.Equal(t, params.SortField, sort)
		assert.Equal(t, params.Limit, limit)
	})
	t.Run("error_direction", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "title"
		direction := "dfajks"
		limit := 12

		_, err := GetParamsWValidate(email, search, sort, direction, limit)

		assert.EqualError(t, err, "provided direction parameter incorrect")
	})

	t.Run("error_SortField", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "ajskda"
		direction := "asc"
		limit := 12

		_, err := GetParamsWValidate(email, search, sort, direction, limit)

		assert.EqualError(t, err, "provided sorting parameter incorrect")
	})

	t.Run("error_limit", func(t *testing.T) {
		email := "jjojas@gmail.com"
		search := "1asdagasd"
		sort := "date"
		direction := "asc"
		limit := 12343

		_, err := GetParamsWValidate(email, search, sort, direction, limit)

		assert.EqualError(t, err, "you should provide limit parameter in range between from 5 to 100")
	})
}
