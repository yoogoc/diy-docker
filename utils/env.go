package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	environ = "/proc/%s/environ"
)

func GetEnvsByPid(pid string) []string {
	envPath := fmt.Sprintf(environ, pid)
	bytes, err := os.ReadFile(envPath)
	if err != nil {
		logrus.Errorf("read environ file error: %v", err)
		return []string{}
	}
	return strings.Split(string(bytes), "\u0000")
}
