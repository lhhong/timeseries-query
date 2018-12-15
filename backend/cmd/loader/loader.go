package main

import (
	"github.com/lhhong/timeseries-query/backend/pkg/config"
	"github.com/lhhong/timeseries-query/backend/pkg/loader"
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
	rootCmd.Flags().StringP("datafile", "f", "", "File of the data")
	rootCmd.Flags().StringP("groupname", "n", "", "Name of data group")
	rootCmd.Flags().IntP("series", "s", -1, "Column number of series name")
	rootCmd.Flags().IntP("date", "d", -1, "Column number of date")
	rootCmd.Flags().IntP("time", "t", -1, "Column number of time")
	rootCmd.Flags().IntP("value", "v", -1, "Column number of series value")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config.LoadConfig(cmd)
	repository.LoadDb()
	loader.LoadData(cmd)
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
}
