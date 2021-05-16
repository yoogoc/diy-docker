package cmd

import (
	"diy-docker/container"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit a container to image",
	Long:  `commit a container to image`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container.CommitContainer(args[0])
	},
}
