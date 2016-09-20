package dbutil

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/xcls/gocard/config"
)

// Connect connects to a database url and returns the handle. Optionally
// accepts a database url to override the url. The connectUrl defaults to
// config.DatabaseUrl()
func Connect(connectUrl ...string) (*sql.DB, error) {
	url := config.DatabaseURL()
	if len(connectUrl) > 0 {
		url = connectUrl[0]
	}
	return sql.Open("postgres", url)
}
