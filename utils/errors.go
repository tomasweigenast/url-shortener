package utils

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgErrorCodeUniqueValidation = "23505"
)

var (
	ErrAlreadyExists = errors.New("already-exists")
	ErrNotFound      = errors.New("not-found")
)

func TransformPgError(pgerror error) error {
	if errors.Is(pgerror, pgx.ErrNoRows) {
		return ErrNotFound
	}

	if err, ok := pgerror.(*pgconn.PgError); ok {
		switch err.Code {
		case pgErrorCodeUniqueValidation:
			return ErrAlreadyExists
		}
	}

	return pgerror
}
