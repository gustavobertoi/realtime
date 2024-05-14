package main

import (
	"log"

	"github.com/common-nighthawk/go-figure"
	"github.com/gustavobertoi/realtime/internal/database"
	"github.com/gustavobertoi/realtime/internal/database/migrations"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "realtime",
	Short:   "Realtime CLI",
	Version: "1.0.0-beta",
	Run: func(cmd *cobra.Command, args []string) {
		figure.NewFigure("Realtime", "isometric1", true).Print()
	},
}

var migrationUp = &cobra.Command{
	Use:   "migrate:up",
	Short: "Create all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.NewTursoDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if err := migrations.MigrationUp(db); err != nil {
			log.Fatal(err)
		}
	},
}

var migrationDown = &cobra.Command{
	Use:   "migrate:down",
	Short: "Drop all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.NewTursoDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if err := migrations.MigrationDown(db); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(migrationUp, migrationDown)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
