package container

var (
	RUNNING                  = "running"
	STOP                     = "stopping"
	EXIT                     = "exited"
	DefaultContainerLocation = "/var/run/diydocker/%s/"
	ConfigName               = "config.json"
	LogFile                  = "container.log"
	RootUrl                  = "/root"
	MntUrl                   = "/root/mnt/%s"
	WriteLayerUrl            = "/root/writeLayer/%s"
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
