package testutil

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/xcls/gocard/config"
	"github.com/xcls/gocard/dbutil"
	"github.com/xcls/gocard/migrations"
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

	runner := nomadpg.NewRunner(testDB, migrations.Migrations)

	// FIXME(xcls): Rollback all migrations?

	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}

	return testDB
}

func ResetDB(t *testing.T, db *sql.DB) {
	tables := []string{
		"answers",
		"reviews",
		"cards",
		"decks",
		"user_sessions",
		"users",
	}
	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := db.Exec(sql); err != nil {
			t.Log(sql)
			t.Fatal(err)
		}
	}
}
