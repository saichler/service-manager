package console

import (
	. "github.com/saichler/utils/golang"
	"net"
	"strconv"
)

type Console struct {
	socket net.Listener
}

func NewConsole(port int) (*Console, error) {
	socket, e := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if e != nil {
		Error("Failed to bind to console port:" + strconv.Itoa(port))
		return nil, e
	}
	console := &Console{}
	console.socket = socket
	return console, nil
}

func (c *Console) waitForConnection() {
	conn,e:=c.socket.Accept()
}