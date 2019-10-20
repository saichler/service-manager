package main

import "github.com/saichler/service-manager/golang/service-manager"

func main() {
	serviceManager, e := service_manager.NewServiceManager()
	if e != nil {
		return
	}
	serviceManager.WaitForShutdown()
}
