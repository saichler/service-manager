package main

import (
	console2 "github.com/saichler/service-manager/golang/service/console"
	"time"
)

func main() {
	console2.NewConsole(20000, &ServiceNode{})
	time.Sleep(time.Second * 60)
}

type ServiceNode struct {
}

func (sn *ServiceNode) Name() string {
	return "Service"
}

func (sn *ServiceNode) HandleCommand(args []string) {

}

func (sn *ServiceNode) CommandList() []string {
	return nil
}
