package common

import (
	. "github.com/saichler/console/golang/console/commands"
)

type IServiceManager interface {
	Shutdown()
	ConsoleId() *ConsoleId
}
