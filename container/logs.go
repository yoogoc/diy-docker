package container

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

func Logs(containerName string) {
	dirUrl := fmt.Sprintf(DefaultContainerLocation, containerName)
	logFilePath := dirUrl + LogFile

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		logrus.Errorf("read log file error: %v", err)
		return
	}

	_, _ = fmt.Fprint(os.Stdout, string(content))
}
