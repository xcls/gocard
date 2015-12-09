package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-12-09_15:31:42",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
			CREATE TABLE user_sessions (
				id serial PRIMARY KEY,
				uid text NOT NULL CHECK(length(uid) < 300) UNIQUE,
				user_id integer NOT NULL REFERENCES users ON DELETE CASCADE,
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec("DROP TABLE user_sessions CASCADE")
			if err != nil {
				return err
			}
			return nil
		},
	}
	Migrations.Add(migration)
}
