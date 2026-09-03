package numbers

import (
	"fmt"
	"math"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

func BigInt64ToInt(numInt64 int64) (int, error) {
	if numInt64 > int64(math.MaxInt) || numInt64 < int64(math.MinInt) {
		err := common_errors.ErrBigInt64ToIntError{
			Err: fmt.Errorf("Error: Value(%v) overflows standard int capacity", &numInt64),
		}
		return 0, err
	}

	safeInt := int(numInt64)
	return safeInt, nil
}
