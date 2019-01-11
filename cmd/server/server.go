package main

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/http"
	"github.com/lhhong/timeseries-query/pkg/querycache"
	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
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
	cs := querycache.NewCacheStore(&conf.Redis)
	http.StartServer(&conf.HTTPServer, repo, cs)
}

func main() {
	log.Println("Starting server")
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
