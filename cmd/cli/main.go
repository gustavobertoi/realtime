package main

import (
	"context"
	"log"

	"github.com/common-nighthawk/go-figure"
	"github.com/gustavobertoi/realtime/internal/database"
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
		figure.NewFigure("Realtime", "isometric1", true).Print()
		ctx := context.Background()
		client, err := database.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close(ctx)
		migrations, err := database.NewMigrations(client.GetConn())
		if err != nil {
			log.Fatal(err)
		}
		if err := migrations.Init(ctx); err != nil {
			log.Fatal(err)
		}
		if err := migrations.Up(ctx); err != nil {
			log.Fatal(err)
		}
	},
}

var migrationDown = &cobra.Command{
	Use:   "migrate:down",
	Short: "Dropping down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		figure.NewFigure("Realtime", "isometric1", true).Print()
		ctx := context.Background()
		client, err := database.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close(ctx)
		migrations, err := database.NewMigrations(client.GetConn())
		if err != nil {
			log.Fatal(err)
		}
		if err := migrations.Init(ctx); err != nil {
			log.Fatal(err)
		}
		if err := migrations.Down(ctx); err != nil {
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
