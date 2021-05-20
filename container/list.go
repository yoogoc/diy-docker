package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
)

func ListContainers() {
	dirUrl := fmt.Sprintf(DefaultContainerLocation, "")
	dirUrl = dirUrl[:len(dirUrl)-1]
	files, err := ioutil.ReadDir(dirUrl)
	if err != nil {
		logrus.Errorf("read container location error: %v", err)
		return
	}
	var containers []*Container

	for _, file := range files {
		container, err := getContainer(file.Name())
		if err != nil {
			logrus.Errorf("read %v error: %v", dirUrl, err)
			return
		}
		containers = append(containers, container)
	}
	writer := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	_, _ = fmt.Fprintf(writer, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")

	for _, c := range containers {
		_, _ = fmt.Fprintf(
			writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			c.Id, c.Name, c.Pid, c.Status, c.Command, c.CreatedTime,
		)
	}

	if err := writer.Flush(); err != nil {
		logrus.Errorf("flush error: %v", err)
		return
	}
}

func getContainer(fileName string) (*Container, error) {
	configDir := fmt.Sprintf(DefaultContainerLocation, fileName)
	configDir = configDir + ConfigName
	content, err := os.ReadFile(configDir)
	if err != nil {
		logrus.Errorf("read container config file error: %v", err)
		return nil, err
	}

	var container Container

	if err := json.Unmarshal(content, &container); err != nil {
		logrus.Errorf("unmarshal container config file error: %v", err)
		return nil, err
	}
	return &container, nil
}
