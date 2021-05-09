// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
					syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(100),
		Gid: uint32(100),
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(-1)
}
