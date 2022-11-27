package db_test

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	migrate "github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/ory/dockertest"
	"github.com/pressly/goose/v3"
)

var (

	//go:embed migrations/*.sql
	embedMigrations embed.FS

	src = migrate.MigrationSource{
		Migrations: embedMigrations,
		Dir:        "migrations",
	}

	port string
	db   *sql.DB // Explicitly for  getting db version w/ goose
)

// Taken from: https://github.com/ory/dockertest#using-dockertest
func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	// Tell docker to hard kill the container in 600 seconds
	if err := resource.Expire(600); err != nil {
		log.Fatalf("Could not set resource expiration time: %s", err)
	}

	// exponential backoff-retry, because the permission in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		port = resource.GetPort("3306/tcp")
		db, err = goose.OpenDBWithDriver("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true&multiStatements=true", port))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// Test Migrate Up
func TestMigrate(t *testing.T) {
	t.Parallel()
	testCFG := migrate.Config{
		User:     "root",
		Password: "secret",
		Host:     "localhost",
		Port:     port,
		Name:     "mysql",
	}
	// The Simple test cases simply migrates the database up to the version set in version.txt, and back down to zero.
	// This also allows you to ensure your down migrations are valid
	t.Run("Simple", func(t *testing.T) {
		t.Run("Up", func(t *testing.T) {
			// Migrate all the way up to the version in version.txt
			migrate.TestEnsureMigrations(t, testCFG)
		})
		t.Run("Down", func(t *testing.T) {
			// Migrate down to 0
			testMigrate(t, testCFG, 0)
		})
	})

	var expectedCurrentVersion int
	t.Run("UpTo16", func(t *testing.T) {
		wantVersion := 16
		t.Run("DRY RUN", func(t *testing.T) {
			testCFG.DryRun = true
			testMigrate(t, testCFG, wantVersion)
			testCFG.DryRun = false

			// Ensure Current db version wasn't changed
			diff(t, expectedCurrentVersion, getDBVersion(t))
		})
		t.Run("WET RUN", func(t *testing.T) {
			testMigrate(t, testCFG, wantVersion)

			// Ensure Current db version was changed
			diff(t, wantVersion, getDBVersion(t))
			expectedCurrentVersion = wantVersion
		})
	})

	t.Run("NoChange", func(t *testing.T) {
		wantVersion := expectedCurrentVersion
		t.Run("DRY RUN", func(t *testing.T) {
			testCFG.DryRun = true
			testMigrate(t, testCFG, wantVersion)
			testCFG.DryRun = false

			// Ensure Current db version wasn't changed
			diff(t, expectedCurrentVersion, getDBVersion(t))
		})
		t.Run("WET RUN", func(t *testing.T) {
			testMigrate(t, testCFG, wantVersion)

			// Ensure Current db version wasn't changed
			diff(t, wantVersion, getDBVersion(t))
		})
	})

	t.Run("DownTo0", func(t *testing.T) {
		wantVersion := 0
		t.Run("DRY RUN", func(t *testing.T) {
			testCFG.DryRun = true
			testMigrate(t, testCFG, wantVersion)
			testCFG.DryRun = false

			// Ensure Current db version wasn't changed
			diff(t, expectedCurrentVersion, getDBVersion(t))
		})
		t.Run("WET RUN", func(t *testing.T) {
			testMigrate(t, testCFG, wantVersion)

			// Ensure Current db version was changed
			diff(t, wantVersion, getDBVersion(t))
		})
	})
	// Tests that the locking mechanism works when multiple migrate
	// calls are coming in in parallel
	t.Run("Concurrent Up Migrations", func(t *testing.T) {
		wantVersion := 16
		t.Parallel()
		for i := 0; i < 5; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				testMigrate(t, testCFG, wantVersion)
				// Ensure Current db version was changed
				diff(t, wantVersion, getDBVersion(t))
			})
		}
		// Ensure db version is where you expect it to be
		diff(t, wantVersion, getDBVersion(t))
	})
	t.Run("Concurrent Down Migrations", func(t *testing.T) {
		wantVersion := 0
		t.Parallel()
		for i := 0; i < 5; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				testMigrate(t, testCFG, wantVersion)
				// Ensure Current db version was changed
				// (at the point in time testMigrate returns,
				//	doesnt mean this instance of testMigrate actually did any migrations)
				diff(t, wantVersion, getDBVersion(t))
			})
		}
		// Ensure db version is where you expect it to be
		diff(t, wantVersion, getDBVersion(t))
	})
}

func diff(t *testing.T, want, got interface{}) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func getDBVersion(t *testing.T) int {
	t.Helper()
	i, err := goose.EnsureDBVersion(db)
	if err != nil {
		t.Error(err)
	}
	return int(i)
}

func testMigrate(t *testing.T, testCFG migrate.Config, version int) {
	t.Helper()
	testCFG.LoggerOverride = migrate.TestingLogger{T: t}
	if err := migrate.Migrate(context.Background(), testCFG, src, int64(version)); err != nil {
		t.Error(err)
	}
}
