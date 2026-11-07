package pgerror

import "github.com/jackc/pgx/v5/pgconn"

type PgError struct {
	Code           Code
	ConstraintName string
	Err            error
}

func HandlePgError(pgConnErr pgconn.PgError, pgErrors []PgError) error {
	for _, pgError := range pgErrors {
		if pgConnErr.Code == pgError.Code.String() && pgConnErr.ConstraintName == pgError.ConstraintName {
			return pgError.Err
		}
	}

	return nil
}
