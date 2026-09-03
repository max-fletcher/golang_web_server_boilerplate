package validator

import (
	"net/url"
	"strconv"
)

type ValidatedPaginationQSData struct {
	FilterString string
	Page         int
	Limit        int
	Offset       int
}

func ValidatePaginationQS(query url.Values) (ValidatedPaginationQSData, map[string]string) {
	page := 1
	limit := 10
	errors := make(map[string]string)

	// Filter string
	filterString := query.Get("filter")

	// Page
	const MaxPage = 1000000
	if value := query.Get("page"); value != "" {
		pageValue, err := strconv.Atoi(value)

		if err != nil {
			errors["page"] = "Page must be a number"
		} else if pageValue >= MaxPage {
			errors["page"] = "Page value out of range"
		} else if pageValue < 1 {
			errors["page"] = "Page must be greater than 0"
		} else {
			page = pageValue
		}
	}

	// Limit
	const MaxLimit = 100
	if value := query.Get("limit"); value != "" {
		limitValue, err := strconv.Atoi(value)

		if err != nil {
			errors["limit"] = "Limit must be a number"
		} else if limitValue >= MaxLimit {
			errors["limit"] = "Limit cannot be greater than 100"
		} else if limitValue < 1 {
			errors["limit"] = "Limit must be greater than 0"
		} else {
			limit = limitValue
		}
	}

	return ValidatedPaginationQSData{
		FilterString: filterString,
		Page:         page,
		Limit:        limit,
		Offset:       (page - 1) * limit,
	}, errors
}
