package container

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

func Remove(name string) {
	container, err := getContainer(name)
	if err != nil {
		logrus.Errorf("rm %v error: ", name, err)
		return
	}

	if !container.IsStop() {
		logrus.Errorf("couldn't rm running container %v", name)
		return
	}
	dirUrl := fmt.Sprintf(DefaultContainerLocation, name)
	if err := os.RemoveAll(dirUrl); err != nil {
		logrus.Errorf("rm files %v error: ", dirUrl, err)
		return
	}

	DeleteWorkSpace(container.Volume, name)
}
