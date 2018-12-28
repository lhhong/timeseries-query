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
	rootCmd.Flags().StringP("config", "c", "", "Configuration file to use")
	rootCmd.Flags().StringP("groupname", "n", "", "Name of data group to smooth")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	conf := config.GetConfig(cmd)
	repo := repository.LoadDb(&conf.Database)
	group, _ := cmd.Flags().GetString("groupname")
	loader.SmoothAllInGroup(repo, group)
}

func main() {
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
