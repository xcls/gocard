package dbutil

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mcls/gocard/config"
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", config.DatabaseUrl())
}
