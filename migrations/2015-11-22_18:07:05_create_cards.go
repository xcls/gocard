package migrations

import (
	"github.com/mcls/nomad"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-11-22_18:07:05",
		Up: func(ctx interface{}) error {
			c := ctx.(*Context)
			_, err := c.DB.Exec(`
			CREATE TABLE cards (
				id serial PRIMARY KEY,
				front text NOT NULL,
				back text NOT NULL,
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*Context)
			_, err := c.DB.Exec("DROP TABLE cards")
			return err
		},
	}
	Migrations.Add(migration)
}
