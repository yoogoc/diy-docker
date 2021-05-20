package cmd

import (
	"diy-docker/cgroups/subsystems"
	"diy-docker/container"
	"diy-docker/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tty    bool
	detach bool
	volume string
	containerName string
	res    = subsystems.ResourceConfig{}
)

func init() {

	runCmd.PersistentFlags().BoolVarP(&tty, "ti", "", false, "enable tty")
	runCmd.PersistentFlags().StringVar(&res.CpuSet, "cpuset", "", "cpu limit")
	runCmd.PersistentFlags().StringVar(&res.CpuShare, "cpushare", "", "cpu share")
	runCmd.PersistentFlags().StringVarP(&res.MemoryLimit, "memory", "m", "", "memory limit")
	runCmd.PersistentFlags().StringVarP(&volume, "volume", "v", "", "volume")
	runCmd.PersistentFlags().BoolVarP(&detach, "detach", "d", false, "detach run")
	runCmd.PersistentFlags().StringVar(&containerName, "name", "", "set name")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Create a container",
	Long:               `Create a container with namespace and cgroups limit diydocker run - t i [command ]`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Args:               cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if detach && tty {
			return fmt.Errorf("ti and d can't both provided")
		}

		Run(args)
		return nil
	},
}

func Run(commands []string) {
	id := utils.RandStringBytes(10)
	if containerName == "" {
		containerName = id
	}
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, commands[0])

	if parent == nil {
		log.Fatal("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("process start error: %v", err)
	}

	if _, err := container.RecordContainer(parent.Process.Pid, commands[1:], containerName, id, volume); err != nil {
		logrus.Errorf("record container error: %v", err)
		return
	}

	// cgroupManager := cgroups.NewCgroupManager("diy-docker-cgroup")
	// defer func(cgroupManager *cgroups.CgroupManager) {
	// 	log.Print("removing cgroup")
	// 	_ = cgroupManager.Destroy()
	// }(cgroupManager)
	//
	// if err := cgroupManager.Set(&res); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
	// 	log.Fatal(err)
	// }

	sendInitCommand(commands, writePipe)

	if tty {
		if err := parent.Wait(); err != nil {
			logrus.Errorf("process wait error: %v", err)
		}
		container.DeleteContainer(containerName)
		container.DeleteWorkSpace(volume, containerName)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Printf("command all is %s", command)
	_, err := writePipe.WriteString(command)
	if err != nil {
		log.Fatal(err)
	}

	err = writePipe.Close()
	if err != nil {
		log.Fatal(err)
	}
}
