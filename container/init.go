package container

import (
	"log"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Fatalln(err.Error())
		return err
	}
	return nil
}
