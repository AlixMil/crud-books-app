package mongodb

import "crud-books/models"

const (
	sortFieldDefaultParam = "title"
	directionDefaultParam = 1
	limitDefaultParam     = 10
	maxSizeOfLimitParam   = 100
)

func ValidateParams(filter models.Filter, sort models.Sort) *models.ParamsAfterValidation {
	p := new(models.ParamsAfterValidation)
	p.Email = filter.Email
	p.Search = filter.Search
	p.Limit = sort.Limit
	p.Offset = sort.Offset

	if sort.SortField != "title" && sort.SortField != "date" {
		p.SortField = sortFieldDefaultParam
	} else if sort.SortField == "title" {
		p.SortField = "title"
	} else if sort.SortField == "date" {
		p.SortField = "date"
	}

	if sort.Direction != "desc" && sort.Direction != "asc" {
		p.Direction = directionDefaultParam
	} else if sort.Direction == "asc" {
		p.Direction = 1
	} else if sort.Direction == "desc" {
		p.Direction = -1
	}

	if sort.Limit > maxSizeOfLimitParam {
		p.Limit = maxSizeOfLimitParam
	}

	if sort.Limit == 0 || sort.Limit < 0 {
		p.Limit = limitDefaultParam
	}
	return p
}
