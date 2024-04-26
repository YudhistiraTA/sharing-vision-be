package response

import (
	"database/sql"
	"errors"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

var (
	ErrNotFound       = sql.ErrNoRows
	ErrInvalidRequest = errors.New("invalid request")
)

func ErrorCode(err error) int {
	switch err {
	case ErrNotFound:
		return 404
	case ErrInvalidRequest:
		return 400
	default:
		return 500
	}
}
