package common_errors

import (
	"errors"

	"github.com/lib/pq"
)

func GetPostgresError(err error) *pq.Error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return nil
	}

	return pqErr
}
