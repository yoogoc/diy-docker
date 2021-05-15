package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

func NewWorkSpace(root string, mnt, volume string) {
	CreateReadOnlyLayer(root)
	CreateWriteOnlyLayer(root)
	CreateMountPoint(root, mnt)
	if volume == "" {
		return
	}
	sourceV, targetV, err := vUrlExtract(volume)
	if err != nil {
		panic(err)
	}
	MountVolume(mnt, sourceV, targetV)
}

func MountVolume(mnt, sourceV, targetV string) {
	if err := os.Mkdir(sourceV, 0777); err != nil {
		logrus.Errorf("mkdir source dir %v error: %v", sourceV, err)
	}

	fullTargetV := mnt + targetV
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
func CreateReadOnlyLayer(root string) {
	bbUrl := root + "busybox/"
	bbtUrl := root + "busybox.tar"

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
func CreateWriteOnlyLayer(root string) {
	writeUrl := root + "writeUrl/"
	if err := os.Mkdir(writeUrl, 0777); err != nil {
		logrus.Errorf("mkdir %v error: %v", writeUrl, err)
	}
}

func CreateMountPoint(root, mnt string) {
	if err := os.Mkdir(mnt, 0777); err != nil {
		logrus.Errorf("mkdir %v error: %v", mnt, err)
	}

	dirs := "dirs=" + root + "writeUrl:" + root + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mnt)
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

func DeleteWorkSpace(rootUrl, mntUrl, volume string) {
	if volume != "" {
		_, targetV, err := vUrlExtract(volume)
		if err == nil {
			UmountVolume(mntUrl, targetV)
		}
	}

	DeleteMountPoint(mntUrl)
	DeleteWriteLayer(rootUrl)
}

func UmountVolume(mntUrl, targetV string) {
	fullTargetUrl := mntUrl + targetV
	cmd := exec.Command("umount", fullTargetUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run();err != nil {
		logrus.Errorf("umount fullTargetUrl error: %v", err)
	}
}

func DeleteMountPoint(mntUrl string) {
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

func DeleteWriteLayer(rootUrl string) {
	writeUrl := rootUrl + "writeUrl/"
	if err := os.RemoveAll(writeUrl); err != nil {
		logrus.Errorf("mkdir %v error: %v", writeUrl, err)
	}
}