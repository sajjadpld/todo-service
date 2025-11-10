package utils

import (
	"microservice/pkg/validator"
	"strconv"
	"time"
)

// StrToTimestamp convert numeric string to time format
func StrToTimestamp(t string) (timestamp time.Time, err error) {
	if err = validator.Var(t, "numeric"); err != nil {
		return
	}

	converted, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}

	timestamp = time.Unix(converted, 0)
	return
}
