package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"gobase/cmd/start"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Go Base Services",
		Short: "gobase",
		Long:  "gobase - Backend",
	}
)

func Execute() {
	start.Cmd().Flags().StringP("config", "c", "config/file", "Config dir i.e. config/file")

	rootCmd.AddCommand(start.Cmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Error: \n", err.Error())
	}
}
