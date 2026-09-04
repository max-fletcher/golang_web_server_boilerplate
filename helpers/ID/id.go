package id_helpers

import (
	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
)

func ParseUUIDRouteParam(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, requests.ErrInvalidUUIDParam{
			Param: "user id",
			Err:   err,
		}
	}

	return id, nil
}

func ParseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, requests.ErrInvalidUUID{
			Field: "user id",
			Err:   err,
		}
	}

	return id, nil
}
