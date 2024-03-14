package pgsql_utils

import (
	"log"
	"strings"

	"github.com/idprm/go-linkit-tsel/internal/utils/errors"
	pq "github.com/lib/pq"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*pq.Error)
	log.Println(sqlErr)
	log.Println(ok)

	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("No record matching given id")
		}
		return errors.NewInternalServerError("Error parsing database response")
	}

	switch sqlErr.Code.Name() {
	case "02001":
		return errors.NewBadRequestError("Invalid data")
	}

	return errors.NewInternalServerError("Error processing request")
}
