package database

import (
	"github.com/gocraft/dbr/v2"
)

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(dbr.SessionRunner) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
//
//nolint:errcheck
func WithTransaction(sess *dbr.Session, fn TxFn) (err error) {
	tx, err := sess.Begin()
	if err != nil {
		return
	}

	defer func() {
		switch p := recover(); p {
		case p != nil:
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		case err != nil:
			// something went wrong, rollback
			tx.Rollback()
		default:
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
