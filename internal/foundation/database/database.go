// Package database provides support for accessing the database.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Config is the required properties to use the database.
type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	Params   string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config, serviceName string) (*sql.DB, error) {
	sqltrace.Register("mysql", &mysql.MySQLDriver{},
		sqltrace.WithServiceName(serviceName),
		sqltrace.WithAnalytics(true),
	)
	dsn := DSN(cfg)
	log.Println(fmt.Sprintf("dsn is: %s", dsn))
	return sqltrace.Open("mysql", dsn)
}

func NewDBR(db *sql.DB) *dbr.Connection {
	return &dbr.Connection{DB: db, EventReceiver: &dbr.NullEventReceiver{}, Dialect: dialect.MySQL}
}

func DSN(cfg Config) string {
	// return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Name, cfg.Params)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Name)
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.SetTag(ext.ManualDrop, true)
	}

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
