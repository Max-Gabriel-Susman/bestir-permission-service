package permission

import (
	"context"

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
	permissionTable = database.NewTable("permission", permission{})
)

func (s *MySQLStorage) Listpermissions(ctx context.Context) ([]permission, error) {
	query := s.sess.Select(permissionTable.Columns...).
		From(permissionTable.Name)

	permissions := []permission{}

	if _, err := query.LoadContext(ctx, &permissions); err != nil {
		return permissions, database.ClassifyError(err)
	}

	return permissions, nil
}

func (s *MySQLStorage) getpermissionByIdempotencyKey(ctx context.Context, idempotencyKey string) (permission, error) {
	var permission permission
	err := s.sess.Select(permissionTable.Columns...).
		From(permissionTable.Name).
		Where("idempotency_key = ?", idempotencyKey).
		LoadOneContext(ctx, &permission)
	return permission, database.ClassifyError(err)
}

func (s *MySQLStorage) Createpermission(ctx context.Context, permission permission) error {
	_, err := s.sess.InsertInto(permissionTable.Name).
		Columns(permissionTable.Columns...).
		Record(permission).
		ExecContext(ctx)
	return database.ClassifyError(err)
}

func (s *MySQLStorage) Deletepermission(ctx context.Context, permission permission) error {
	_, err := s.sess.InsertInto(permissionTable.Name).
		Columns(permissionTable.Columns...).
		Record(permission).
		ExecContext(ctx)
	return database.ClassifyError(err)
}

func (s *MySQLStorage) Updatepermission(ctx context.Context, permission permission) error {
	_, err := s.sess.InsertInto(permissionTable.Name).
		Columns(permissionTable.Columns...).
		Record(permission).
		ExecContext(ctx)
	return database.ClassifyError(err)
}
