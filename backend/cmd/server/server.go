package main

import (
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"github.com/lhhong/timeseries-query/backend/pkg/http"
	"github.com/lhhong/timeseries-query/backend/pkg/repository"
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
	config.LoadConfig(cmd)
	repository.LoadDb()
}

func init() {
	log.Println("Started init")
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished init")
}
func main() {
	log.Println("Starting server")
	http.StartServer()
}
