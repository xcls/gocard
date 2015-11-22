package main

import (
	"log"
	"net/http"

	"github.com/mcls/gocard/migrations"

	"github.com/mcls/nomad"
	"github.com/spf13/cobra"
	"github.com/unrolled/render"
)

var renderer = render.New(render.Options{
	Directory: "templates",
	Layout:    "layout",
})

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
		startServer()
	},
}

var NewMigrationCmd = &cobra.Command{
	Use:   "migration:new",
	Short: "create migration",
	Run: func(cmd *cobra.Command, args []string) {
		m := nomad.NewMigrator("./migrations")
		m.Create(args[0])
	},
}

var RunMigrationCmd = &cobra.Command{
	Use:   "migration:run",
	Short: "run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrations.Run()
	},
}

var RollbackMigrationCmd = &cobra.Command{
	Use:   "migration:rollback",
	Short: "rollsback the latest migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrations.Rollback()
	},
}

func main() {
	GocardCmd.AddCommand(ServerCmd)
	GocardCmd.AddCommand(NewMigrationCmd)
	GocardCmd.AddCommand(RunMigrationCmd)
	GocardCmd.AddCommand(RollbackMigrationCmd)
	GocardCmd.Execute()
}

func startServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/cards/new", NewCardHandler)

	port := ":8080"
	log.Printf("Starting server on %q\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderer.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"lol": "lolsies",
		"hah": []string{"hello", "hello"},
	})
}

func NewCardHandler(w http.ResponseWriter, r *http.Request) {
	renderer.HTML(w, http.StatusOK, "cards/new", nil)
}
