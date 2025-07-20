package jobs

import controllers "github.com/ankur12345678/uptime-monitor/Controllers"

type JobName string

const (
	MonitorWesbitesJob JobName = "monitor-websites"
	NotificationJob    JobName = "notify-users"
)

type JobInput struct {
	controllers.BaseController
}
