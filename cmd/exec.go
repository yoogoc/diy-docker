package cmd

import (
	"diy-docker/container"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "container exec",
	Long:  `container exec`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv(container.ENV_EXEC_PID) == "" {
			logrus.Info("1 exec")
			container.Exec(args[0], args[1:])
		} else {
			logrus.Info("2 exec")
		}
	},
}

