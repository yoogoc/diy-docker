package container

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
)

func Logs(containerName string) {
	dirUrl := fmt.Sprintf(DefaultContainerLocation, containerName)
	logFilePath := dirUrl + LogFile

	file, err := os.Open(logFilePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		logrus.Errorf("open log file error: %v", err)
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read log file error: %v", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, string(content))
}
