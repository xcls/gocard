package main

import (
	"github.com/mcls/gocard/app"
	"github.com/mcls/gocard/migrations"

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

var MigrationCmd = nomad.NewMigrationCmd(migrations.Migrations, "./migrations")

func main() {
	GocardCmd.AddCommand(ServerCmd)
	GocardCmd.AddCommand(MigrationCmd)
	GocardCmd.Execute()
}
