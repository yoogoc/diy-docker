package container

import (
	"diy-docker/utils"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	ENV_EXEC_PID = "docker_pid"
	ENV_EXEC_CMD = "docker_cmd"
)

func Exec(name string, commands []string) {
	pid, err := GetContainerPidByName(name)
	if err != nil {
		logrus.Errorf("get container pid error: %v", err)
		return
	}
	cmdStr := strings.Join(commands, " ")

	logrus.Infof("container pid %s", pid)
	logrus.Infof("command %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := os.Setenv(ENV_EXEC_PID, pid); err != nil {
		logrus.Errorf("set env ENV_EXEC_CMD error: %v", err)
		return
	}
	if err := os.Setenv(ENV_EXEC_CMD, cmdStr); err != nil {
		logrus.Errorf("set env ENV_EXEC_CMD error: %v", err)
		return
	}

	envs := utils.GetEnvsByPid(pid)
	cmd.Env = append(os.Environ(), envs...)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("run exec error: %v", err)
		return
	}
}

func GetContainerPidByName(name string) (string, error) {
	dirUrl := fmt.Sprintf(DefaultContainerLocation, name)
	configFilePath := dirUrl + ConfigName
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		logrus.Errorf("read container config file error: %v", err)
		return "", err
	}
	var container Container

	if err := json.Unmarshal(content, &container); err != nil {
		logrus.Errorf("unmarshal config file error: %v", err)
		return "", err
	}
	return container.Pid, nil
}