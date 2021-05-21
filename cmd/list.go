package cmd

import (
	"diy-docker/container"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:                "ps",
	Short:              "list all the containers",
	Long:               `list all the containers`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Args:               cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		container.ListContainers()
	},
}
