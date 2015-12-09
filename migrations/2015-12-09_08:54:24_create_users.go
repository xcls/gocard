package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-12-09_08:54:24",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
			CREATE TABLE users (
				id serial PRIMARY KEY,
				email text NOT NULL CHECK(length(email) < 400) UNIQUE,
				encrypted_password text NOT NULL CHECK(length(encrypted_password) < 5000),
				created_at timestamp with time zone DEFAULT(current_timestamp),
				updated_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec("DROP TABLE users CASCADE")
			return err
		},
	}
	Migrations.Add(migration)
}
