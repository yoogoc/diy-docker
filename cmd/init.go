package cmd

import (
	"diy-docker/container"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init container",
	Long:  `Init container process run user’s process in container . Do not call it outside`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("init come on")
		log.Printf("command %s", args[0])
		_ = container.RunContainerInitProcess(args[0], nil)
	},
}
