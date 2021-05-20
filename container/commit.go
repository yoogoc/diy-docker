package container

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func CommitContainer(imageName, containerName string)  {
	mntUrl := fmt.Sprintf(MntUrl, containerName) + "/"
	imageTar := RootUrl + "/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntUrl, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder %v error: %v", imageTar, err)
	}
}
