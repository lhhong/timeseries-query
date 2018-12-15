package main

import (
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"github.com/lhhong/timeseries-query/backend/pkg/http"
	"github.com/lhhong/timeseries-query/backend/pkg/repository"
	"github.com/spf13/cobra"
	"log"
)

func rootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use: "example",
		Run: run,
	}

	rootCmd.Flags().StringP("config", "c", "", "Configuration file to use")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	conf := config.GetConfig(cmd)
	repo := repository.LoadDb(&conf.Database)
	http.StartServer(&conf.HTTPServer, repo)
}

func main() {
	log.Println("Starting server")
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
