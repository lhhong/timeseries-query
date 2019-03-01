package main

import (
	"log"

	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/loader"
	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use: "example",
		Run: run,
	}

	// this is where we will configure everything!
	rootCmd.Flags().Bool("swift-data", false, "Load SWIFT astro dataset")
	rootCmd.Flags().Bool("index-only", false, "Index data from database")
	rootCmd.Flags().String("dir", "", "directory of desired dataset dataset")
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
	conf := config.GetConfig(cmd)
	repo := repository.LoadDb(&conf.Database)
	loader.LoadData(cmd, &conf.App, repo)
}

func main() {
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
