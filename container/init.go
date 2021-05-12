package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {

	commands := readUserCommands()
	if commands == nil || len(commands) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}
	log.Println(commands)

	// setUpMount()

	path, err := exec.LookPath(commands[0])
	if err != nil {
		log.Fatalf("Exec loop path error %v", err)
		return err
	}
	log.Printf("Find path %s", path)

	if err := syscall.Exec(path, commands[0:], os.Environ()); err != nil {
		log.Fatalln(err.Error())
		return err
	}
	return nil
}

func setUpMount() {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
}

// readUserCommands 从管道中获取参数
func readUserCommands() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Fatalf("init read pipe error %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

