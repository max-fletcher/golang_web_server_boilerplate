package pagination

import (
	validation "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
)

type PaginatedData = struct {
	Page  int
	Limit int
	Total int
	Next  bool
	Prev  bool
	Data  any
}

func GeneratePaginationFormat(pagination validation.ValidatedPaginationQSData, total int, data any) PaginatedData {
	next := pagination.Offset+pagination.Limit < total
	prev := pagination.Offset > 0 && pagination.Offset < total+pagination.Limit

	return PaginatedData{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Total: total,
		Next:  next,
		Prev:  prev,
		Data:  data,
	}
}
