package main

import (
	"github.com/saichler/service-manager/golang/common"
	"github.com/saichler/service-manager/golang/file-manager/service"
)

var Service common.IService = &service.FileManager{}
