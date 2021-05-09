package cmd

import (
	"diy-docker/container"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	tty bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&tty, "ti", false, "enable tty")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Create a container",
	Long:  `Create a container with namespace and cgroups limit diydocker run - t i [command ]`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing container command")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		Run(tty, args[0])
	},
}

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}
	_ = parent.Wait()
	os.Exit(-1)
}

