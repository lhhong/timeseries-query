package main

import (
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"github.com/lhhong/timeseries-query/backend/pkg/http"
	"github.com/spf13/cobra"
	"log"
)

func RootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use: "example",
		Run: run,
	}

	// this is where we will configure everything!
	rootCmd.Flags().StringP("config", "c", "", "Configuration file to use")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	c, _ := config.GetConfig(cmd)
	log.Println(c.Database.Hostname)
	log.Println(c.HttpServer.Port)
	http.StartServer(&c.HttpServer)
}

func main() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
