package container

import (
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func CommitContainer(imageName string)  {
	mnt := "/root/mnt"
	imageTar := "/root/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mnt, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder %v error: %v", imageTar, err)
	}
}