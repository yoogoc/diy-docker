package cmd

import (
	"diy-docker/container"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "container stop",
	Long:  `container stop`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container.Remove(args[0])
	},
}
