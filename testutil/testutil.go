package testutil

import (
	"database/sql"
	"testing"

	"github.com/mcls/gocard/config"
	"github.com/mcls/gocard/dbutil"
	"github.com/mcls/gocard/migrations"
	nomadpg "github.com/mcls/nomad/pg"
)

var testDB *sql.DB

func ConnectDB(t *testing.T) *sql.DB {
	if testDB != nil {
		return testDB
	}
	return reconnectDB(t)
}

func reconnectDB(t *testing.T) *sql.DB {
	var err error
	testDB, err = dbutil.Connect(config.DatabaseTestURL())
	if err != nil {
		t.Fatal(err)
	}
	return testDB
}

func ResetDB(t *testing.T, db *sql.DB) {
	runner := nomadpg.NewRunner(db, migrations.Migrations)
	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM cards;"); err != nil {
		t.Fatal(err)
	}
}
