package migrations

import (
	"fmt"

	"github.com/mcls/nomad"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-11-22_18:07:05",
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
