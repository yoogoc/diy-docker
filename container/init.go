package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func RunContainerInitProcess() error {

	commands := readUserCommands()
	if commands == nil || len(commands) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	setUpMount()

	path, err := exec.LookPath(commands[0])
	if err != nil {
		log.Fatalf("Exec loop path error %v", err)
		return err
	}
	log.Printf("Find path %s", path)

	if err := syscall.Exec(path, commands[0:], os.Environ()); err != nil {
		log.Fatalln(err.Error())
		return err
	}
	return nil
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("get pwd error: %v", err))
	}
	logrus.Infof("pwd is: %v", pwd)

	if err := pivotRoot(pwd); err != nil {
		// panic(fmt.Errorf("pivot root error: %v", err))
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount proc error: %v", err)
	}

	_ = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

// readUserCommands 从管道中获取参数
func readUserCommands() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Fatalf("init read pipe error %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

// pivotRoot 将进程切换至新的文件系统
func pivotRoot(root string) error {

	// 1. 先重新mount一下旧的root
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to self error: %v", err)
	}

	// 2. 创建一个临时文件夹来存旧的root
	pivotPath := filepath.Join(root, ".pivot_path")
	if err := os.Mkdir(pivotPath, 0777); err != nil {
		logrus.Errorf("pivotPath: %v, error: %v", pivotPath, err.Error())
		return err
	}

	// 3. 切换root文件系统
	if err := syscall.PivotRoot(root, pivotPath); err != nil {
		logrus.Errorf("root: %v, pivotPath: %v, pivot root err: %v", root, pivotPath, err)
		return fmt.Errorf("pivot root err: %v", err)
	}

	// 4. 切换至新fs的根目录
	if err := syscall.Chdir("/"); err != nil {
		logrus.Errorf("chdir / err: %v", err)
		return fmt.Errorf("chdir / err: %v", err)
	}

	// 5. 卸载旧的fs
	pivotPath = filepath.Join("/", ".pivot_path")
	if err := syscall.Unmount(pivotPath, syscall.MNT_DETACH); err != nil {
		logrus.Errorf("unmount old pivot_path err: %v", err)
		return fmt.Errorf("unmount old pivot_path err: %v", err)
	}

	// 6. 移除旧的目录
	return os.Remove(pivotPath)
}
