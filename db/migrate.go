package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"testing"

	_ "embed"

	"go.uber.org/zap"

	"github.com/pressly/goose/v3"

	"github.com/Max-Gabriel-Susman/bestir-go-kit/bestirlog"
)

var (
	//go:embed version.txt
	version string
	//go:embed migrations/*.sql
	embedMigrations embed.FS

	src = MigrationSource{
		Migrations: embedMigrations,
		Dir:        "migrations", // embedded directory name
	}

	// Parse version from version.txt as an int64
	DesiredVersion = func() int64 {
		i, err := strconv.ParseInt(strings.Trim(version, "\n"), 10, 64)
		if err != nil {
			panic(fmt.Errorf("Unable to parse version from version.txt : %w", err))
		}
		return i
	}()
)

// EnsureMigrations is a simple wrapper function meant to be called by main.go
func EnsureMigrations(ctx context.Context, lgr *bestirlog.ZapLogger, cfg Config) error {
	fmt.Println("ensure migrations invoked")
	migrationLogger := lgr.Logger.WithOptions(zap.AddCallerSkip(-1))

	// Override Gooses default logger with a copy of our existing logger
	cfg.LoggerOverride = WrapLogger(migrationLogger)

	return Migrate(ctx, cfg, src, DesiredVersion)
}

// MigrationSource wraps the embedded migrations filesystem and the directory name (its a convenience thing)
type MigrationSource struct {
	Migrations embed.FS
	Dir        string // Name of the embedded migrations directory
}

// Config holds the information needed to run the migrations
type Config struct {
	User             string
	Password         string
	Host             string
	Port             string
	Name             string
	AdditionalParams map[string]string

	LoggerOverride goose.Logger
	DryRun         bool
}

// Migrate connects to the database described in the Config,
// and migrates it up or down based on if the current database version is different from the passed
// in desiredVersion
func Migrate(ctx context.Context, cfg Config, src MigrationSource, desiredVersion int64) error {
	fmt.Println("ensure migrate")
	dsn := func() string {
		q := make(url.Values)
		for k, v := range cfg.AdditionalParams {
			q.Add(k, v)
		}
		// Note: for MySQL parseTime flag must be enabled.
		q.Set("parseTime", "true")
		// Note: for MySQL multiStatements must be enabled.
		//	This is required when writing multiple queries separated by ';' characters in a single sql file.
		q.Set("multiStatements", "true")

		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, q.Encode())
	}()
	fmt.Println(fmt.Sprintf("dsn is: %s", dsn))

	db, err := goose.OpenDBWithDriver("mysql", dsn) // This is just runs sql.Open and tells goose which driver to use
	if err != nil {
		return fmt.Errorf("unable to open db for migrating: %w", err)
	}

	// Tell Goose which filesystem to use
	goose.SetBaseFS(src.Migrations)

	// This isn't very pretty, but it works.
	var l goose.Logger
	if cfg.LoggerOverride != nil {
		l = cfg.LoggerOverride
		goose.SetLogger(l)
	} else {
		l = log.Default()
	}

	// you now have l to use to log any messages from here on (it will show up inline with goose log messages)

	// advisory lock to prevent multiple instances from attempting to migrate at the same time
	// locker := gomysqllock.NewMysqlLocker(db)
	// fmt.Println("prior to lock assignment")
	// // ObtainTimeoutContext tries to acquire lock and gives up when the given context is cancelled or 60 seconds has passed
	// lock, err := locker.ObtainTimeoutContext(ctx, "goose-migration", 60)
	// if err != nil {
	// 	return err
	// }
	// defer func() {
	// 	// be sure to release the lock when we are done
	// 	if err := lock.Release(); err != nil {
	// 		l.Fatal("Error Releasing Lock:", err)
	// 	}
	// }()
	// fmt.Println("post lock assignment")

	// grab current migration version of the database
	dbVersion, err := goose.EnsureDBVersion(db)
	if err != nil {
		return fmt.Errorf("Unable to get current db version: %w", err)
	}

	if cfg.DryRun {
		l.Print("DRY RUN MODE: Current DB Version: ", dbVersion, "  Desired DB Version: ", desiredVersion)
		return nil
	}

	if dbVersion > desiredVersion {
		// migrate down to desired version
		return goose.DownTo(db, src.Dir, desiredVersion)
	}

	fmt.Println("ensure migrations reached final return stmt")
	// migrate up to desired version
	// if dbVersion == desiredVersion, goose will log "no migration needed"
	return goose.UpTo(db, src.Dir, desiredVersion)
}

// GooseLogger is intended to wrap the services regular logger so that the migrations show up correctly in datadog
type GooseLogger struct {
	*zap.SugaredLogger
}

// https://docs.datadoghq.com/getting_started/tagging/unified_service_tagging/
func WrapLogger(z *zap.Logger) *GooseLogger {
	return &GooseLogger{SugaredLogger: z.Sugar()}
}

func (l *GooseLogger) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *GooseLogger) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *GooseLogger) Printf(format string, v ...interface{}) {
	t := strings.Trim(fmt.Sprintf(format, v...), "\n") // hack to rm trailing new line char
	l.Info(t)
}

// TestingLogger is intended to be used while testing, so that the migration logs show up as test logs instead of rando std.out logs
type TestingLogger struct {
	*testing.T
}

func (t TestingLogger) Print(v ...interface{})                 { t.Log(v...) }
func (t TestingLogger) Println(v ...interface{})               { t.Log(v...) }
func (t TestingLogger) Printf(format string, v ...interface{}) { t.Logf(format, v...) }

// TestEnsureMigrations is a function ment to be called from a test case.
// It allows tests to access to the embedded data (version and src)
// This allows us to use the real migrations in our tests.
// (Not intended for use in the migrate_test.go file, but in the ./testing directory)
func TestEnsureMigrations(t *testing.T, cfg Config) {
	t.Helper()
	cfg.LoggerOverride = TestingLogger{T: t}
	if err := Migrate(context.Background(), cfg, src, DesiredVersion); err != nil {
		t.Fatal(err)
	}
}
