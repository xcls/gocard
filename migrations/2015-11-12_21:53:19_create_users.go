package migrations

import (
	"fmt"

	"github.com/mcls/nomad"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-11-12_21:53:19",
		Up: func(ctx interface{}) error {
			c := ctx.(Context)
			fmt.Println("Up")
			fmt.Println(c)
			return nil
		},
		Down: func(ctx interface{}) error {
			c := ctx.(Context)
			fmt.Println("Down")
			fmt.Println(c)
			return nil
		},
	}
	Migrations.Add(migration)
}
