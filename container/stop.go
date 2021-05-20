package container

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func Stop(name string) {
	pid, err := GetContainerPidByName(name)
	if err != nil {
		logrus.Errorf("get container error: %v", err)
		return
	}

	pidInt, _ := strconv.Atoi(pid)


	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		logrus.Errorf("stop container %v error: %v", name, err)
		return
	}

	container, err := getContainer(name)
	if err != nil {
		logrus.Errorf("get container %v meta error: %v", name, err)
		return 
	}

	container.Status = STOP
	container.Pid = ""
	marshal, err := json.Marshal(container)
	if err != nil {
		logrus.Errorf("marshal container %v meta error: %v", name, err)
		return
	}
	dirUrl := fmt.Sprintf(DefaultContainerLocation, name)
	configFilePath := dirUrl + ConfigName

	if err := os.WriteFile(configFilePath, marshal, 0622); err != nil {
		logrus.Errorf("write container %v meta config error: %v", name, err)
	}
}
