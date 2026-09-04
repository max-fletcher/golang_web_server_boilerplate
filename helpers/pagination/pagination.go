package pagination

import (
	validation "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
)

type PaginatedData = struct {
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
	Total   int  `json:"total"`
	Next    bool `json:"next"`
	Prev    bool `json:"prev"`
	Data    any  `json:"data"`
}

func GeneratePaginationFormat(pagination validation.ValidatedPaginationQSData, total int, data any) PaginatedData {
	next := pagination.Offset+pagination.Limit < total
	prev := pagination.Offset > 0 && pagination.Offset < total

	return PaginatedData{
		Page:    pagination.Page,
		PerPage: pagination.Limit,
		Total:   total,
		Next:    next,
		Prev:    prev,
		Data:    data,
	}
}
