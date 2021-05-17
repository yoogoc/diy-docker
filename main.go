package main

import (
	"diy-docker/cmd"
	_ "diy-docker/nsenter"
)

func main() {
	cmd.Execute()
}
