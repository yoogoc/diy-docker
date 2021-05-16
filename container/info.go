package container

import (
	"diy-docker/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

func RecordContainer(pid int, commands []string, name string) (string, error) {
	id := utils.RandStringBytes(10)
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	commandString := strings.Join(commands, ",")
	if name == "" {
		name = id
	}
	container := &Container{
		Pid:         strconv.Itoa(pid),
		Id:          id,
		Name:        name,
		Command:     commandString,
		CreatedTime: createdAt,
		Status:      RUNNING,
	}

	marshal, err := json.Marshal(container)
	if err != nil {
		logrus.Errorf("marshal container error: %v", err)
		return "", err
	}
	jsonStr := string(marshal)

	dirUrl := fmt.Sprintf(DefaultContainerLocation, name)

	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		logrus.Errorf("mkdir %v error: %v", dirUrl, err)
		return "", err
	}

	fileName := dirUrl + "/" + ConfigName

	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("create %v error: %v", fileName, err)
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Errorf("close %v error: %v", file.Name(), err)
		}
	}(file)

	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("write file %v error: %v", file.Name(), err)
		return "", err
	}
	return name, nil
}

func DeleteContainer(name string) {
	dirUrl := fmt.Sprintf(DefaultContainerLocation, name)
	_ = os.RemoveAll(dirUrl)
}