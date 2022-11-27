package role

type API struct {
	// Logger *bestirlog.Logger
	// Store CockroachDBStorage // we'll do cockroach l8r
	Store *MySQLStorage
}

// we may want to parameterize storage and logging later
// func NewAPI(conn *pgx.Conn) *API {
func NewAPI(store *MySQLStorage) *API {
	return &API{
		Store: store,
	}
}

// return &API{Store: *NewCockroachDBStorage(conn)} // we'll do cockroach l8r
