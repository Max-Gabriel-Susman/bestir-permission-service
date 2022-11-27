package database

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/bestirerror"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
)

func ErrNotFound(err error) error {
	return bestirerror.WithStatusCode(err, http.StatusNotFound)
}

func ErrDuplicateUnique(err error) error {
	return bestirerror.WithCodeAndMessage(err, http.StatusConflict, err.Error())
}

func ErrForeignKeyConstraint(err error) error {
	return bestirerror.WithCodeAndMessage(err, http.StatusBadRequest, "foreign key constraint fails")
}

func ClassifyError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound(err)
	}
	if errors.Is(err, dbr.ErrNotFound) {
		return ErrNotFound(err)
	}
	if mse := new(mysql.MySQLError); errors.As(err, &mse) {
		if mse.Number == uint16(1062) {
			// Error 1062: Duplicate entry 'ch_123' for key 'stripe_charge_id'
			return ErrDuplicateUnique(err)
		}
		if mse.Number == uint16(1452) {
			// Error 1452: Cannot add or update a child row: a foreign key constraint fails
			return ErrForeignKeyConstraint(err)
		}
	}

	return err
}
