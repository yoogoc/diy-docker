package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func NewParentProcess(tty bool, volume string, containerName, imageName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	if tty {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
	} else {
		dirUrl := fmt.Sprintf(DefaultContainerLocation, containerName)

		if err := os.MkdirAll(dirUrl, 0622); err != nil {
			logrus.Errorf("mkdir %v error %v", dirUrl, err)
			return nil, nil
		}

		stdLogFilePath := dirUrl + LogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		logrus.Infof("log by %v", stdLogFilePath)
		if err != nil {
			logrus.Errorf("create log file error %v", err)
			return nil, nil
		}

		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}

	NewWorkSpace(volume, imageName, containerName)

	cmd.Dir = fmt.Sprintf(MntUrl, containerName)
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
