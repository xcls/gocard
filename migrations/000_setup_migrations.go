package migrations

import (
	"database/sql"
	"log"

	"github.com/mcls/gocard/config"
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
	// Setup postgres driver
	_ "github.com/lib/pq"
)

var Migrations *nomad.List

// Context will be available to each migration and should be used to provide
// access to the database
type Context struct {
	DB *sql.DB
}

// This struct will be used as an argument to each migrations Up/Down func.
// You can use this to get access to your database.
var context = &Context{}

func init() {
	db, err := sql.Open("postgres", config.DatabaseUrl())
	if err != nil {
		log.Fatal(err)
	}
	context.DB = db
	Migrations = nomad.NewList(pg.NewVersionStore(db), context)
}
