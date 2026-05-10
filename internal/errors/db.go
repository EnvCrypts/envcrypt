package errors

import (
	"database/sql"
	"errors"
)

func FromDB(err error, notFoundMsg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return NotFound(notFoundMsg, "")
	}
	return Internal(err)
}
