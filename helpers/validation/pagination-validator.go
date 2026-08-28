package validator

import (
	"net/url"
	"strconv"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

func BaseValidatePagination(query url.Values) (Pagination, map[string]string) {
	page := 1
	limit := 10
	errors := make(map[string]string)

	// Page
	if value := query.Get("page"); value != "" {
		pageValue, err := strconv.Atoi(value)

		if err != nil {
			errors["page"] = "Page must be a number"
		} else if pageValue < 1 {
			errors["page"] = "Page must be greater than 0"
		} else {
			page = pageValue
		}
	}

	// Limit
	if value := query.Get("limit"); value != "" {
		limitValue, err := strconv.Atoi(value)

		if err != nil {
			errors["limit"] = "Limit must be a number"
		} else if limitValue < 1 {
			errors["limit"] = "Limit must be greater than 0"
		} else if limitValue > 100 {
			errors["limit"] = "Limit cannot be greater than 100"
		} else {
			limit = limitValue
		}
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}, errors
}
