package role

import (
	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/database"
	"github.com/gocraft/dbr/v2"
)

func NewMySQLStore(conn *dbr.Connection) *MySQLStorage {
	return &MySQLStorage{conn: conn, sess: conn.NewSession(nil)}
}

type MySQLStorage struct {
	conn *dbr.Connection
	sess *dbr.Session
}

var (
	permissionTable = database.NewTable("role", Role{})
)
