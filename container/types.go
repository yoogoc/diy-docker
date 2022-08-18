package container

var (
	RUNNING                  = "running"
	STOP                     = "stopping"
	EXIT                     = "exited"
	DefaultContainerLocation = "/var/run/diydocker/%s/"
	ConfigName               = "config.json"
	LogFile                  = "container.log"
	RootUrl                  = "/root/diydocker"
	MntUrl                   = "/root/diydocker/mnt/%s"
	TmpWorkUrl               = "/root/diydocker/tmp/work/%s"
	WriteLayerUrl            = "/root/diydocker/writeLayer/%s"
)

type Container struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"created_time"`
	Status      string `json:"status"`
	Volume      string `json:"volume"`
}

func (c *Container) IsStop() bool {
	return c.Status == STOP
}
