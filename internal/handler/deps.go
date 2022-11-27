package handler

import (
	"database/sql"
)

type Deps struct {
	// Logger // must have do eet
	// Conn *pgx.Conn
	DB *sql.DB
}
