package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteOnlyLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume == "" {
		return
	}
	sourceV, targetV, err := vUrlExtract(volume)
	if err != nil {
		panic(err)
	}
	MountVolume(containerName, sourceV, targetV)
}

func MountVolume(containerName, sourceV, targetV string) {
	if err := os.Mkdir(sourceV, 0777); err != nil {
		logrus.Errorf("mkdir source dir %v error: %v", sourceV, err)
	}

	mntUrl := fmt.Sprintf(MntUrl, containerName)
	fullTargetV := mntUrl + "/" + targetV
	if err := os.Mkdir(fullTargetV, 0777); err != nil {
		logrus.Errorf("mkdir target dir %v error: %v", fullTargetV, err)
	}

	dirs := "dirs=" + sourceV
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", fullTargetV)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:= cmd.Run(); err !=nil{
		logrus.Errorf("mount volume error: %v", fullTargetV, err)
	}
}

func vUrlExtract(volume string) (string, string, error) {
	vos := strings.Split(volume, ":")
	if len(vos) == 2 && vos[0] != "" && vos[1] != "" {
		return vos[0], vos[1], nil
	} else {
		return "", "", fmt.Errorf("volumn format error")
	}
}

// CreateReadOnlyLayer 创建只读层
func CreateReadOnlyLayer(imageName string) {
	bbUrl := RootUrl + "/" + imageName + "/"
	bbtUrl := RootUrl + "/" + imageName + ".tar"

	fileExists, err := FileExists(bbUrl)
	if err != nil {
		logrus.Errorf("check file exists error: %v", err)
	}
	if fileExists {
		return
	}

	if err := os.Mkdir(bbUrl, 0777); err != nil {
		logrus.Errorf("mkdir %s error: %v", bbUrl, err)
	}

	if _, err:= exec.Command("tar", "-xvf", bbtUrl, "-C", bbUrl).CombinedOutput(); err!=nil {
		logrus.Errorf("tar %s error: %v", bbtUrl, err)
	}
}

// CreateWriteOnlyLayer 创建只写层
func CreateWriteOnlyLayer(containerName string) {
	writeUrl := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeUrl, 0777); err != nil {
		logrus.Errorf("mkdir %v error: %v", writeUrl, err)
	}
}

func CreateMountPoint(containerName, imageName string) {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.Mkdir(mntUrl, 0777); err != nil {
		logrus.Errorf("mkdir %v error: %v", mntUrl, err)
	}

	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	tmpImageLocation := RootUrl + "/" + imageName

	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount error: %v", err)
	}
}

func FileExists(fileUrl string) (bool, error) {
	_, err := os.Stat(fileUrl)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		_, targetV, err := vUrlExtract(volume)
		if err == nil {
			UmountVolume(containerName, targetV)
		}
	}

	DeleteMountPoint(containerName)
	DeleteWriteLayer(containerName)
}

func UmountVolume(containerName, targetV string) {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	fullTargetUrl := mntUrl + targetV
	cmd := exec.Command("umount", fullTargetUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run();err != nil {
		logrus.Errorf("umount fullTargetUrl error: %v", err)
	}
}

func DeleteMountPoint(containerName string) {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run();err != nil {
		logrus.Errorf("umount error: %v", err)
	}

	if err := os.RemoveAll(mntUrl); err != nil {
		logrus.Errorf("rm dir %v error: %v", mntUrl, err)
	}
}

func DeleteWriteLayer(containerName string) {
	writeUrl := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeUrl); err != nil {
		logrus.Errorf("mkdir %v error: %v", writeUrl, err)
	}
}