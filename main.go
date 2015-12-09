package main

import (
	"log"

	"github.com/mcls/gocard/app"
	"github.com/mcls/gocard/migrations"
	"github.com/mcls/gocard/seed"

	"github.com/mcls/nomad"
	"github.com/spf13/cobra"
)

var GocardCmd = &cobra.Command{
	Use:   "gocard",
	Short: "gocard command",
	Long:  `Long gocard description`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "starts server",
	Run: func(cmd *cobra.Command, args []string) {
		app.StartServer()
	},
}

var MigrationCmd = nomad.NewMigrationCmd(migrations.NewRunner(), "./migrations")

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeds the database with dev data",
	Run: func(cmd *cobra.Command, args []string) {
		err := seed.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	GocardCmd.AddCommand(ServerCmd)
	GocardCmd.AddCommand(MigrationCmd)
	GocardCmd.AddCommand(SeedCmd)
	GocardCmd.Execute()
}
