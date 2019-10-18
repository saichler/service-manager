package main

import "github.com/saichler/service-manager/golang/service"

func main() {
	serviceManager,e:=service.NewServiceManager()
	if e!=nil {
		return
	}
	serviceManager.WaitForShutdown()
}