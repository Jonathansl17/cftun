package servicecmd

import "github.com/Jonathansl17/cftun/internal/service"

const (
	useService = "service"
	useInstall = "install"

	orderInstall = 1
)

var actionOrder = map[service.Action]int{
	service.Status:  2,
	service.Restart: 3,
	service.Start:   4,
	service.Stop:    5,
	service.Enable:  6,
	service.Disable: 7,
	service.Logs:    8,
}
