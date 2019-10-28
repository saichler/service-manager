package main

import (
	"fmt"
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/commands"
	message_handlers "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	"github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

var Service IService = &service.FileManager{}
var Commands IServiceCommands = &commands.FileManagerCommands{}
var Handlers IServiceMessageHandlers = &message_handlers.FileManagerHandlers{}

func NewVars() *message_handlers.FileManagerHandlers {

}

func main() {
	str:="[M=0,Ip=10.0.2.15,P=52000][T=Management Service,D=1]"
	sid:=&protocol.ServiceID{}
	sid.Parse(str)
	fmt.Println(sid)
}
