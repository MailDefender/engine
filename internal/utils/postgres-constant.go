package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolationErrCode = "23505"
)

func IsUniqueViolationErr(err error) bool {
	var perr *pgconn.PgError
	errors.As(err, &perr)
	res := perr.Code == UniqueViolationErrCode
	return res
}
