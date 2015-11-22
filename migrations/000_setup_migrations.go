package migrations

import (
	"database/sql"
	"log"
	"os"

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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	context.DB = db
	Migrations = nomad.NewList(pg.NewVersionStore(db))
}

// Run pending migrations
func Run() {
	Migrations.Run(context)
}
