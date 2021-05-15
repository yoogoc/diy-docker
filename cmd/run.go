package cmd

import (
	"diy-docker/cgroups"
	"diy-docker/cgroups/subsystems"
	"diy-docker/container"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tty bool
	volume string
	res = subsystems.ResourceConfig{}
)

func init() {

	runCmd.PersistentFlags().BoolVarP(&tty, "ti", "", false, "enable tty")
	runCmd.PersistentFlags().StringVar(&res.CpuSet, "cpuset", "", "cpu limit")
	runCmd.PersistentFlags().StringVar(&res.CpuShare, "cpushare", "", "cpu share")
	runCmd.PersistentFlags().StringVarP(&res.MemoryLimit, "memory", "m", "", "memory limit")
	runCmd.PersistentFlags().StringVarP(&volume, "volume", "v", "", "volume")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Create a container",
	Long:  `Create a container with namespace and cgroups limit diydocker run - t i [command ]`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		Run(tty, args)
	},
}

func Run(tty bool, commands []string) {
	parent, writePipe := container.NewParentProcess(tty, volume)

	if parent == nil {
		log.Fatal("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("process start error: %v", err)
	}

	cgroupManager := cgroups.NewCgroupManager("diy-docker-cgroup")
	defer func(cgroupManager *cgroups.CgroupManager) {
		log.Print("removing cgroup")
		_ = cgroupManager.Destroy()
	}(cgroupManager)

	if err := cgroupManager.Set(&res); err != nil {
		log.Fatal(err)
	}

	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Fatal(err)
	}

	sendInitCommand(commands, writePipe)

	if err := parent.Wait(); err != nil {
		logrus.Errorf("process wait error: %v", err)
	}

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	log.Print("exit")
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
