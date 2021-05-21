package cmd

import (
	"diy-docker/container"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "container logs",
	Long:  `container logs`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container.Logs(args[0])
	},
}
